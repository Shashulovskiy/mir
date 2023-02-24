package brbchannelpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbchannelpb "github.com/filecoin-project/mir/pkg/pb/brbchannelpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() brbchannelpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb brbchannelpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *brbchannelpb.Event_Request:
		return &Event_Request{Request: BroadcastRequestFromPb(pb.Request)}
	case *brbchannelpb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	}
	return nil
}

type Event_Request struct {
	Request *BroadcastRequest
}

func (*Event_Request) isEvent_Type() {}

func (w *Event_Request) Unwrap() *BroadcastRequest {
	return w.Request
}

func (w *Event_Request) Pb() brbchannelpb.Event_Type {
	return &brbchannelpb.Event_Request{Request: (w.Request).Pb()}
}

func (*Event_Request) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.Event_Request]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() brbchannelpb.Event_Type {
	return &brbchannelpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.Event_Deliver]()}
}

func EventFromPb(pb *brbchannelpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *brbchannelpb.Event {
	return &brbchannelpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.Event]()}
}

type BroadcastRequest struct {
	Id   int64
	Data []uint8
}

func BroadcastRequestFromPb(pb *brbchannelpb.BroadcastRequest) *BroadcastRequest {
	return &BroadcastRequest{
		Id:   pb.Id,
		Data: pb.Data,
	}
}

func (m *BroadcastRequest) Pb() *brbchannelpb.BroadcastRequest {
	return &brbchannelpb.BroadcastRequest{
		Id:   m.Id,
		Data: m.Data,
	}
}

func (*BroadcastRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.BroadcastRequest]()}
}

type Deliver struct {
	Id   int64
	Data []uint8
}

func DeliverFromPb(pb *brbchannelpb.Deliver) *Deliver {
	return &Deliver{
		Id:   pb.Id,
		Data: pb.Data,
	}
}

func (m *Deliver) Pb() *brbchannelpb.Deliver {
	return &brbchannelpb.Deliver{
		Id:   m.Id,
		Data: m.Data,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() brbchannelpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb brbchannelpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *brbchannelpb.Message_AnyMessage:
		return &Message_AnyMessage{AnyMessage: AnyMessageFromPb(pb.AnyMessage)}
	}
	return nil
}

type Message_AnyMessage struct {
	AnyMessage *AnyMessage
}

func (*Message_AnyMessage) isMessage_Type() {}

func (w *Message_AnyMessage) Unwrap() *AnyMessage {
	return w.AnyMessage
}

func (w *Message_AnyMessage) Pb() brbchannelpb.Message_Type {
	return &brbchannelpb.Message_AnyMessage{AnyMessage: (w.AnyMessage).Pb()}
}

func (*Message_AnyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.Message_AnyMessage]()}
}

func MessageFromPb(pb *brbchannelpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *brbchannelpb.Message {
	return &brbchannelpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.Message]()}
}

type AnyMessage struct {
	Id      int64
	Message *anypb.Any
}

func AnyMessageFromPb(pb *brbchannelpb.AnyMessage) *AnyMessage {
	return &AnyMessage{
		Id:      pb.Id,
		Message: pb.Message,
	}
}

func (m *AnyMessage) Pb() *brbchannelpb.AnyMessage {
	return &brbchannelpb.AnyMessage{
		Id:      m.Id,
		Message: m.Message,
	}
}

func (*AnyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbchannelpb.AnyMessage]()}
}
