package codingpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	codingpb "github.com/filecoin-project/mir/pkg/pb/codingpb"
	events "github.com/filecoin-project/mir/pkg/pb/codingpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func EncodeRequest(m dsl.Module, destModule types.ModuleID, totalShards int64, dataShards int64, paddedData []uint8, origin *codingpb.EncodeOrigin) {
	dsl.EmitMirEvent(m, events.EncodeRequest(destModule, totalShards, dataShards, paddedData, origin))
}

func EncodeResult(m dsl.Module, destModule types.ModuleID, encoded [][]uint8, origin *codingpb.EncodeOrigin) {
	dsl.EmitMirEvent(m, events.EncodeResult(destModule, encoded, origin))
}

func DecodeRequest(m dsl.Module, destModule types.ModuleID, totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.DecodeOrigin) {
	dsl.EmitMirEvent(m, events.DecodeRequest(destModule, totalShards, dataShards, shares, origin))
}

func DecodeResult(m dsl.Module, destModule types.ModuleID, success bool, decoded []uint8, origin *codingpb.DecodeOrigin) {
	dsl.EmitMirEvent(m, events.DecodeResult(destModule, success, decoded, origin))
}

func RebuildRequest(m dsl.Module, destModule types.ModuleID, totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.RebuildOrigin) {
	dsl.EmitMirEvent(m, events.RebuildRequest(destModule, totalShards, dataShards, shares, origin))
}

func RebuildResult(m dsl.Module, destModule types.ModuleID, success bool, decoded []uint8, origin *codingpb.RebuildOrigin) {
	dsl.EmitMirEvent(m, events.RebuildResult(destModule, success, decoded, origin))
}
