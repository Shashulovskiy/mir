package merkletreepbevents

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/merkletreepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func VerifyRequest(destModule types.ModuleID, chunk []uint8, rootHash []uint8, proof [][]uint8) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Merkletree{
			Merkletree: &types2.Event{
				Type: &types2.Event_VerifyRequest{
					VerifyRequest: &types2.VerifyRequest{
						Chunk:    chunk,
						RootHash: rootHash,
						Proof:    proof,
					},
				},
			},
		},
	}
}

func VerifyResult(destModule types.ModuleID, success bool) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Merkletree{
			Merkletree: &types2.Event{
				Type: &types2.Event_VerifyResult{
					VerifyResult: &types2.VerifyResult{
						Success: success,
					},
				},
			},
		},
	}
}
