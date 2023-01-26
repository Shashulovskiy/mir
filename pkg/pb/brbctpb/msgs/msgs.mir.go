package brbctpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/brbctpb/types"
	commonpb "github.com/filecoin-project/mir/pkg/pb/commonpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func StartMessage(destModule types.ModuleID, chunk []uint8, rootHash []uint8, proof *commonpb.MerklePath) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbct{
			Brbct: &types2.Message{
				Type: &types2.Message_StartMessage{
					StartMessage: &types2.StartMessage{
						Chunk:    chunk,
						RootHash: rootHash,
						Proof:    proof,
					},
				},
			},
		},
	}
}

func EchoMessage(destModule types.ModuleID, chunk []uint8, rootHash []uint8, proof *commonpb.MerklePath) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbct{
			Brbct: &types2.Message{
				Type: &types2.Message_EchoMessage{
					EchoMessage: &types2.EchoMessage{
						Chunk:    chunk,
						RootHash: rootHash,
						Proof:    proof,
					},
				},
			},
		},
	}
}

func ReadyMessage(destModule types.ModuleID, rootHash []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbct{
			Brbct: &types2.Message{
				Type: &types2.Message_ReadyMessage{
					ReadyMessage: &types2.ReadyMessage{
						RootHash: rootHash,
					},
				},
			},
		},
	}
}
