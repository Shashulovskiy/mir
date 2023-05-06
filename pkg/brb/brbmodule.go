package brb

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	brbpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
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
		Self:     "brb",
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

type Accumulator struct {
	data  []byte
	count int
}

// brbModuleState represents the state of the brb module.
type state struct {
	n int64

	sentEcho  bool
	sentReady bool
	delivered bool

	echos               map[t.NodeID]bool
	echoMessagesCount   map[string]int
	echosMaxAccumulator Accumulator

	readys              map[t.NodeID]bool
	readyMessagesCount  map[string]int
	readyMaxAccumulator Accumulator
}

type brbModuleState state

func incrementAndUpdateAccumulator(message []byte, counts map[string]int, accumulator *Accumulator) {
	var stringRepr = string(message)

	counts[stringRepr]++

	if counts[stringRepr] > (*accumulator).count {
		(*accumulator).count = counts[stringRepr]
		(*accumulator).data = message
	}
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := make(map[int64]*brbModuleState)
	lastId := int64(0)

	brbpbdsl.UponBroadcastRequest(m, func(id, n int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		initialize(state, id, n)
		eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.StartMessage(mc.Self, id, n, data), params.AllNodes[:n])
		return nil
	})

	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id, n int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		if from == params.Leader {
			initialize(state, id, n)
			if state[id].sentEcho == false {
				state[id].sentEcho = true
				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.EchoMessage(mc.Self, id, n, data), params.AllNodes[:n])
				return nil
			} else {
				return fmt.Errorf("already sent echo")
			}
		} else {
			return fmt.Errorf("received start message not from leader")
		}
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id, n int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(state, id, n)
		if _, ok := state[id].echos[from]; !ok {
			state[id].echos[from] = true

			incrementAndUpdateAccumulator(data, state[id].echoMessagesCount, &state[id].echosMaxAccumulator)
		}
		return nil
	})

	brbpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id, n int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(state, id, n)
		if _, ok := state[id].readys[from]; !ok {
			state[id].readys[from] = true

			incrementAndUpdateAccumulator(data, state[id].readyMessagesCount, &state[id].readyMaxAccumulator)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		for id, currentState := range state {
			if id <= lastId {
				continue
			}
			if (currentState.echosMaxAccumulator.count > int(currentState.n+getF(currentState.n))/2 || currentState.readyMaxAccumulator.count > int(getF(currentState.n))) && currentState.sentReady == false {
				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.ReadyMessage(mc.Self, id, currentState.n, currentState.echosMaxAccumulator.data), params.AllNodes[:currentState.n])
			}

			if currentState.readyMaxAccumulator.count > 2*int(getF(currentState.n)) && currentState.delivered == false {
				currentState.delivered = true
				brbpbdsl.Deliver(m, mc.Consumer, id, currentState.readyMaxAccumulator.data)
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

	return m
}

func getF(n int64) int64 {
	return (n - 1) / 3
}

func initialize(state map[int64]*brbModuleState, id int64, n int64) {
	if _, ok := state[id]; !ok {
		state[id] = &brbModuleState{
			n:                   n,
			sentEcho:            false,
			sentReady:           false,
			delivered:           false,
			echos:               make(map[t.NodeID]bool),
			echoMessagesCount:   make(map[string]int),
			echosMaxAccumulator: Accumulator{data: nil, count: 0},
			readys:              make(map[t.NodeID]bool),
			readyMessagesCount:  make(map[string]int),
			readyMaxAccumulator: Accumulator{data: nil, count: 0},
		}
	}
}
