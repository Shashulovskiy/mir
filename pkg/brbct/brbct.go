package brbct

import (
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbctpbdsl "github.com/filecoin-project/mir/pkg/pb/brbctpb/dsl"
	brbmsgs "github.com/filecoin-project/mir/pkg/pb/brbctpb/msgs"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/codingpb"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/mathutil"
	rs "github.com/vivint/infectious"
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
	Coder               t.ModuleID
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
		Coder:               "coder",
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

	echos                    map[string][]rs.Share
	receivedEcho             []bool
	echoMessagesCount        map[string]int
	echoMessagesAccumulator  Accumulator
	readyMessagesReceived    map[string][]bool
	readyMessagesCount       map[string]int
	readyMessagesAccumulator Accumulator
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

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

		dsl.EncodeRequest(m, mc.Coder, int64(params.GetN()), int64(nDataShards), data, &encodeDataContext{id: id})

		return nil
	})

	dsl.UponEncodeResult(m, func(encoded [][]byte, context *encodeDataContext) error {
		if context.id <= lastId {
			return nil
		}

		dsl.MerkleBuildRequest(m, mc.MerkleProofVerifier, encoded, &hashInitialMessageContext{id: context.id, data: encoded})

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
				currentState.echos[hash] = append(currentState.echos[hash], rs.Share{
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
				shares := make([]*codingpb.Share, 0)
				for _, rd := range currentState.echos[currentState.echoMessagesAccumulator.value] {
					shares = append(shares, &codingpb.Share{
						Number: int64(rd.Number),
						Chunk:  rd.Data,
					})
				}

				dsl.RebuildRequest(m, mc.Coder, int64(params.GetN()), int64(params.GetN()-2*params.GetF()), shares, &rebuildDataContext{id: id})
			}
		}
		return nil
	})

	dsl.UponRebuildResult(m, func(success bool, decoded []byte, context *rebuildDataContext) error {
		if context.id <= lastId {
			return nil
		}
		if !success {
			panic("Unable to rebuild?")
		}

		size := binary.LittleEndian.Uint32(decoded[:4])
		state[context.id].delivered = true
		decoded = decoded[4 : 4+size]
		brbpbdsl.Deliver(m, mc.Consumer, context.id, decoded)
		if context.id > lastId {
			for i := lastId; i <= context.id; i++ {
				delete(state, i)
			}
			lastId = context.id
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
			echos:                    make(map[string][]rs.Share),
			receivedEcho:             make([]bool, n),
			echoMessagesCount:        make(map[string]int),
			echoMessagesAccumulator:  Accumulator{value: "", count: -1},
			readyMessagesCount:       make(map[string]int),
			readyMessagesReceived:    make(map[string][]bool),
			readyMessagesAccumulator: Accumulator{value: "", count: -1},
		}
	}
}

type encodeDataContext struct {
	id int64
}

type rebuildDataContext struct {
	id int64
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
