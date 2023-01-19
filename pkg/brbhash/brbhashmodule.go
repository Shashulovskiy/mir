package brbhash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbctpb/dsl"
	brbmsgs "github.com/filecoin-project/mir/pkg/pb/brbctpb/msgs"
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
		Self:     "brbhash",
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
	value []byte
	count int
}

// brbHashModuleState represents the state of the brb module.
type brbHashModuleState struct {
	// this variable is not part of the original protocol description, but it greatly simplifies the code
	proposed []byte

	sentEcho  bool
	sentReady bool
	delivered bool

	echos              [][]byte
	echoMessagesCount  int
	readyMessagesCount int
	startReceived      bool
	echoReceived       []bool
	readyReceived      []bool
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	encoder, err := rs.New(params.GetN()-params.GetF(), params.GetF())

	if err != nil {
		return nil, errors.Wrap(err, "Unable to create coder")
	}

	// upon event <brb, Init> do
	state := brbHashModuleState{
		proposed: nil,

		sentEcho:           false,
		sentReady:          false,
		delivered:          false,
		startReceived:      false,
		echos:              make([][]byte, params.GetN()),
		echoMessagesCount:  0,
		echoReceived:       make([]bool, params.GetN()),
		readyMessagesCount: 0,
		readyReceived:      make([]bool, params.GetN()),
	}

	// upon event <brb, Broadcast | m> do
	brbpbdsl.UponBroadcastRequest(m, func(hdata []byte) error {
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}

		data := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(data, uint32(len(hdata)))
		data = append(data, hdata...)

		chunkSize := int(math.Ceil(float64(len(data)) / float64(params.GetN()-params.GetF())))
		encoded := make([][]byte, params.GetN())
		for i := range encoded {
			encoded[i] = make([]byte, chunkSize)
		}

		for i, in := range encoded[:params.GetN()-params.GetF()] {
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

		toHash := make([][][]byte, params.GetN())
		for i := 0; i < params.GetN(); i++ {
			toHash[i] = [][]byte{encoded[i]}
		}

		dsl.HashRequest(m, mc.Hasher, toHash, &hashInitialMessageContext{data: encoded})

		return nil
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashInitialMessageContext) error {
		for i, node := range params.AllNodes {
			eventpbdsl.SendMessage(m, mc.Net, brbmsgs.StartMessage(mc.Self, context.data[i], hashes), []t.NodeID{node})
		}
		return nil
	})

	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, chunk []uint8, hashes [][]uint8) error {
		if from == params.Leader && !state.startReceived {
			dsl.HashOneMessage(m, mc.Hasher, [][]byte{chunk}, &hashStartMessageContext{chunk: chunk, hashes: hashes})

			return nil
		} else {
			return fmt.Errorf("received start message not from leader or already received start from leader")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashStartMessageContext) error {
		thisId, err := strconv.ParseInt(nodeID.Pb(), 10, 64)
		if err != nil {
			return err
		}
		if bytes.Equal(hashes[0], context.hashes[thisId]) && !state.startReceived {
			state.startReceived = true

			for _, node := range params.AllNodes {
				eventpbdsl.SendMessage(m, mc.Net, brbmsgs.EchoMessage(mc.Self, context.chunk, context.hashes), []t.NodeID{node})
			}

			return nil
		} else {
			return fmt.Errorf("received invalid chunk or already sent echo")
		}
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, chunk []uint8, hashes [][]uint8) error {
		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}

		if state.echoReceived[fromId] == false {
			dsl.HashOneMessage(m, mc.Hasher, [][]byte{chunk}, &hashEchoMessageContext{fromId: fromId, chunk: chunk, hashes: hashes})

			return nil
		} else {
			return fmt.Errorf("received duplicate echo")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashEchoMessageContext) error {
		if bytes.Equal(hashes[0], context.hashes[context.fromId]) {
			state.echos[context.fromId] = context.chunk
			state.echoReceived[context.fromId] = true
			state.echoMessagesCount += 1

			if state.echoMessagesCount == int(math.Max(math.Ceil(float64((params.GetN()+params.GetF()+1)/2)), float64(params.GetN()-2*params.GetF()))) && state.readyMessagesCount < params.GetN()-2*params.GetF() {
				decoded := make([][]byte, params.GetN())
				for i := 0; i < params.GetN(); i++ {
					decoded[i] = state.echos[i]
				}

				err := encoder.Reconstruct(decoded)
				if err != nil {
					return err
				}

				toHash := make([][][]byte, params.GetN())
				for i := 0; i < params.GetN(); i++ {
					toHash[i] = [][]byte{decoded[i]}
				}

				dsl.HashRequest(m, mc.Hasher, toHash, &hashEchoResultContext{decoded: decoded, hashes: context.hashes})
			}
			return nil
		} else {
			return fmt.Errorf("received invalid echo")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashEchoResultContext) error {
		thisId, err := strconv.ParseInt(nodeID.Pb(), 10, 64)
		if err != nil {
			return err
		}

		for i := 0; i < params.GetN(); i++ {
			if !bytes.Equal(hashes[i], context.hashes[i]) {
				return errors.New("invalid hash receive with echo")
			}
		}

		for _, node := range params.AllNodes {
			eventpbdsl.SendMessage(m, mc.Net, brbmsgs.ReadyMessage(mc.Self, context.decoded[thisId], context.hashes), []t.NodeID{node})
		}

		return nil
	})

	brbpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, chunk []uint8, hashes [][]uint8) error {
		fromId, err := strconv.ParseInt(from.Pb(), 10, 64)
		if err != nil {
			return err
		}

		if state.readyReceived[fromId] == false {
			dsl.HashOneMessage(m, mc.Hasher, [][]byte{chunk}, &hashReadyMessageContext{fromId: fromId, chunk: chunk, hashes: hashes})

			return nil
		} else {
			return fmt.Errorf("received duplicate echo")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashReadyMessageContext) error {
		if bytes.Equal(hashes[0], context.hashes[context.fromId]) {
			state.echos[context.fromId] = context.chunk
			state.readyReceived[context.fromId] = true
			state.readyMessagesCount += 1

			if state.echoMessagesCount < int(math.Max(math.Ceil(float64((params.GetN()+params.GetF()+1)/2)), float64(params.GetN()-2*params.GetF()))) && state.readyMessagesCount == params.GetN()-2*params.GetF() {
				decoded := make([][]byte, params.GetN())
				for i := 0; i < params.GetN(); i++ {
					decoded[i] = state.echos[i]
				}

				err := encoder.Reconstruct(decoded)
				if err != nil {
					return err
				}

				toHash := make([][][]byte, params.GetN())
				for i := 0; i < params.GetN(); i++ {
					toHash[i] = [][]byte{decoded[i]}
				}

				dsl.HashRequest(m, mc.Hasher, toHash, &hashReadyResultContext{decoded: decoded, hashes: context.hashes})
			} else if state.readyMessagesCount == params.GetN()-params.GetF() {
				decoded := make([][]byte, params.GetN())
				for i := 0; i < params.GetN(); i++ {
					decoded[i] = state.echos[i]
				}

				err := encoder.Reconstruct(decoded)
				if err != nil {
					return err
				}

				output := bytes.Join(decoded, []byte{})

				size := binary.LittleEndian.Uint32(output[:4])
				output = output[4 : 4+size]
				brbpbdsl.Deliver(m, mc.Consumer, output)
			}
			return nil
		} else {
			return fmt.Errorf("received invalid echo")
		}
	})

	dsl.UponHashResult(m, func(hashes [][]byte, context *hashReadyResultContext) error {
		thisId, err := strconv.ParseInt(nodeID.Pb(), 10, 64)
		if err != nil {
			return err
		}

		for i := 0; i < params.GetN(); i++ {
			if !bytes.Equal(hashes[i], context.hashes[i]) {
				return errors.New("invalid hash receive with ready")
			}
		}

		for _, node := range params.AllNodes {
			eventpbdsl.SendMessage(m, mc.Net, brbmsgs.ReadyMessage(mc.Self, context.decoded[thisId], context.hashes), []t.NodeID{node})
		}

		return nil
	})

	return m, nil
}

type hashInitialMessageContext struct {
	data [][]byte
}

type hashStartMessageContext struct {
	chunk  []byte
	hashes [][]byte
}

type hashEchoMessageContext struct {
	fromId int64
	chunk  []byte
	hashes [][]byte
}

type hashReadyMessageContext struct {
	fromId int64
	chunk  []byte
	hashes [][]byte
}

type hashEchoResultContext struct {
	decoded [][]byte
	hashes  [][]byte
}

type hashReadyResultContext struct {
	decoded [][]byte
	hashes  [][]byte
}
