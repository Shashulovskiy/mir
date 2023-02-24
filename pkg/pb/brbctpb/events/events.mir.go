package brbctpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/brbctpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func BroadcastRequest(destModule types.ModuleID, id int64, data []uint8) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Brbct{
			Brbct: &types2.Event{
				Type: &types2.Event_Request{
					Request: &types2.BroadcastRequest{
						Id:   id,
						Data: data,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, id int64, data []uint8) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Brbct{
			Brbct: &types2.Event{
				Type: &types2.Event_Deliver{
					Deliver: &types2.Deliver{
						Id:   id,
						Data: data,
					},
				},
			},
		},
	}
}
