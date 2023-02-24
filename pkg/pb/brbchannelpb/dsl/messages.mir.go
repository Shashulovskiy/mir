package brbchannelpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/brbchannelpb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Brbchannel](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponAnyMessageReceived(m dsl.Module, handler func(from types1.NodeID, id int64, message *anypb.Any) error) {
	UponMessageReceived[*types.Message_AnyMessage](m, func(from types1.NodeID, msg *types.AnyMessage) error {
		return handler(from, msg.Id, msg.Message)
	})
}
