package brbct

import (
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbctpbdsl "github.com/filecoin-project/mir/pkg/pb/brbctpb/dsl"
	brbmsgs "github.com/filecoin-project/mir/pkg/pb/brbctpb/msgs"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/mathutil"
	"github.com/pkg/errors"
	rs_ezpwd "github.com/vivint/infectious"
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

type Accumulator struct {
	value string
	count int
}

// brbHashModuleState represents the state of the brb module.
type brbHashModuleState struct {
	sentEcho  bool
	sentReady bool
	delivered bool

	echos                    map[string][]rs_ezpwd.Share
	receivedEcho             []bool
	echoMessagesCount        map[string]int
	echoMessagesAccumulator  Accumulator
	readyMessagesReceived    map[string][]bool
	readyMessagesCount       map[string]int
	readyMessagesAccumulator Accumulator
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	encoder, err := rs_ezpwd.NewFEC(params.GetN()-2*params.GetF(), params.GetN())

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

		nDataShards := params.GetN() - 2*params.GetF()

		data := make([]byte, mathutil.Pad(4+len(hdata), nDataShards))
		binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
		copy(data[4:], hdata)

		shardSize := len(data) / nDataShards

		encoded := make([][]byte, params.GetN())

		dataWithPadding := make([]byte, nDataShards*shardSize)
		copy(dataWithPadding, data)

		output := func(s rs_ezpwd.Share) {
			encoded[s.Number] = make([]byte, len(s.Data))
			copy(encoded[s.Number], s.Data)
		}

		err := encoder.Encode(dataWithPadding, output)

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
	brbctpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
		if id <= lastId {
			return nil
		}
		if from == params.Leader {
			initialize(state, id, params.GetN())

			if !state[id].sentEcho {

				dsl.MerkleProofVerifyRequest(m, mc.MerkleProofVerifier, rootHash, chunk, proof, &hashStartMessageContext{id: id, chunk: chunk, rootHash: rootHash, proof: proof})

				return nil
			} else {
				//fmt.Println("already received start from leader")
				return nil
			}
		}
		//fmt.Println("received start message not from leader")
		return nil
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
			//fmt.Println("received invalid chunk or already sent echo")
			return nil
		}
	})

	brbctpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
		if id <= lastId {
			return nil
		}
		initialize(state, id, params.GetN())

		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}

		dsl.MerkleProofVerifyRequest(m, mc.MerkleProofVerifier, rootHash, chunk, proof, &hashEchoMessageContext{id: id, fromId: fromId, chunk: chunk, rootHash: rootHash, proof: proof})
		return nil
	})

	dsl.MerkleProofVerifyResult(m, func(result bool, context *hashEchoMessageContext) error {
		if context.id <= lastId {
			return nil
		}
		currentState := state[context.id]
		if result {
			accumulator := &currentState.echoMessagesAccumulator
			hash := string(context.rootHash)
			if !currentState.receivedEcho[context.fromId] {
				currentState.receivedEcho[context.fromId] = true
				currentState.echos[hash] = append(currentState.echos[hash], rs_ezpwd.Share{
					Number: int(context.fromId),
					Data:   context.chunk,
				})
				currentState.echoMessagesCount[hash]++
				if currentState.echoMessagesCount[hash] > accumulator.count {
					(*accumulator).count = currentState.echoMessagesCount[hash]
					(*accumulator).value = hash
				}
			}
			return nil
		} else {
			//fmt.Println("received invalid echo")
			return nil
		}
	})

	brbctpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, rootHash []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(state, id, params.GetN())

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

			if currentState.readyMessagesAccumulator.count >= params.GetN()-params.GetF() && currentState.echoMessagesAccumulator.count >= params.GetN()-2*params.GetF() && !currentState.delivered {
				output := make([]byte, len(currentState.echos[currentState.echoMessagesAccumulator.value][0].Data)*(params.GetN()-2*params.GetF()))
				err := encoder.Rebuild(currentState.echos[currentState.echoMessagesAccumulator.value], func(s rs_ezpwd.Share) {
					copy(output[s.Number*len(s.Data):], s.Data)
				})
				if err != nil {
					return err
				}

				size := binary.LittleEndian.Uint32(output[:4])
				currentState.delivered = true
				output = output[4 : 4+size]
				brbpbdsl.Deliver(m, mc.Consumer, id, output)
				if id > lastId {
					for i := lastId; i <= id; i++ {
						delete(state, i)
					}
					lastId = id
				}
			}
		}
		return nil
	})

	return m, nil
}

func initialize(state map[int64]*brbHashModuleState, id int64, n int) {
	if _, ok := state[id]; !ok {
		state[id] = &brbHashModuleState{
			sentEcho:                 false,
			sentReady:                false,
			delivered:                false,
			echos:                    make(map[string][]rs_ezpwd.Share),
			receivedEcho:             make([]bool, n),
			echoMessagesCount:        make(map[string]int),
			echoMessagesAccumulator:  Accumulator{value: "", count: -1},
			readyMessagesCount:       make(map[string]int),
			readyMessagesReceived:    make(map[string][]bool),
			readyMessagesAccumulator: Accumulator{value: "", count: -1},
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
