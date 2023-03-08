package brbdxr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbdxrpbdsl "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/dsl"
	brbdxrpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	rs "github.com/klauspost/reedsolomon"
	"github.com/pkg/errors"
	"math"
	"strconv"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self     t.ModuleID // id of this module
	Consumer t.ModuleID // id of the module to send the "Deliver" event to
	Net      t.ModuleID
	Crypto   t.ModuleID
	Hasher   t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:     "brbdxr",
		Consumer: consumer,
		Net:      "net",
		Crypto:   "crypto",
		Hasher:   "hasher",
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

type Accumulator = struct {
	hash  *[]byte
	chunk *[]byte
	count int
}

// brbDxrModuleState represents the state of the brb module.
type brbDxrModuleState struct {
	sentEcho  bool
	sentReady bool
	delivered bool

	echos               [][]byte
	echoMessagesCount   map[string]map[string]int
	echosMaxAccumulator Accumulator

	readys                [][]byte
	readyMessagesCount    map[string]map[string]int
	readyMessagesReceived int
	readyMaxAccumulator   Accumulator
}

func incrementAndUpdateAccumulator(hash, chunk *[]byte, counts *map[string]map[string]int, accumulator *Accumulator) {
	hashString := string(*hash)
	chunkString := string(*chunk)

	if _, ok := (*counts)[hashString]; !ok {
		(*counts)[hashString] = make(map[string]int)
	}

	(*counts)[hashString][chunkString]++

	if (*counts)[hashString][chunkString] > (*accumulator).count {
		(*accumulator).count = (*counts)[hashString][chunkString]
		(*accumulator).hash = hash
		(*accumulator).chunk = chunk
	}
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	encoder, err := rs.New(params.GetN()-2*params.GetF(), 2*params.GetF())

	if err != nil {
		return nil, errors.Wrap(err, "Unable to create coder")
	}

	state := make(map[int64]*brbDxrModuleState)
	lastId := int64(0)

	// upon event <brb, Broadcast | m> do
	brbdxrpbdsl.UponBroadcastRequest(m, func(id int64, data []byte) error {
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

				data := []byte{0, 0, 0, 0}
				binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
				data = append(data, hdata...)

				dsl.HashOneMessage(m, mc.Hasher, [][]byte{hdata}, &hashInitialMessageContext{id: id, data: data})
				return nil
			} else {
				return fmt.Errorf("already sent echo")
			}

		} else {
			return fmt.Errorf("received start message not from leader")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashInitialMessageContext) error {
		if context.id <= lastId {
			return nil
		}
		chunkSize := int(math.Ceil(float64(len(context.data)) / float64(params.GetN()-2*params.GetF())))
		encoded := make([][]byte, params.GetN())
		for i := range encoded {
			encoded[i] = make([]byte, chunkSize)
		}

		for i, in := range encoded[:params.GetN()-2*params.GetF()] {
			for j := range in {
				if i*chunkSize+j < len(context.data) {
					in[j] = context.data[i*chunkSize+j]
				} else {
					in[j] = 0
				}
			}
		}

		err := encoder.Encode(encoded)

		if err != nil {
			return err
		}

		for i, node := range params.AllNodes {
			eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.EchoMessage(mc.Self, context.id, hashes[0], encoded[i]), []t.NodeID{node})
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

			incrementAndUpdateAccumulator(&hash, &chunk, &state[id].echoMessagesCount, &state[id].echosMaxAccumulator)
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
		if state[id].readys[fromId] == nil {
			state[id].readys[fromId] = chunk
			state[id].readyMessagesReceived += 1

			incrementAndUpdateAccumulator(&hash, &chunk, &state[id].readyMessagesCount, &state[id].readyMaxAccumulator)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		for id, currentState := range state {
			if id <= lastId {
				continue
			}
			if (currentState.echosMaxAccumulator.count > (params.GetN()+params.GetF())/2 ||
				currentState.readyMaxAccumulator.count > params.GetF()) && currentState.sentReady == false {

				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(
					mc.Self, id, *currentState.echosMaxAccumulator.hash,
					*currentState.echosMaxAccumulator.chunk,
				), params.AllNodes)
			}

			if currentState.readyMessagesReceived > 2*params.GetF() && currentState.delivered == false {

				err := encoder.Reconstruct(currentState.readys)
				if err == nil {
					output := bytes.Join(currentState.readys, []byte{})

					size := binary.LittleEndian.Uint32(output[:4])
					output = output[4 : 4+size]

					dsl.HashOneMessage(m, mc.Hasher, [][]byte{output}, &hashVerificationContext{id: id, output: output})
				}
			}
		}
		return nil
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashVerificationContext) error {
		if context.id <= lastId {
			return nil
		}
		currentState := state[context.id]
		if bytes.Equal(hashes[0], *currentState.readyMaxAccumulator.hash) {
			if !currentState.delivered {
				currentState.delivered = true
				if context.id > lastId {
					lastId = context.id
				}
				brbdxrpbdsl.Deliver(m, mc.Consumer, context.id, context.output)
				delete(state, context.id)
			}
		} else {
			panic("?")
		}

		return nil
	})

	return m, nil
}

func initialize(state *map[int64]*brbDxrModuleState, id int64, n int) {
	if _, ok := (*state)[id]; !ok {
		(*state)[id] = &brbDxrModuleState{
			sentEcho:              false,
			sentReady:             false,
			delivered:             false,
			echos:                 make([][]byte, n),
			echoMessagesCount:     make(map[string]map[string]int),
			echosMaxAccumulator:   Accumulator{},
			readys:                make([][]byte, n),
			readyMessagesCount:    make(map[string]map[string]int),
			readyMessagesReceived: 0,
			readyMaxAccumulator:   Accumulator{},
		}
	}
}

type hashInitialMessageContext struct {
	id   int64
	data []byte
}

type hashVerificationContext struct {
	id     int64
	output []byte
}
