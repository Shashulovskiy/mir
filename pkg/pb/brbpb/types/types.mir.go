package brbpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbpb "github.com/filecoin-project/mir/pkg/pb/brbpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() brbpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb brbpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *brbpb.Event_Request:
		return &Event_Request{Request: BroadcastRequestFromPb(pb.Request)}
	case *brbpb.Event_Deliver:
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

func (w *Event_Request) Pb() brbpb.Event_Type {
	return &brbpb.Event_Request{Request: (w.Request).Pb()}
}

func (*Event_Request) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Event_Request]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() brbpb.Event_Type {
	return &brbpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Event_Deliver]()}
}

func EventFromPb(pb *brbpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *brbpb.Event {
	return &brbpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Event]()}
}

type BroadcastRequest struct {
	Id   int64
	N    int64
	Data []uint8
}

func BroadcastRequestFromPb(pb *brbpb.BroadcastRequest) *BroadcastRequest {
	return &BroadcastRequest{
		Id:   pb.Id,
		N:    pb.N,
		Data: pb.Data,
	}
}

func (m *BroadcastRequest) Pb() *brbpb.BroadcastRequest {
	return &brbpb.BroadcastRequest{
		Id:   m.Id,
		N:    m.N,
		Data: m.Data,
	}
}

func (*BroadcastRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.BroadcastRequest]()}
}

type Deliver struct {
	Id   int64
	Data []uint8
}

func DeliverFromPb(pb *brbpb.Deliver) *Deliver {
	return &Deliver{
		Id:   pb.Id,
		Data: pb.Data,
	}
}

func (m *Deliver) Pb() *brbpb.Deliver {
	return &brbpb.Deliver{
		Id:   m.Id,
		Data: m.Data,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() brbpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb brbpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *brbpb.Message_StartMessage:
		return &Message_StartMessage{StartMessage: StartMessageFromPb(pb.StartMessage)}
	case *brbpb.Message_EchoMessage:
		return &Message_EchoMessage{EchoMessage: EchoMessageFromPb(pb.EchoMessage)}
	case *brbpb.Message_ReadyMessage:
		return &Message_ReadyMessage{ReadyMessage: ReadyMessageFromPb(pb.ReadyMessage)}
	}
	return nil
}

type Message_StartMessage struct {
	StartMessage *StartMessage
}

func (*Message_StartMessage) isMessage_Type() {}

func (w *Message_StartMessage) Unwrap() *StartMessage {
	return w.StartMessage
}

func (w *Message_StartMessage) Pb() brbpb.Message_Type {
	return &brbpb.Message_StartMessage{StartMessage: (w.StartMessage).Pb()}
}

func (*Message_StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Message_StartMessage]()}
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage
}

func (*Message_EchoMessage) isMessage_Type() {}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_EchoMessage) Pb() brbpb.Message_Type {
	return &brbpb.Message_EchoMessage{EchoMessage: (w.EchoMessage).Pb()}
}

func (*Message_EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Message_EchoMessage]()}
}

type Message_ReadyMessage struct {
	ReadyMessage *ReadyMessage
}

func (*Message_ReadyMessage) isMessage_Type() {}

func (w *Message_ReadyMessage) Unwrap() *ReadyMessage {
	return w.ReadyMessage
}

func (w *Message_ReadyMessage) Pb() brbpb.Message_Type {
	return &brbpb.Message_ReadyMessage{ReadyMessage: (w.ReadyMessage).Pb()}
}

func (*Message_ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Message_ReadyMessage]()}
}

func MessageFromPb(pb *brbpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *brbpb.Message {
	return &brbpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.Message]()}
}

type StartMessage struct {
	Id   int64
	N    int64
	Data []uint8
}

func StartMessageFromPb(pb *brbpb.StartMessage) *StartMessage {
	return &StartMessage{
		Id:   pb.Id,
		N:    pb.N,
		Data: pb.Data,
	}
}

func (m *StartMessage) Pb() *brbpb.StartMessage {
	return &brbpb.StartMessage{
		Id:   m.Id,
		N:    m.N,
		Data: m.Data,
	}
}

func (*StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.StartMessage]()}
}

type EchoMessage struct {
	Id   int64
	N    int64
	Data []uint8
}

func EchoMessageFromPb(pb *brbpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		Id:   pb.Id,
		N:    pb.N,
		Data: pb.Data,
	}
}

func (m *EchoMessage) Pb() *brbpb.EchoMessage {
	return &brbpb.EchoMessage{
		Id:   m.Id,
		N:    m.N,
		Data: m.Data,
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.EchoMessage]()}
}

type ReadyMessage struct {
	Id   int64
	N    int64
	Data []uint8
}

func ReadyMessageFromPb(pb *brbpb.ReadyMessage) *ReadyMessage {
	return &ReadyMessage{
		Id:   pb.Id,
		N:    pb.N,
		Data: pb.Data,
	}
}

func (m *ReadyMessage) Pb() *brbpb.ReadyMessage {
	return &brbpb.ReadyMessage{
		Id:   m.Id,
		N:    m.N,
		Data: m.Data,
	}
}

func (*ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbpb.ReadyMessage]()}
}
