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

func NewByzantineModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, strategy string) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	// -----------
	brbctpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				eventpbdsl.SendMessage(m, mc.Net, brbmsgs.EchoMessage(mc.Self, id, brb.Corrupt(chunk), rootHash, proof), params.AllNodes)
			}
		}
		return nil
	})

	brbctpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, chunk []byte, rootHash []byte, proof *commonpb.MerklePath) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				eventpbdsl.SendMessage(m, mc.Net, brbmsgs.ReadyMessage(mc.Self, id, brb.Corrupt(rootHash)), params.AllNodes)
			}
		}
		return nil
	})

	brbctpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, rootHash []byte) error {
		return nil
	})

	return m, nil
}
