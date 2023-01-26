package brbctpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbctpb "github.com/filecoin-project/mir/pkg/pb/brbctpb"
	commonpb "github.com/filecoin-project/mir/pkg/pb/commonpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() brbctpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb brbctpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *brbctpb.Event_Request:
		return &Event_Request{Request: BroadcastRequestFromPb(pb.Request)}
	case *brbctpb.Event_Deliver:
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

func (w *Event_Request) Pb() brbctpb.Event_Type {
	return &brbctpb.Event_Request{Request: (w.Request).Pb()}
}

func (*Event_Request) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Event_Request]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() brbctpb.Event_Type {
	return &brbctpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Event_Deliver]()}
}

func EventFromPb(pb *brbctpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *brbctpb.Event {
	return &brbctpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Event]()}
}

type BroadcastRequest struct {
	Data []uint8
}

func BroadcastRequestFromPb(pb *brbctpb.BroadcastRequest) *BroadcastRequest {
	return &BroadcastRequest{
		Data: pb.Data,
	}
}

func (m *BroadcastRequest) Pb() *brbctpb.BroadcastRequest {
	return &brbctpb.BroadcastRequest{
		Data: m.Data,
	}
}

func (*BroadcastRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.BroadcastRequest]()}
}

type Deliver struct {
	Data []uint8
}

func DeliverFromPb(pb *brbctpb.Deliver) *Deliver {
	return &Deliver{
		Data: pb.Data,
	}
}

func (m *Deliver) Pb() *brbctpb.Deliver {
	return &brbctpb.Deliver{
		Data: m.Data,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() brbctpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb brbctpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *brbctpb.Message_StartMessage:
		return &Message_StartMessage{StartMessage: StartMessageFromPb(pb.StartMessage)}
	case *brbctpb.Message_EchoMessage:
		return &Message_EchoMessage{EchoMessage: EchoMessageFromPb(pb.EchoMessage)}
	case *brbctpb.Message_ReadyMessage:
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

func (w *Message_StartMessage) Pb() brbctpb.Message_Type {
	return &brbctpb.Message_StartMessage{StartMessage: (w.StartMessage).Pb()}
}

func (*Message_StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Message_StartMessage]()}
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage
}

func (*Message_EchoMessage) isMessage_Type() {}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_EchoMessage) Pb() brbctpb.Message_Type {
	return &brbctpb.Message_EchoMessage{EchoMessage: (w.EchoMessage).Pb()}
}

func (*Message_EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Message_EchoMessage]()}
}

type Message_ReadyMessage struct {
	ReadyMessage *ReadyMessage
}

func (*Message_ReadyMessage) isMessage_Type() {}

func (w *Message_ReadyMessage) Unwrap() *ReadyMessage {
	return w.ReadyMessage
}

func (w *Message_ReadyMessage) Pb() brbctpb.Message_Type {
	return &brbctpb.Message_ReadyMessage{ReadyMessage: (w.ReadyMessage).Pb()}
}

func (*Message_ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Message_ReadyMessage]()}
}

func MessageFromPb(pb *brbctpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *brbctpb.Message {
	return &brbctpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.Message]()}
}

type StartMessage struct {
	Chunk    []uint8
	RootHash []uint8
	Proof    *commonpb.MerklePath
}

func StartMessageFromPb(pb *brbctpb.StartMessage) *StartMessage {
	return &StartMessage{
		Chunk:    pb.Chunk,
		RootHash: pb.RootHash,
		Proof:    pb.Proof,
	}
}

func (m *StartMessage) Pb() *brbctpb.StartMessage {
	return &brbctpb.StartMessage{
		Chunk:    m.Chunk,
		RootHash: m.RootHash,
		Proof:    m.Proof,
	}
}

func (*StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.StartMessage]()}
}

type EchoMessage struct {
	Chunk    []uint8
	RootHash []uint8
	Proof    *commonpb.MerklePath
}

func EchoMessageFromPb(pb *brbctpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		Chunk:    pb.Chunk,
		RootHash: pb.RootHash,
		Proof:    pb.Proof,
	}
}

func (m *EchoMessage) Pb() *brbctpb.EchoMessage {
	return &brbctpb.EchoMessage{
		Chunk:    m.Chunk,
		RootHash: m.RootHash,
		Proof:    m.Proof,
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.EchoMessage]()}
}

type ReadyMessage struct {
	RootHash []uint8
}

func ReadyMessageFromPb(pb *brbctpb.ReadyMessage) *ReadyMessage {
	return &ReadyMessage{
		RootHash: pb.RootHash,
	}
}

func (m *ReadyMessage) Pb() *brbctpb.ReadyMessage {
	return &brbctpb.ReadyMessage{
		RootHash: m.RootHash,
	}
}

func (*ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.ReadyMessage]()}
}
