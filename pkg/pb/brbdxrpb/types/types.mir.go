package brbdxrpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbdxrpb "github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

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
	Id   int64
	N    int64
	Data []uint8
}

func StartMessageFromPb(pb *brbdxrpb.StartMessage) *StartMessage {
	return &StartMessage{
		Id:   pb.Id,
		N:    pb.N,
		Data: pb.Data,
	}
}

func (m *StartMessage) Pb() *brbdxrpb.StartMessage {
	return &brbdxrpb.StartMessage{
		Id:   m.Id,
		N:    m.N,
		Data: m.Data,
	}
}

func (*StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.StartMessage]()}
}

type EchoMessage struct {
	Id    int64
	N     int64
	Hash  []uint8
	Chunk []uint8
}

func EchoMessageFromPb(pb *brbdxrpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		Id:    pb.Id,
		N:     pb.N,
		Hash:  pb.Hash,
		Chunk: pb.Chunk,
	}
}

func (m *EchoMessage) Pb() *brbdxrpb.EchoMessage {
	return &brbdxrpb.EchoMessage{
		Id:    m.Id,
		N:     m.N,
		Hash:  m.Hash,
		Chunk: m.Chunk,
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.EchoMessage]()}
}

type ReadyMessage struct {
	Id    int64
	N     int64
	Hash  []uint8
	Chunk []uint8
}

func ReadyMessageFromPb(pb *brbdxrpb.ReadyMessage) *ReadyMessage {
	return &ReadyMessage{
		Id:    pb.Id,
		N:     pb.N,
		Hash:  pb.Hash,
		Chunk: pb.Chunk,
	}
}

func (m *ReadyMessage) Pb() *brbdxrpb.ReadyMessage {
	return &brbdxrpb.ReadyMessage{
		Id:    m.Id,
		N:     m.N,
		Hash:  m.Hash,
		Chunk: m.Chunk,
	}
}

func (*ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbdxrpb.ReadyMessage]()}
}
