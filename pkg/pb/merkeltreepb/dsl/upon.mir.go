package merkletreepbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/merkletreepb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Merkletree](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponVerifyRequest(m dsl.Module, handler func(chunk []uint8, rootHash []uint8, proof [][]uint8) error) {
	UponEvent[*types.Event_VerifyRequest](m, func(ev *types.VerifyRequest) error {
		return handler(ev.Chunk, ev.RootHash, ev.Proof)
	})
}

func UponVerifyResult(m dsl.Module, handler func(success bool) error) {
	UponEvent[*types.Event_VerifyResult](m, func(ev *types.VerifyResult) error {
		return handler(ev.Success)
	})
}
