package brbdxr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	//"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	brbdxrpbdsl "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/dsl"
	brbdxrpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/msgs"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/codingpb"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	rs "github.com/filecoin-project/mir/pkg/rs_ezpwd"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/mathutil"
	"github.com/pkg/errors"
	rs1 "github.com/vivint/infectious"
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

// GetN returns the total number of nodes.
func (params *ModuleParams) GetN() int {
	return len(params.AllNodes)
}

// GetF returns the maximum tolerated number of faulty nodes.
func (params *ModuleParams) GetF() int {
	return (params.GetN() - 1) / 3
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
	sentEcho  bool
	sentReady bool
	delivered bool

	echos                 [][]byte
	echoMessagesCount     map[string]map[string]int
	echoAccumulatorByHash map[string]*SingleAccumulator
	echosMaxAccumulator   DualAccumulator

	readys              []rs1.Share
	receivedReady       []bool
	readyMessagesCount  map[string]int
	readyMaxAccumulator SingleAccumulator
	lastDecodeAttempt   int
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

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	state := make(map[int64]*moduleState)
	lastId := int64(0)

	// upon event <brb, Broadcast | m> do
	brbpbdsl.UponBroadcastRequest(m, func(id int64, data []byte) error {
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.StartMessage(mc.Self, id, data), params.AllNodes)
		return nil
	})

	brbdxrpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, hdata []byte) error {
		if id <= lastId {
			return nil
		}
		if from == params.Leader {
			initialize(&state, id, params.GetN())

			if !state[id].sentEcho {
				state[id].sentEcho = true

				nDataShards := params.GetN() - 2*params.GetF()

				data := make([]byte, mathutil.Pad(4+len(hdata), nDataShards))
				binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
				copy(data[4:], hdata)

				dsl.HashOneMessage(m, mc.Hasher, [][]byte{hdata}, &hashInitialMessageContext{id: id, data: data})
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

		encoded, err := encode(int64(params.GetN()), int64(params.GetF()), context.data)
		if err != nil {
			return err
		}
		for i, node := range params.AllNodes {
			eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.EchoMessage(mc.Self, context.id, hash, encoded[i]), []t.NodeID{node})
		}
		return nil
	})

	brbdxrpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, hash, chunk []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id, params.GetN())
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

	brbdxrpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, hash, chunk []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id, params.GetN())
		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}
		if !state[id].receivedReady[fromId] {
			state[id].receivedReady[fromId] = true
			state[id].readys = append(state[id].readys, rs1.Share{
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
			if (currentState.echosMaxAccumulator.count > (params.GetN()+params.GetF())/2) && currentState.sentReady == false {
				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(
					mc.Self, id, currentState.echosMaxAccumulator.hash,
					currentState.echosMaxAccumulator.chunk,
				), params.AllNodes)
			}

			// upon receiving t + 1 ⟨READY, *, h⟩ messages and not having sent a
			// READY message do
			if currentState.readyMaxAccumulator.count > params.GetF() && currentState.sentReady == false {
				// Wait for t + 1 matching ⟨ECHO,m_i,h⟩
				accumulator := currentState.echoAccumulatorByHash[string(currentState.readyMaxAccumulator.value)]
				if accumulator != nil && accumulator.count > params.GetF() {
					currentState.sentReady = true
					eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(
						mc.Self, id, currentState.readyMaxAccumulator.value,
						accumulator.value,
					), params.AllNodes)
				}
			}

			// Online error correction
			if len(currentState.readys) > 2*params.GetF() && len(currentState.readys) > currentState.lastDecodeAttempt && currentState.delivered == false {
				currentState.lastDecodeAttempt = len(currentState.readys)
				shares := make([]*codingpb.Share, 0)
				for _, rd := range currentState.readys {
					shares = append(shares, &codingpb.Share{
						Number: int64(rd.Number),
						Chunk:  rd.Data,
					})
				}

				success, decoded := decode(int64(params.GetN()), int64(params.GetF()), shares)
				if success {
					size := binary.LittleEndian.Uint32(decoded[:4])
					decoded = decoded[4 : 4+size]

					dsl.HashOneMessage(m, mc.Hasher, [][]byte{decoded}, &hashVerificationContext{id: id, output: decoded})
				}
			}
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

func encode(n, f int64, data []byte) ([][]byte, error) {
	encoded := rs.EncodeWrapper(n, f, data)

	if len(encoded) == 0 {
		return nil, errors.New("Failed to encode")
	}

	return encoded, nil
}

func decode(n, f int64, inputShares []*codingpb.Share) (bool, []byte) {
	shares := make([][]byte, n)
	missing := make([]int, 0)
	for _, share := range inputShares {
		shares[share.Number] = share.Chunk
	}
	for i := range shares {
		if shares[i] == nil {
			missing = append(missing, i)
		}
	}
	decoded := rs.DecodeWrapper(n, f, shares, missing)

	if len(decoded) == 0 {
		return false, nil
	} else {
		return true, decoded
	}
}

func initialize(state *map[int64]*moduleState, id int64, n int) {
	if _, ok := (*state)[id]; !ok {
		(*state)[id] = &moduleState{
			sentEcho:              false,
			sentReady:             false,
			delivered:             false,
			echos:                 make([][]byte, n),
			echoMessagesCount:     make(map[string]map[string]int),
			echoAccumulatorByHash: make(map[string]*SingleAccumulator),
			echosMaxAccumulator:   DualAccumulator{},
			readys:                make([]rs1.Share, 0),
			receivedReady:         make([]bool, n),
			readyMessagesCount:    make(map[string]int),
			readyMaxAccumulator:   SingleAccumulator{},
			lastDecodeAttempt:     0,
		}
	}
}

type hashInitialMessageContext struct {
	id   int64
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
