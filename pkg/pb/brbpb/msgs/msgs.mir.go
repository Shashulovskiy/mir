package brbpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/brbpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func StartMessage(destModule types.ModuleID, id int64, data []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brb{
			Brb: &types2.Message{
				Type: &types2.Message_StartMessage{
					StartMessage: &types2.StartMessage{
						Id:   id,
						Data: data,
					},
				},
			},
		},
	}
}

func EchoMessage(destModule types.ModuleID, id int64, data []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brb{
			Brb: &types2.Message{
				Type: &types2.Message_EchoMessage{
					EchoMessage: &types2.EchoMessage{
						Id:   id,
						Data: data,
					},
				},
			},
		},
	}
}

func ReadyMessage(destModule types.ModuleID, id int64, data []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brb{
			Brb: &types2.Message{
				Type: &types2.Message_ReadyMessage{
					ReadyMessage: &types2.ReadyMessage{
						Id:   id,
						Data: data,
					},
				},
			},
		},
	}
}
