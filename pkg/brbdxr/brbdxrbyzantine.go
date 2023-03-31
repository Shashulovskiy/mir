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

func NewByzantineModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, strategy string) (modules.PassiveModule, error) {
	m := dsl.NewModule(mc.Self)

	brbdxrpbdsl.UponStartMessageReceived(m, func(from t.NodeID, id int64, hdata []byte) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				nDataShards := params.GetN() - 2*params.GetF()

				data := make([]byte, mathutil.Pad(4+len(hdata), nDataShards))
				shardSize := len(data) / nDataShards

				eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.EchoMessage(mc.Self, id, brb.Corrupt(hdata[:shardSize]), brb.Corrupt(hdata[:shardSize])), params.AllNodes)
			}
		}
		return nil
	})

	brbdxrpbdsl.UponEchoMessageReceived(m, func(from t.NodeID, id int64, hash, chunk []byte) error {
		switch strategy {
		case "ignore":
			{
				return nil
			}
		case "corrupt":
			{
				eventpbdsl.SendMessage(m, mc.Net, brbdxrpbmsgs.ReadyMessage(mc.Self, id, hash, brb.Corrupt(chunk)), params.AllNodes)
			}
		}
		return nil
	})

	brbdxrpbdsl.UponReadyMessageReceived(m, func(from t.NodeID, id int64, hash, chunk []byte) error {
		return nil
	})

	return m, nil
}
