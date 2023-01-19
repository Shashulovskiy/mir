package brbencoded

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbencpbdsl "github.com/filecoin-project/mir/pkg/pb/brbencodedpb/dsl"
	brbencpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbencodedpb/msgs"
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
		Self:     "brbencoded",
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
	hash  []byte
	chunk []byte
	count int
}

type messageContent struct {
	hash  string
	chunk string
}

// brbEncodedModuleState represents the state of the brb module.
type brbEncodedModuleState struct {
	sentEcho  bool
	sentReady bool
	delivered bool

	echos               map[t.NodeID][]byte
	echoMessagesCount   map[messageContent]int
	echosMaxAccumulator Accumulator

	readys              map[t.NodeID][]byte
	readyMessagesCount  map[messageContent]int
	readyMaxAccumulator Accumulator
}

func incrementAndUpdateAccumulator(hash, chunk *[]byte, counts *map[messageContent]int, accumulator *Accumulator) {
	// TODO cant store []byte in a map
	var content = messageContent{
		hash:  string(*hash),
		chunk: string(*chunk),
	}

	(*counts)[content]++

	if (*counts)[content] > (*accumulator).count {
		(*accumulator).count = (*counts)[content]
		(*accumulator).hash = *hash
		(*accumulator).chunk = *chunk
	}
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	encoder, err := rs.New(params.GetN()-params.GetF(), params.GetF())

	if err != nil {
		return nil, errors.Wrap(err, "Unable to create coder")
	}

	// upon event <brb, Init> do
	state := brbEncodedModuleState{
		sentEcho:            false,
		sentReady:           false,
		delivered:           false,
		echos:               make(map[t.NodeID][]byte),
		echoMessagesCount:   make(map[messageContent]int),
		echosMaxAccumulator: Accumulator{},
		readys:              make(map[t.NodeID][]byte),
		readyMessagesCount:  make(map[messageContent]int),
		readyMaxAccumulator: Accumulator{},
	}

	// upon event <brb, Broadcast | m> do
	brbencpbdsl.UponBroadcastRequest(m, func(data []byte) error {
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		eventpbdsl.SendMessage(m, mc.Net, brbencpbmsgs.StartMessage(mc.Self, data), params.AllNodes)
		return nil
	})

	brbencpbdsl.UponStartMessageReceived(m, func(from t.NodeID, hdata []byte) error {
		if from == params.Leader && state.sentEcho == false {
			state.sentEcho = true

			data := []byte{0, 0, 0, 0}
			binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
			data = append(data, hdata...)

			dsl.HashOneMessage(m, mc.Hasher, [][]byte{hdata}, &hashInitialMessageContext{data: data})

			return nil
		} else {
			return fmt.Errorf("received start message not from leader or already sent echo")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashInitialMessageContext) error {
		chunkSize := int(math.Ceil(float64(len(context.data)) / float64(params.GetN()-params.GetF())))
		encoded := make([][]byte, params.GetN())
		for i := range encoded {
			encoded[i] = make([]byte, chunkSize)
		}

		for i, in := range encoded[:params.GetN()-params.GetF()] {
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
			eventpbdsl.SendMessage(m, mc.Net, brbencpbmsgs.EchoMessage(mc.Self, hashes[0], encoded[i]), []t.NodeID{node})
		}
		return nil
	})

	brbencpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, hash, chunk []byte) error {
		if state.echos[from] == nil {
			state.echos[from] = chunk

			incrementAndUpdateAccumulator(&hash, &chunk, &state.echoMessagesCount, &state.echosMaxAccumulator)
		}
		return nil
	})

	brbencpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, hash, chunk []byte) error {
		if state.readys[from] == nil {
			state.readys[from] = chunk

			incrementAndUpdateAccumulator(&hash, &chunk, &state.readyMessagesCount, &state.readyMaxAccumulator)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		if (state.echosMaxAccumulator.count > (params.GetN()+params.GetF())/2 ||
			state.readyMaxAccumulator.count > params.GetF()) && state.sentReady == false {

			state.sentReady = true
			eventpbdsl.SendMessage(m, mc.Net, brbencpbmsgs.ReadyMessage(
				mc.Self, state.echosMaxAccumulator.hash,
				state.echosMaxAccumulator.chunk,
			), params.AllNodes)
		}

		if len(state.readys) > 2*params.GetF() && state.delivered == false {

			decoded := make([][]byte, params.GetN())
			for key, value := range state.readys {
				id, err := strconv.ParseInt(key.Pb(), 10, 64)
				if err != nil {
					return err
				}
				decoded[id] = value
			}

			err := encoder.Reconstruct(decoded)
			if err != nil {
				return err
			}

			output := bytes.Join(decoded, []byte{})

			size := binary.LittleEndian.Uint32(output[:4])
			output = output[4 : 4+size]

			dsl.HashOneMessage(m, mc.Hasher, [][]byte{output}, &hashVerificationContext{output: output})
		}
		return nil
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashVerificationContext) error {
		if bytes.Equal(hashes[0], state.readyMaxAccumulator.hash) {
			if !state.delivered {
				state.delivered = true
				brbencpbdsl.Deliver(m, mc.Consumer, context.output)
			}
		} else {
			panic("?")
		}

		return nil
	})

	return m, nil
}

type hashInitialMessageContext struct {
	data []byte
}

type hashVerificationContext struct {
	output []byte
}
