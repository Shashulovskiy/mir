package brbencodedpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbencodedpb "github.com/filecoin-project/mir/pkg/pb/brbencodedpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() brbencodedpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb brbencodedpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *brbencodedpb.Event_Request:
		return &Event_Request{Request: BroadcastRequestFromPb(pb.Request)}
	case *brbencodedpb.Event_Deliver:
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

func (w *Event_Request) Pb() brbencodedpb.Event_Type {
	return &brbencodedpb.Event_Request{Request: (w.Request).Pb()}
}

func (*Event_Request) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Event_Request]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() brbencodedpb.Event_Type {
	return &brbencodedpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Event_Deliver]()}
}

func EventFromPb(pb *brbencodedpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *brbencodedpb.Event {
	return &brbencodedpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Event]()}
}

type BroadcastRequest struct {
	Data []uint8
}

func BroadcastRequestFromPb(pb *brbencodedpb.BroadcastRequest) *BroadcastRequest {
	return &BroadcastRequest{
		Data: pb.Data,
	}
}

func (m *BroadcastRequest) Pb() *brbencodedpb.BroadcastRequest {
	return &brbencodedpb.BroadcastRequest{
		Data: m.Data,
	}
}

func (*BroadcastRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.BroadcastRequest]()}
}

type Deliver struct {
	Data []uint8
}

func DeliverFromPb(pb *brbencodedpb.Deliver) *Deliver {
	return &Deliver{
		Data: pb.Data,
	}
}

func (m *Deliver) Pb() *brbencodedpb.Deliver {
	return &brbencodedpb.Deliver{
		Data: m.Data,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() brbencodedpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb brbencodedpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *brbencodedpb.Message_StartMessage:
		return &Message_StartMessage{StartMessage: StartMessageFromPb(pb.StartMessage)}
	case *brbencodedpb.Message_EchoMessage:
		return &Message_EchoMessage{EchoMessage: EchoMessageFromPb(pb.EchoMessage)}
	case *brbencodedpb.Message_ReadyMessage:
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

func (w *Message_StartMessage) Pb() brbencodedpb.Message_Type {
	return &brbencodedpb.Message_StartMessage{StartMessage: (w.StartMessage).Pb()}
}

func (*Message_StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Message_StartMessage]()}
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage
}

func (*Message_EchoMessage) isMessage_Type() {}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_EchoMessage) Pb() brbencodedpb.Message_Type {
	return &brbencodedpb.Message_EchoMessage{EchoMessage: (w.EchoMessage).Pb()}
}

func (*Message_EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Message_EchoMessage]()}
}

type Message_ReadyMessage struct {
	ReadyMessage *ReadyMessage
}

func (*Message_ReadyMessage) isMessage_Type() {}

func (w *Message_ReadyMessage) Unwrap() *ReadyMessage {
	return w.ReadyMessage
}

func (w *Message_ReadyMessage) Pb() brbencodedpb.Message_Type {
	return &brbencodedpb.Message_ReadyMessage{ReadyMessage: (w.ReadyMessage).Pb()}
}

func (*Message_ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Message_ReadyMessage]()}
}

func MessageFromPb(pb *brbencodedpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *brbencodedpb.Message {
	return &brbencodedpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.Message]()}
}

type StartMessage struct {
	Data []uint8
}

func StartMessageFromPb(pb *brbencodedpb.StartMessage) *StartMessage {
	return &StartMessage{
		Data: pb.Data,
	}
}

func (m *StartMessage) Pb() *brbencodedpb.StartMessage {
	return &brbencodedpb.StartMessage{
		Data: m.Data,
	}
}

func (*StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.StartMessage]()}
}

type EchoMessage struct {
	Hash  []uint8
	Chunk []uint8
}

func EchoMessageFromPb(pb *brbencodedpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		Hash:  pb.Hash,
		Chunk: pb.Chunk,
	}
}

func (m *EchoMessage) Pb() *brbencodedpb.EchoMessage {
	return &brbencodedpb.EchoMessage{
		Hash:  m.Hash,
		Chunk: m.Chunk,
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.EchoMessage]()}
}

type ReadyMessage struct {
	Hash  []uint8
	Chunk []uint8
}

func ReadyMessageFromPb(pb *brbencodedpb.ReadyMessage) *ReadyMessage {
	return &ReadyMessage{
		Hash:  pb.Hash,
		Chunk: pb.Chunk,
	}
}

func (m *ReadyMessage) Pb() *brbencodedpb.ReadyMessage {
	return &brbencodedpb.ReadyMessage{
		Hash:  m.Hash,
		Chunk: m.Chunk,
	}
}

func (*ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbencodedpb.ReadyMessage]()}
}
