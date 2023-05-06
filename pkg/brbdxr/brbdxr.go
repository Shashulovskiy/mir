package brbdxr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbdxrpbdsl "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/dsl"
	brbdxrpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/msgs"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/mathutil"
	rs "github.com/vivint/infectious"
	"strconv"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self     t.ModuleID // id of this module
	Consumer t.ModuleID // id of the module to send the "Deliver" event to
	Net      t.ModuleID
	Crypto   t.ModuleID
	Hasher   t.ModuleID
	Coder    t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:     "brbdxr",
		Consumer: consumer,
		Net:      "net",
		Crypto:   "crypto",
		Hasher:   "hasher",
		Coder:    "coder",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of BRB, used to prevent cross-instance replay attacks
	AllNodes    []t.NodeID // the list of participating nodes
	Leader      t.NodeID   // the id of the leader of the instance
}

type DualAccumulator struct {
	hash  []byte
	chunk []byte
	count int
}

type SingleAccumulator struct {
	value []byte
	count int
}

// moduleState represents the state of the brb module.
type moduleState struct {
	n int

	sentEcho  bool
	sentReady bool
	delivered bool

	echos                 [][]byte
	echoMessagesCount     map[string]map[string]int
	echoAccumulatorByHash map[string]*SingleAccumulator
	echosMaxAccumulator   DualAccumulator

	readys              []rs.Share
	receivedReady       []bool
	readyMessagesCount  map[string]int
	readyMaxAccumulator SingleAccumulator
	nextDecodeAttempt   int
}

func incrementAndUpdateEchoAccumulator(hash, chunk []byte, counts map[string]map[string]int, byHash map[string]*SingleAccumulator, accumulator *DualAccumulator) {
	hashString := string(hash)
	chunkString := string(chunk)

	if _, ok := counts[hashString]; !ok {
		counts[hashString] = make(map[string]int)
	}

	if _, ok := byHash[hashString]; !ok {
		byHash[hashString] = &SingleAccumulator{}
	}

	counts[hashString][chunkString]++

	if counts[hashString][chunkString] > (*accumulator).count {
		(*accumulator).count = counts[hashString][chunkString]
		(*accumulator).hash = hash
		(*accumulator).chunk = chunk
	}

	if counts[hashString][chunkString] > byHash[hashString].count {
		byHash[hashString].count = counts[hashString][chunkString]
		byHash[hashString].value = chunk
	}
}

func incrementAndUpdateReadyAccumulator(hash []byte, counts map[string]int, accumulator *SingleAccumulator) {
	hashString := string(hash)

	if _, ok := counts[hashString]; !ok {
		// TODO necessary?
		counts[hashString] = 0
	}

	counts[hashString]++

	if counts[hashString] > (*accumulator).count {
		(*accumulator).count = counts[hashString]
		(*accumulator).value = hash
	}
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, decodingStrategy string) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	state := make(map[int64]*moduleState)
	lastId := int64(0)

	// upon event <brb, Broadcast | m> do
	brbpbdsl.UponBroadcastRequest(m, func(id, n int64, data []byte) error {
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.StartMessage(mc.Self, id, n, data), params.AllNodes[:n])
		return nil
	})

	brbdxrpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id, n int64, hdata []byte) error {
		if id <= lastId {
			return nil
		}
		if from == params.Leader {
			initialize(&state, id, int(n), int(getF(n)))

			if !state[id].sentEcho {
				state[id].sentEcho = true

				nDataShards := int(n) - 2*int(getF(n))

				data := make([]byte, mathutil.Pad(4+len(hdata), nDataShards))
				binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
				copy(data[4:], hdata)

				dsl.HashOneMessage(m, mc.Hasher, [][]byte{hdata}, &hashInitialMessageContext{id: id, n: n, data: data})
				return nil
			} else {
				//fmt.Println("already sent echo")
				return nil
			}

		} else {
			//fmt.Println("received start message not from leader")
			return nil
		}
	})

	dsl.UponOneHashResult(m, func(hash []byte,
		context *hashInitialMessageContext) error {
		if context.id <= lastId {
			return nil
		}

		nDataShards := int(context.n - 2*getF(context.n))
		shardSize := len(context.data) / nDataShards

		dataWithPadding := make([]byte, nDataShards*shardSize)
		copy(dataWithPadding, context.data)

		encoded, err := encode(context.n, getF(context.n), dataWithPadding)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		for i, node := range params.AllNodes[:context.n] {
			eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.EchoMessage(mc.Self, context.id, context.n, hash, encoded[i]), []t.NodeID{node})
		}
		return nil
	})

	brbdxrpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id, n int64, hash, chunk []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id, int(n), int(getF(n)))
		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}
		if state[id].echos[fromId] == nil {
			state[id].echos[fromId] = chunk

			incrementAndUpdateEchoAccumulator(hash, chunk, state[id].echoMessagesCount, state[id].echoAccumulatorByHash, &state[id].echosMaxAccumulator)
		}
		return nil
	})

	brbdxrpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id, n int64, hash, chunk []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id, int(n), int(getF(n)))
		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}
		if !state[id].receivedReady[fromId] {
			state[id].receivedReady[fromId] = true
			state[id].readys = append(state[id].readys, rs.Share{
				Number: int(fromId),
				Data:   chunk,
			})

			incrementAndUpdateReadyAccumulator(hash, state[id].readyMessagesCount, &state[id].readyMaxAccumulator)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		for id, currentState := range state {
			if id <= lastId {
				continue
			}

			// upon receiving 2t + 1 ⟨ECHO,m_i,h⟩ matching messages and not having
			// sent a READY message do
			if (currentState.echosMaxAccumulator.count > (currentState.n+int(getF(int64(currentState.n))))/2) && currentState.sentReady == false {
				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(
					mc.Self, id, int64(currentState.n), currentState.echosMaxAccumulator.hash,
					currentState.echosMaxAccumulator.chunk,
				), params.AllNodes[:currentState.n])
			}

			// upon receiving t + 1 ⟨READY, *, h⟩ messages and not having sent a
			// READY message do
			if currentState.readyMaxAccumulator.count > int(getF(int64(currentState.n))) && currentState.sentReady == false {
				// Wait for t + 1 matching ⟨ECHO,m_i,h⟩
				accumulator := currentState.echoAccumulatorByHash[string(currentState.readyMaxAccumulator.value)]
				if accumulator.count > int(getF(int64(currentState.n))) {
					currentState.sentReady = true
					eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(
						mc.Self, id, int64(currentState.n), currentState.readyMaxAccumulator.value,
						accumulator.value,
					), params.AllNodes[:currentState.n])
				}
			}

			// Online error correction
			onlineErrorCorrection(decodingStrategy, currentState, m, mc, id, currentState.n, int(getF(int64(currentState.n))))
		}
		return nil
	})

	dsl.UponOneHashResult(m, func(hash []byte, context *hashVerificationContext) error {
		if context.id <= lastId {
			return nil
		}
		currentState := state[context.id]
		if bytes.Equal(hash, currentState.readyMaxAccumulator.value) {
			if !currentState.delivered {
				currentState.delivered = true
				brbpbdsl.Deliver(m, mc.Consumer, context.id, context.output)
				if context.id > lastId {
					for i := lastId; i <= context.id; i++ {
						delete(state, i)
					}
					lastId = context.id
				}
			}
		}

		return nil
	})

	return m, nil
}

func onlineErrorCorrection(strategy string, currentState *moduleState, m dsl.Module, mc *ModuleConfig, id int64, n, f int) {
	if strategy == "classic" {
		// Attempt error correction for each message received n-f...f (f times)
		if len(currentState.readys) > 2*f && currentState.delivered == false {
			tryCorrectErrors(currentState, m, mc, id, int64(n), int64(f))
		}
	} else if strategy == "optimized" {
		//if len(currentState.readys) == params.GetN()-2*params.GetF() && currentState.delivered == false {
		//	tryCorrectErasures(currentState, params, encoder, m, mc, id)
		//}
		if len(currentState.readys) >= currentState.nextDecodeAttempt && currentState.delivered == false {
			currentState.nextDecodeAttempt = currentState.nextDecodeAttempt + (n-currentState.nextDecodeAttempt+1)/2

			tryCorrectErrors(currentState, m, mc, id, int64(n), int64(f))
		}
	} else {
		panic("Unknown strategy")
	}
}

func tryCorrectErasures(currentState *moduleState, params *ModuleParams, encoder *rs.FEC, m dsl.Module, mc *ModuleConfig, id int64, n, f int) {
	res := make([]byte, len(currentState.readys[0].Data)*(n-2*f))

	readys := make([]rs.Share, 0)
	for _, rd := range currentState.readys {
		readys = append(readys, rd.DeepCopy())
	}

	err := encoder.Rebuild(currentState.readys, func(s rs.Share) {
		copy(res[s.Number*len(s.Data):], s.Data)
	})
	if err == nil {
		size := binary.LittleEndian.Uint32(res[:4])
		if int(size) > len(res) {
			return
		}
		res = res[4 : 4+size]

		dsl.HashOneMessage(m, mc.Hasher, [][]byte{res}, &hashVerificationContext{id: id, output: res})
	}
}

func tryCorrectErrors(currentState *moduleState, m dsl.Module, mc *ModuleConfig, id, n, f int64) {
	readys := make([]rs.Share, 0)
	for _, rd := range currentState.readys {
		readys = append(readys, rd.DeepCopy())
	}

	success, res := decode(n, f, currentState.readys)
	if success {
		size := binary.LittleEndian.Uint32(res[:4])
		if int(size) > len(res) {
			return
		}
		res = res[4 : 4+size]

		dsl.HashOneMessage(m, mc.Hasher, [][]byte{res}, &hashVerificationContext{id: id, output: res})
	}
}

func initialize(state *map[int64]*moduleState, id int64, n, f int) {
	if _, ok := (*state)[id]; !ok {
		(*state)[id] = &moduleState{
			n:                     n,
			sentEcho:              false,
			sentReady:             false,
			delivered:             false,
			echos:                 make([][]byte, n),
			echoMessagesCount:     make(map[string]map[string]int),
			echoAccumulatorByHash: make(map[string]*SingleAccumulator),
			echosMaxAccumulator:   DualAccumulator{},
			readys:                make([]rs.Share, 0),
			receivedReady:         make([]bool, n),
			readyMessagesCount:    make(map[string]int),
			readyMaxAccumulator:   SingleAccumulator{},
			nextDecodeAttempt:     n - f,
		}
	}
}

type hashInitialMessageContext struct {
	id   int64
	n    int64
	data []byte
}

type encodeMessageContext struct {
	id   int64
	hash []byte
}

type decodeMessageContext struct {
	id int64
}

type hashVerificationContext struct {
	id     int64
	output []byte
}

func getF(n int64) int64 {
	return (n - 1) / 3
}
