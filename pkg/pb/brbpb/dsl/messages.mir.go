package brbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/brbpb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Brb](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponStartMessageReceived(m dsl.Module, handler func(from types1.NodeID, id int64, data []uint8) error) {
	UponMessageReceived[*types.Message_StartMessage](m, func(from types1.NodeID, msg *types.StartMessage) error {
		return handler(from, msg.Id, msg.Data)
	})
}

func UponEchoMessageReceived(m dsl.Module, handler func(from types1.NodeID, id int64, data []uint8) error) {
	UponMessageReceived[*types.Message_EchoMessage](m, func(from types1.NodeID, msg *types.EchoMessage) error {
		return handler(from, msg.Id, msg.Data)
	})
}

func UponReadyMessageReceived(m dsl.Module, handler func(from types1.NodeID, id int64, data []uint8) error) {
	UponMessageReceived[*types.Message_ReadyMessage](m, func(from types1.NodeID, msg *types.ReadyMessage) error {
		return handler(from, msg.Id, msg.Data)
	})
}
