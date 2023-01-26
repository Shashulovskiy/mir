package merkeltreepbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/merkeltreepb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func VerifyRequest(m dsl.Module, destModule types.ModuleID, chunk []uint8, rootHash []uint8, proof [][]uint8) {
	dsl.EmitMirEvent(m, events.VerifyRequest(destModule, chunk, rootHash, proof))
}

func VerifyResult(m dsl.Module, destModule types.ModuleID, success bool) {
	dsl.EmitMirEvent(m, events.VerifyResult(destModule, success))
}
