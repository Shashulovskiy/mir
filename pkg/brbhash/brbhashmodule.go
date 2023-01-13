package brbhash

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	brbpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"hash"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self     t.ModuleID // id of this module
	Consumer t.ModuleID // id of the module to send the "Deliver" event to
	Net      t.ModuleID
	Crypto   t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:     "brbhash",
		Consumer: consumer,
		Net:      "net",
		Crypto:   "crypto",
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

	echos               map[t.NodeID][]byte
	echoMessagesCount   map[string]int
	echosMaxAccumulator Accumulator

	readys              map[t.NodeID][]byte
	readyMessagesCount  map[string]int
	readyMaxAccumulator Accumulator

	hasher hash.Hash
}

func incrementAndUpdateAccumulator(message *[]byte, counts *map[string]int, accumulator *Accumulator) {
	// TODO cant store []byte in a map
	var stringRepr = string(*message)

	(*counts)[stringRepr]++

	if (*counts)[stringRepr] > (*accumulator).count {
		(*accumulator).count = (*counts)[stringRepr]
		(*accumulator).value = *message
	}
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	// upon event <brb, Init> do
	state := brbHashModuleState{
		proposed: nil,

		sentEcho:            false,
		sentReady:           false,
		delivered:           false,
		echos:               make(map[t.NodeID][]byte),
		echoMessagesCount:   make(map[string]int),
		echosMaxAccumulator: Accumulator{},
		readys:              make(map[t.NodeID][]byte),
		readyMessagesCount:  make(map[string]int),
		readyMaxAccumulator: Accumulator{},
		hasher:              sha1.New(),
	}

	// upon event <brb, Broadcast | m> do
	brbpbdsl.UponBroadcastRequest(m, func(data []byte) error {
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		state.proposed = data
		eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.StartMessage(mc.Self, data), params.AllNodes)
		return nil
	})

	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, data []byte) error {
		if from == params.Leader && state.sentEcho == false {
			state.proposed = data
			state.sentEcho = true
			state.hasher.Write(data)
			eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.EchoMessage(mc.Self, state.hasher.Sum(nil)), params.AllNodes)
			state.hasher.Reset()
			return nil
		} else {
			return fmt.Errorf("received start message not from leader or already sent echo")
		}
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, hash []byte) error {
		if state.echos[from] == nil {
			state.echos[from] = hash

			incrementAndUpdateAccumulator(&hash, &state.echoMessagesCount, &state.echosMaxAccumulator)
		}
		return nil
	})

	brbpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, hash []byte) error {
		if state.readys[from] == nil {
			state.readys[from] = hash

			incrementAndUpdateAccumulator(&hash, &state.readyMessagesCount, &state.readyMaxAccumulator)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		if (state.echosMaxAccumulator.count > (params.GetN()+params.GetF())/2 || state.readyMaxAccumulator.count > params.GetF()) && state.sentReady == false {
			state.sentReady = true
			eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.ReadyMessage(mc.Self, state.echosMaxAccumulator.value), params.AllNodes)
		}

		if state.readyMaxAccumulator.count > 2*params.GetF() && state.delivered == false {
			state.hasher.Write(state.proposed)
			var proposedHash = state.hasher.Sum(nil)
			state.hasher.Reset()
			if bytes.Equal(proposedHash, state.readyMaxAccumulator.value) {
				state.delivered = true
				brbpbdsl.Deliver(m, mc.Consumer, state.proposed)
			} else {
				panic("?")
			}
		}
		return nil
	})

	return m
}
