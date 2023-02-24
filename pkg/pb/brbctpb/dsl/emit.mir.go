package brbctpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/brbctpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func BroadcastRequest(m dsl.Module, destModule types.ModuleID, id int64, data []uint8) {
	dsl.EmitMirEvent(m, events.BroadcastRequest(destModule, id, data))
}

func Deliver(m dsl.Module, destModule types.ModuleID, id int64, data []uint8) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, id, data))
}
