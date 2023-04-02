package codingpbevents

import (
	codingpb "github.com/filecoin-project/mir/pkg/pb/codingpb"
	types2 "github.com/filecoin-project/mir/pkg/pb/codingpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func EncodeRequest(destModule types.ModuleID, totalShards int64, dataShards int64, paddedData []uint8, origin *codingpb.EncodeOrigin) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_EncodeRequest{
			EncodeRequest: &types2.EncodeRequest{
				TotalShards: totalShards,
				DataShards:  dataShards,
				PaddedData:  paddedData,
				Origin:      origin,
			},
		},
	}
}

func EncodeResult(destModule types.ModuleID, encoded [][]uint8, origin *codingpb.EncodeOrigin) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_EncodeResult{
			EncodeResult: &types2.EncodeResult{
				Encoded: encoded,
				Origin:  origin,
			},
		},
	}
}

func DecodeRequest(destModule types.ModuleID, totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.DecodeOrigin) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_DecodeRequest{
			DecodeRequest: &types2.DecodeRequest{
				TotalShards: totalShards,
				DataShards:  dataShards,
				Shares:      shares,
				Origin:      origin,
			},
		},
	}
}

func DecodeResult(destModule types.ModuleID, success bool, decoded []uint8, origin *codingpb.DecodeOrigin) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_DecodeResult{
			DecodeResult: &types2.DecodeResult{
				Success: success,
				Decoded: decoded,
				Origin:  origin,
			},
		},
	}
}

func RebuildRequest(destModule types.ModuleID, totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.RebuildOrigin) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_RebuildRequest{
			RebuildRequest: &types2.RebuildRequest{
				TotalShards: totalShards,
				DataShards:  dataShards,
				Shares:      shares,
				Origin:      origin,
			},
		},
	}
}

func RebuildResult(destModule types.ModuleID, success bool, decoded []uint8, origin *codingpb.RebuildOrigin) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_RebuildResult{
			RebuildResult: &types2.RebuildResult{
				Success: success,
				Decoded: decoded,
				Origin:  origin,
			},
		},
	}
}
