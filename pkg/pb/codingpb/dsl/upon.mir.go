package codingpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	codingpb "github.com/filecoin-project/mir/pkg/pb/codingpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/codingpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// Module-specific dsl functions for processing events.

func UponEncodeRequest(m dsl.Module, handler func(totalShards int64, dataShards int64, paddedData []uint8, origin *codingpb.EncodeOrigin) error) {
	dsl.UponMirEvent[*types.Event_EncodeRequest](m, func(ev *types1.EncodeRequest) error {
		return handler(ev.TotalShards, ev.DataShards, ev.PaddedData, ev.Origin)
	})
}

func UponEncodeResult(m dsl.Module, handler func(encoded [][]uint8, origin *codingpb.EncodeOrigin) error) {
	dsl.UponMirEvent[*types.Event_EncodeResult](m, func(ev *types1.EncodeResult) error {
		return handler(ev.Encoded, ev.Origin)
	})
}

func UponDecodeRequest(m dsl.Module, handler func(totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.DecodeOrigin) error) {
	dsl.UponMirEvent[*types.Event_DecodeRequest](m, func(ev *types1.DecodeRequest) error {
		return handler(ev.TotalShards, ev.DataShards, ev.Shares, ev.Origin)
	})
}

func UponDecodeResult(m dsl.Module, handler func(success bool, decoded []uint8, origin *codingpb.DecodeOrigin) error) {
	dsl.UponMirEvent[*types.Event_DecodeResult](m, func(ev *types1.DecodeResult) error {
		return handler(ev.Success, ev.Decoded, ev.Origin)
	})
}

func UponRebuildRequest(m dsl.Module, handler func(totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.RebuildOrigin) error) {
	dsl.UponMirEvent[*types.Event_RebuildRequest](m, func(ev *types1.RebuildRequest) error {
		return handler(ev.TotalShards, ev.DataShards, ev.Shares, ev.Origin)
	})
}

func UponRebuildResult(m dsl.Module, handler func(success bool, decoded []uint8, origin *codingpb.RebuildOrigin) error) {
	dsl.UponMirEvent[*types.Event_RebuildResult](m, func(ev *types1.RebuildResult) error {
		return handler(ev.Success, ev.Decoded, ev.Origin)
	})
}
