package brbctpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	brbctpb "github.com/filecoin-project/mir/pkg/pb/brbctpb"
	commonpb "github.com/filecoin-project/mir/pkg/pb/commonpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

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
	Id       int64
	N        int64
	Chunk    []uint8
	RootHash []uint8
	Proof    *commonpb.MerklePath
}

func StartMessageFromPb(pb *brbctpb.StartMessage) *StartMessage {
	return &StartMessage{
		Id:       pb.Id,
		N:        pb.N,
		Chunk:    pb.Chunk,
		RootHash: pb.RootHash,
		Proof:    pb.Proof,
	}
}

func (m *StartMessage) Pb() *brbctpb.StartMessage {
	return &brbctpb.StartMessage{
		Id:       m.Id,
		N:        m.N,
		Chunk:    m.Chunk,
		RootHash: m.RootHash,
		Proof:    m.Proof,
	}
}

func (*StartMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.StartMessage]()}
}

type EchoMessage struct {
	Id       int64
	N        int64
	Chunk    []uint8
	RootHash []uint8
	Proof    *commonpb.MerklePath
}

func EchoMessageFromPb(pb *brbctpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		Id:       pb.Id,
		N:        pb.N,
		Chunk:    pb.Chunk,
		RootHash: pb.RootHash,
		Proof:    pb.Proof,
	}
}

func (m *EchoMessage) Pb() *brbctpb.EchoMessage {
	return &brbctpb.EchoMessage{
		Id:       m.Id,
		N:        m.N,
		Chunk:    m.Chunk,
		RootHash: m.RootHash,
		Proof:    m.Proof,
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.EchoMessage]()}
}

type ReadyMessage struct {
	Id       int64
	N        int64
	RootHash []uint8
}

func ReadyMessageFromPb(pb *brbctpb.ReadyMessage) *ReadyMessage {
	return &ReadyMessage{
		Id:       pb.Id,
		N:        pb.N,
		RootHash: pb.RootHash,
	}
}

func (m *ReadyMessage) Pb() *brbctpb.ReadyMessage {
	return &brbctpb.ReadyMessage{
		Id:       m.Id,
		N:        m.N,
		RootHash: m.RootHash,
	}
}

func (*ReadyMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*brbctpb.ReadyMessage]()}
}
