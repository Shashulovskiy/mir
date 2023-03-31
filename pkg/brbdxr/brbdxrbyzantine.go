package brbdxr

import (
	"github.com/filecoin-project/mir/pkg/brb"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbdxrpbdsl "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/dsl"
	brbdxrpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/mathutil"
)

type byzantineModuleState struct {
	sentEcho  bool
	sentReady bool
}

func NewByzantineModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, strategy string) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	lastId := int64(-1)
	states := make(map[int64]*byzantineModuleState)

	brbdxrpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, hdata []byte) error {
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

					nDataShards := params.GetN() - 2*params.GetF()

					data := make([]byte, mathutil.Pad(4+len(hdata), nDataShards))
					shardSize := len(data) / nDataShards

					eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.EchoMessage(mc.Self, id, brb.Corrupt(hdata[:shardSize]), brb.Corrupt(hdata[:shardSize])), params.AllNodes)
				}
			}
		}
		return nil
	})

	brbdxrpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, hash, chunk []byte) error {
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
					eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(mc.Self, id, hash, brb.Corrupt(chunk)), params.AllNodes)
				}
			}
		}
		return nil
	})

	brbdxrpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, hash, chunk []byte) error {
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
