package brbdxrpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/brbdxrpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func StartMessage(destModule types.ModuleID, id int64, n int64, data []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbdxr{
			Brbdxr: &types2.Message{
				Type: &types2.Message_StartMessage{
					StartMessage: &types2.StartMessage{
						Id:   id,
						N:    n,
						Data: data,
					},
				},
			},
		},
	}
}

func EchoMessage(destModule types.ModuleID, id int64, n int64, hash []uint8, chunk []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbdxr{
			Brbdxr: &types2.Message{
				Type: &types2.Message_EchoMessage{
					EchoMessage: &types2.EchoMessage{
						Id:    id,
						N:     n,
						Hash:  hash,
						Chunk: chunk,
					},
				},
			},
		},
	}
}

func ReadyMessage(destModule types.ModuleID, id int64, n int64, hash []uint8, chunk []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbdxr{
			Brbdxr: &types2.Message{
				Type: &types2.Message_ReadyMessage{
					ReadyMessage: &types2.ReadyMessage{
						Id:    id,
						N:     n,
						Hash:  hash,
						Chunk: chunk,
					},
				},
			},
		},
	}
}
