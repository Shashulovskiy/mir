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

// GetN returns the total number of nodes.
func (params *ModuleParams) GetN() int {
	return len(params.AllNodes)
}

// GetF returns the maximum tolerated number of faulty nodes.
func (params *ModuleParams) GetF() int {
	return (params.GetN() - 1) / 3
}

type MessageContent []byte

// brbModuleState represents the state of the brb module.
type brbModuleState struct {
	// this variable is not part of the original protocol description, but it greatly simplifies the code
	request []byte

	sentEcho  bool
	sentReady bool
	delivered bool

	echos               map[t.NodeID][]byte
	echoMessagesCount   map[string]int
	echosMaxAccumulator struct {
		MessageContent
		int
	}

	readys              map[t.NodeID][]byte
	readyMessagesCount  map[string]int
	readyMaxAccumulator struct {
		MessageContent
		int
	}
}

func incrementAndUpdateAccumulator(message *[]byte, counts *map[string]int, accumulator *struct {
	MessageContent
	int
}) {
	// TODO cant store []byte in a map
	var stringRepr = string(*message)

	(*counts)[stringRepr]++

	if (*counts)[stringRepr] > (*accumulator).int {
		(*accumulator).int = (*counts)[stringRepr]
		(*accumulator).MessageContent = *message
	}
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	// upon event <brb, Init> do
	state := make(map[int64]*brbModuleState)
	lastId := int64(0)

	// upon event <brb, Broadcast | m> do
	brbpbdsl.UponBroadcastRequest(m, func(id int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		if nodeID != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		initialize(&state, id)
		state[id].request = data
		eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.StartMessage(mc.Self, id, data), params.AllNodes)
		return nil
	})

	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		if from == params.Leader {
			initialize(&state, id)
			if state[id].sentEcho == false {
				state[id].sentEcho = true
				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.EchoMessage(mc.Self, id, data), params.AllNodes)
				return nil
			} else {
				return fmt.Errorf("already sent echo")
			}
		} else {
			return fmt.Errorf("received start message not from leader")
		}
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id)
		if state[id].echos[from] == nil {
			state[id].echos[from] = data

			incrementAndUpdateAccumulator(&data, &state[id].echoMessagesCount, &state[id].echosMaxAccumulator)
		}
		return nil
	})

	brbpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		if id <= lastId {
			return nil
		}
		initialize(&state, id)
		if state[id].readys[from] == nil {
			state[id].readys[from] = data

			incrementAndUpdateAccumulator(&data, &state[id].readyMessagesCount, &state[id].readyMaxAccumulator)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		for id, currentState := range state {
			if id <= lastId {
				continue
			}
			if (currentState.echosMaxAccumulator.int > (params.GetN()+params.GetF())/2 || currentState.readyMaxAccumulator.int > params.GetF()) && currentState.sentReady == false {
				currentState.sentReady = true
				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.ReadyMessage(mc.Self, id, currentState.echosMaxAccumulator.MessageContent), params.AllNodes)
			}

			if currentState.readyMaxAccumulator.int > 2*params.GetF() && currentState.delivered == false {
				currentState.delivered = true
				if id > lastId {
					lastId = id
				}
				brbpbdsl.Deliver(m, mc.Consumer, id, currentState.readyMaxAccumulator.MessageContent)
				delete(state, id)
			}
		}
		return nil
	})

	return m
}

func initialize(state *map[int64]*brbModuleState, id int64) {
	if _, ok := (*state)[id]; !ok {
		(*state)[id] = &brbModuleState{
			request: nil,

			sentEcho:          false,
			sentReady:         false,
			delivered:         false,
			echos:             make(map[t.NodeID][]byte),
			echoMessagesCount: make(map[string]int),
			echosMaxAccumulator: struct {
				MessageContent
				int
			}{MessageContent: nil, int: 0},
			readys:             make(map[t.NodeID][]byte),
			readyMessagesCount: make(map[string]int),
			readyMaxAccumulator: struct {
				MessageContent
				int
			}{MessageContent: nil, int: 0},
		}
	}
}
