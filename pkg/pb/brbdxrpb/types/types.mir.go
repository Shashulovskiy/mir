package brbdxrpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbdxrpb "github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() brbdxrpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb brbdxrpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *brbdxrpb.Event_Request:
		return &Event_Request{Request: BroadcastRequestFromPb(pb.Request)}
	case *brbdxrpb.Event_Deliver:
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

func (w *Event_Request) Pb() brbdxrpb.Event_Type {
	return &brbdxrpb.Event_Request{Request: (w.Request).Pb()}
}

func (*Event_Request) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Event_Request]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() brbdxrpb.Event_Type {
	return &brbdxrpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Event_Deliver]()}
}

func EventFromPb(pb *brbdxrpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *brbdxrpb.Event {
	return &brbdxrpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Event]()}
}

type BroadcastRequest struct {
	Data []uint8
}

func BroadcastRequestFromPb(pb *brbdxrpb.BroadcastRequest) *BroadcastRequest {
	return &BroadcastRequest{
		Data: pb.Data,
	}
}

func (m *BroadcastRequest) Pb() *brbdxrpb.BroadcastRequest {
	return &brbdxrpb.BroadcastRequest{
		Data: m.Data,
	}
}

func (*BroadcastRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.BroadcastRequest]()}
}

type Deliver struct {
	Data []uint8
}

func DeliverFromPb(pb *brbdxrpb.Deliver) *Deliver {
	return &Deliver{
		Data: pb.Data,
	}
}

func (m *Deliver) Pb() *brbdxrpb.Deliver {
	return &brbdxrpb.Deliver{
		Data: m.Data,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() brbdxrpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb brbdxrpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *brbdxrpb.Message_StartMessage:
		return &Message_StartMessage{StartMessage: StartMessageFromPb(pb.StartMessage)}
	case *brbdxrpb.Message_EchoMessage:
		return &Message_EchoMessage{EchoMessage: EchoMessageFromPb(pb.EchoMessage)}
	case *brbdxrpb.Message_ReadyMessage:
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

func (w *Message_StartMessage) Pb() brbdxrpb.Message_Type {
	return &brbdxrpb.Message_StartMessage{StartMessage: (w.StartMessage).Pb()}
}

func (*Message_StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Message_StartMessage]()}
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage
}

func (*Message_EchoMessage) isMessage_Type() {}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_EchoMessage) Pb() brbdxrpb.Message_Type {
	return &brbdxrpb.Message_EchoMessage{EchoMessage: (w.EchoMessage).Pb()}
}

func (*Message_EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Message_EchoMessage]()}
}

type Message_ReadyMessage struct {
	ReadyMessage *ReadyMessage
}

func (*Message_ReadyMessage) isMessage_Type() {}

func (w *Message_ReadyMessage) Unwrap() *ReadyMessage {
	return w.ReadyMessage
}

func (w *Message_ReadyMessage) Pb() brbdxrpb.Message_Type {
	return &brbdxrpb.Message_ReadyMessage{ReadyMessage: (w.ReadyMessage).Pb()}
}

func (*Message_ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Message_ReadyMessage]()}
}

func MessageFromPb(pb *brbdxrpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *brbdxrpb.Message {
	return &brbdxrpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.Message]()}
}

type StartMessage struct {
	Data []uint8
}

func StartMessageFromPb(pb *brbdxrpb.StartMessage) *StartMessage {
	return &StartMessage{
		Data: pb.Data,
	}
}

func (m *StartMessage) Pb() *brbdxrpb.StartMessage {
	return &brbdxrpb.StartMessage{
		Data: m.Data,
	}
}

func (*StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.StartMessage]()}
}

type EchoMessage struct {
	Hash  []uint8
	Chunk []uint8
}

func EchoMessageFromPb(pb *brbdxrpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		Hash:  pb.Hash,
		Chunk: pb.Chunk,
	}
}

func (m *EchoMessage) Pb() *brbdxrpb.EchoMessage {
	return &brbdxrpb.EchoMessage{
		Hash:  m.Hash,
		Chunk: m.Chunk,
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.EchoMessage]()}
}

type ReadyMessage struct {
	Hash  []uint8
	Chunk []uint8
}

func ReadyMessageFromPb(pb *brbdxrpb.ReadyMessage) *ReadyMessage {
	return &ReadyMessage{
		Hash:  pb.Hash,
		Chunk: pb.Chunk,
	}
}

func (m *ReadyMessage) Pb() *brbdxrpb.ReadyMessage {
	return &brbdxrpb.ReadyMessage{
		Hash:  m.Hash,
		Chunk: m.Chunk,
	}
}

func (*ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.ReadyMessage]()}
}
