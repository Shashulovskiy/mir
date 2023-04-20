package brb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	brbpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type byzantineModuleState struct {
	sentEcho  bool
	sentReady bool
}

func NewByzantineModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, strategy string) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	lastId := int64(-1)
	states := make(map[int64]*byzantineModuleState)

	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		if id < lastId {
			return nil
		}
		if id > lastId {
			for i := id; i < lastId; i++ {
				delete(states, i)
			}
			lastId = id
		}
		initializeByzantine(id, states)

		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				if !states[id].sentEcho {
					state := states[id]
					state.sentEcho = true
					eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.EchoMessage(mc.Self, id, Corrupt(data)), params.AllNodes)
				}
			}
		}
		return nil
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		if id < lastId {
			return nil
		}
		if id > lastId {
			for i := id; i < lastId; i++ {
				delete(states, i)
			}
			lastId = id
		}
		initializeByzantine(id, states)

		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				if !states[id].sentReady {
					state := states[id]
					state.sentReady = true
					eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.ReadyMessage(mc.Self, id, Corrupt(data)), params.AllNodes)
				}
			}
		}
		return nil
	})

	brbpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		return nil
	})

	return m
}

func Corrupt(data []byte) []byte {
	corrupted := make([]byte, len(data))
	for i := range corrupted {
		if i == len(data)/2 {
			corrupted[i] = 1
		} else {
			corrupted[i] = data[i]
		}
	}

	return corrupted
}

func initializeByzantine(id int64, states map[int64]*byzantineModuleState) {
	if _, ok := states[id]; !ok {
		states[id] = &byzantineModuleState{
			sentEcho:  false,
			sentReady: false,
		}
	}
}
