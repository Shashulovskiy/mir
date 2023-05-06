package brbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/brbpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Brb](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponBroadcastRequest(m dsl.Module, handler func(id int64, n int64, data []uint8) error) {
	UponEvent[*types.Event_Request](m, func(ev *types.BroadcastRequest) error {
		return handler(ev.Id, ev.N, ev.Data)
	})
}

func UponDeliver(m dsl.Module, handler func(id int64, data []uint8) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Id, ev.Data)
	})
}
