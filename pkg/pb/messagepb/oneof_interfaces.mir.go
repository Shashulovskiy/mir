package messagepb

import (
	mscpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	brbctpb "github.com/filecoin-project/mir/pkg/pb/brbctpb"
	brbdxrpb "github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	brbpb "github.com/filecoin-project/mir/pkg/pb/brbpb"
	checkpointpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	ordererspb "github.com/filecoin-project/mir/pkg/pb/ordererspb"
	pingpongpb "github.com/filecoin-project/mir/pkg/pb/pingpongpb"
)

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Iss) Unwrap() *isspb.ISSMessage {
	return w.Iss
}

func (w *Message_Bcb) Unwrap() *bcbpb.Message {
	return w.Bcb
}

func (w *Message_MultisigCollector) Unwrap() *mscpb.Message {
	return w.MultisigCollector
}

func (w *Message_Pingpong) Unwrap() *pingpongpb.Message {
	return w.Pingpong
}

func (w *Message_Checkpoint) Unwrap() *checkpointpb.Message {
	return w.Checkpoint
}

func (w *Message_SbMessage) Unwrap() *ordererspb.SBInstanceMessage {
	return w.SbMessage
}

func (w *Message_Brb) Unwrap() *brbpb.Message {
	return w.Brb
}

func (w *Message_Brbdxr) Unwrap() *brbdxrpb.Message {
	return w.Brbdxr
}

func (w *Message_Brbct) Unwrap() *brbctpb.Message {
	return w.Brbct
}
