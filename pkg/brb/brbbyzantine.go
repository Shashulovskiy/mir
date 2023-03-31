package brb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	brbpbdsl "github.com/filecoin-project/mir/pkg/pb/brbpb/dsl"
	brbpbmsgs "github.com/filecoin-project/mir/pkg/pb/brbpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

func NewByzantineModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, strategy string) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	brbpbdsl.UponBroadcastRequest(m, func(id int64, data []byte) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{

				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.StartMessage(mc.Self, id, Corrupt(data)), params.AllNodes)
			}
		}
		return nil
	})

	brbpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.EchoMessage(mc.Self, id, Corrupt(data)), params.AllNodes)
			}
		}
		return nil
	})

	brbpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, data []byte) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				eventpbdsl.SendMessage(m, mc.Net, brbpbmsgs.ReadyMessage(mc.Self, id, Corrupt(data)), params.AllNodes)
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
	for _, i := range corrupted {
		corrupted[i] = 1
	}

	return corrupted
}
