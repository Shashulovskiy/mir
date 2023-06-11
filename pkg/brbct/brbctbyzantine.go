package brbct

import (
	"github.com/filecoin-project/mir/pkg/brb"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbctpbdsl "github.com/filecoin-project/mir/pkg/pb/brbctpb/dsl"
	brbmsgs "github.com/filecoin-project/mir/pkg/pb/brbctpb/msgs"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type byzantineModuleState struct {
	sentEcho  bool
	sentReady bool
}

func NewByzantineModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, strategy string) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	lastId := int64(-1)
	states := make(map[int64]*byzantineModuleState)

	// -----------
	brbctpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
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
					eventpbdsl.SendMessage(m, mc.Net, brbmsgs.EchoMessage(mc.Self, id, brb.Corrupt(chunk), rootHash, proof), params.AllNodes)
				}
			}
		}
		return nil
	})

	brbctpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
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
					eventpbdsl.SendMessage(m, mc.Net, brbmsgs.ReadyMessage(mc.Self, id, brb.Corrupt(rootHash)), params.AllNodes)
				}
			}
		}
		return nil
	})

	brbctpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, rootHash []byte) error {
		return nil
	})

	return m, nil
}

func initializeByzantine(id int64, states map[int64]*byzantineModuleState) {
	if _, ok := states[id]; !ok {
		states[id] = &byzantineModuleState{
			sentEcho:  false,
			sentReady: false,
		}
	}
}
