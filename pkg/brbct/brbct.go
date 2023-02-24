package brbct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbctpb/dsl"
	brbmsgs "github.com/filecoin-project/mir/pkg/pb/brbctpb/msgs"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	rs "github.com/klauspost/reedsolomon"
	"github.com/pkg/errors"
	"math"
	"strconv"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self                t.ModuleID // id of this module
	Consumer            t.ModuleID // id of the module to send the "Deliver" event to
	Net                 t.ModuleID
	Crypto              t.ModuleID
	Hasher              t.ModuleID
	MerkleProofVerifier t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:                "brbct",
		Consumer:            consumer,
		Net:                 "net",
		Crypto:              "crypto",
		Hasher:              "hasher",
		MerkleProofVerifier: "merkleverify",
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
	value string
	count int
}

// brbHashModuleState represents the state of the brb module.
type brbHashModuleState struct {
	// this variable is not part of the original protocol description, but it greatly simplifies the code
	proposed []byte

	sentEcho  bool
	sentReady bool
	delivered bool

	echos                    map[string][][]byte
	echoMessagesCount        map[string]int
	echoMessagesAccumulator  Accumulator
	readyMessagesReceived    map[string][]bool
	readyMessagesCount       map[string]int
	readyMessagesAccumulator Accumulator
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	encoder, err := rs.New(params.GetN()-2*params.GetF(), 2*params.GetF())

	if err != nil {
		return nil, errors.Wrap(err, "Unable to create coder")
	}

	// upon event <brb, Init> do
	state := make(map[int64]*brbHashModuleState)
	lastId := int64(0)

	// upon event <brb, Broadcast | m> do
	brbpbdsl.UponBroadcastRequest(m, func(id int64, hdata []byte) error {
		if id <= lastId {
			return nil
		}
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}

		data := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
		data = append(data, hdata...)

		chunkSize := int(math.Ceil(float64(len(data)) / float64(params.GetN()-2*params.GetF())))
		encoded := make([][]byte, params.GetN())
		for i := range encoded {
			encoded[i] = make([]byte, chunkSize)
		}

		for i, in := range encoded[:params.GetN()-2*params.GetF()] {
			for j := range in {
				if i*chunkSize+j < len(data) {
					in[j] = data[i*chunkSize+j]
				} else {
					in[j] = 0
				}
			}
		}

		err := encoder.Encode(encoded)
		if err != nil {
			return err
		}

		dsl.MerkleBuildRequest(m, mc.MerkleProofVerifier, encoded, &hashInitialMessageContext{id: id, data: encoded})

		return nil
	})

	dsl.MerkleBuildResult(m, func(rootHash []byte, proofs []*commonpb.MerklePath, context *hashInitialMessageContext) error {
		if context.id <= lastId {
			return nil
		}
		for i, node := range params.AllNodes {
			eventpbdsl.SendMessage(m, mc.Net, brbmsgs.StartMessage(mc.Self, context.id, context.data[i], rootHash, proofs[i]), []t.NodeID{node})
		}
		return nil
	})

	// -----------
	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
		if id <= lastId {
			return nil
		}
		if from == params.Leader {
			initialize(&state, id)

			if !state[id].sentEcho {

				dsl.MerkleProofVerifyRequest(m, mc.MerkleProofVerifier, rootHash, proof, &hashStartMessageContext{id: id, chunk: chunk, rootHash: rootHash, proof: proof})

				return nil
			} else {
				return fmt.Errorf("already received start from leader")
			}
		}
		return fmt.Errorf("received start message not from leader")
	})

	dsl.MerkleProofVerifyResult(m, func(result bool, context *hashStartMessageContext) error {
		if context.id <= lastId {
			return nil
		}
		currentState := state[context.id]
		if result && !currentState.sentEcho {
			currentState.sentEcho = true

			eventpbdsl.SendMessage(m, mc.Net, brbmsgs.EchoMessage(mc.Self, context.id, context.chunk, context.rootHash, context.proof), params.AllNodes)

			return nil
		} else {
			return fmt.Errorf("received invalid chunk or already sent echo")
		}
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id)

		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}

		dsl.MerkleProofVerifyRequest(m, mc.MerkleProofVerifier, rootHash, proof, &hashEchoMessageContext{id: id, fromId: fromId, chunk: chunk, rootHash: rootHash, proof: proof})
		return nil
	})

	dsl.MerkleProofVerifyResult(m, func(result bool, context *hashEchoMessageContext) error {
		if context.id <= lastId {
			return nil
		}
		currentState := state[context.id]
		if result {
			updateAccumulator(
				&currentState.echos,
				&currentState.echoMessagesCount,
				&currentState.echoMessagesAccumulator,
				string(context.rootHash),
				context.chunk,
				context.fromId,
				params,
			)
			return nil
		} else {
			return fmt.Errorf("received invalid echo")
		}
	})

	brbpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, rootHash []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id)

		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}

		stringRepl := string(rootHash)

		currentState := state[id]

		if _, ok := currentState.readyMessagesReceived[stringRepl]; !ok {
			currentState.readyMessagesReceived[stringRepl] = make([]bool, params.GetN())
		}
		if !currentState.readyMessagesReceived[stringRepl][fromId] {
			currentState.readyMessagesReceived[stringRepl][fromId] = true

			currentState.readyMessagesCount[stringRepl]++
			if currentState.readyMessagesCount[stringRepl] > currentState.readyMessagesAccumulator.count {
				currentState.readyMessagesAccumulator.count = currentState.readyMessagesCount[stringRepl]
				currentState.readyMessagesAccumulator.value = stringRepl
			}
		}

		return nil
	})

	dsl.UponCondition(m, func() error {
		for id, currentState := range state {
			if id <= lastId {
				continue
			}
			if (currentState.echoMessagesAccumulator.count >= params.GetN()-params.GetF()) && !currentState.sentReady {
				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbmsgs.ReadyMessage(mc.Self, id, []byte(currentState.echoMessagesAccumulator.value)), params.AllNodes)
			}

			if (currentState.readyMessagesAccumulator.count >= params.GetN()-2*params.GetF()) && !currentState.sentReady {
				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbmsgs.ReadyMessage(mc.Self, id, []byte(currentState.readyMessagesAccumulator.value)), params.AllNodes)
			}

			if currentState.readyMessagesAccumulator.count >= params.GetN()-params.GetF() && currentState.echoMessagesAccumulator.count > params.GetN()-2*params.GetF() && !currentState.delivered {
				currentState.delivered = true
				decoded := make([][]byte, params.GetN())
				for i := 0; i < params.GetN(); i++ {
					decoded[i] = currentState.echos[currentState.echoMessagesAccumulator.value][i]
				}

				err := encoder.Reconstruct(decoded)
				if err != nil {
					return err
				}

				output := bytes.Join(decoded, []byte{})

				size := binary.LittleEndian.Uint32(output[:4])
				output = output[4 : 4+size]
				if id > lastId {
					lastId = id
				}
				brbpbdsl.Deliver(m, mc.Consumer, id, output)
				delete(state, id)
			}
		}
		return nil
	})

	return m, nil
}

func initialize(state *map[int64]*brbHashModuleState, id int64) {
	if _, ok := (*state)[id]; !ok {
		(*state)[id] = &brbHashModuleState{
			proposed: nil,

			sentEcho:                 false,
			sentReady:                false,
			delivered:                false,
			echos:                    make(map[string][][]byte),
			echoMessagesCount:        make(map[string]int),
			echoMessagesAccumulator:  Accumulator{value: "", count: -1},
			readyMessagesCount:       make(map[string]int),
			readyMessagesReceived:    make(map[string][]bool),
			readyMessagesAccumulator: Accumulator{value: "", count: -1},
		}
	}
}

func updateAccumulator(values *map[string][][]byte, counter *map[string]int, accumulator *Accumulator, hash string, chunk []byte, id int64, params *ModuleParams) {
	if _, ok := (*values)[hash]; !ok {
		(*values)[hash] = make([][]byte, params.GetN())
	}
	if (*values)[hash][id] == nil {
		(*values)[hash][id] = chunk
		(*counter)[hash]++
		if (*counter)[hash] > accumulator.count {
			(*accumulator).count = (*counter)[hash]
			(*accumulator).value = hash
		}
	}
}

type hashInitialMessageContext struct {
	id   int64
	data [][]byte
}

type hashStartMessageContext struct {
	id       int64
	chunk    []byte
	rootHash []byte
	proof    *commonpb.MerklePath
}

type hashEchoMessageContext struct {
	id       int64
	fromId   int64
	chunk    []byte
	rootHash []byte
	proof    *commonpb.MerklePath
}
