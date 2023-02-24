package brbchannelpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/brbchannelpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

func AnyMessage(destModule types.ModuleID, id int64, message *anypb.Any) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Brbchannel{
			Brbchannel: &types2.Message{
				Type: &types2.Message_AnyMessage{
					AnyMessage: &types2.AnyMessage{
						Id:      id,
						Message: message,
					},
				},
			},
		},
	}
}
