package merkeltreepbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	merkeltreepb "github.com/filecoin-project/mir/pkg/pb/merkeltreepb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() merkeltreepb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb merkeltreepb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *merkeltreepb.Event_VerifyRequest:
		return &Event_VerifyRequest{VerifyRequest: VerifyRequestFromPb(pb.VerifyRequest)}
	case *merkeltreepb.Event_VerifyResult:
		return &Event_VerifyResult{VerifyResult: VerifyResultFromPb(pb.VerifyResult)}
	}
	return nil
}

type Event_VerifyRequest struct {
	VerifyRequest *VerifyRequest
}

func (*Event_VerifyRequest) isEvent_Type() {}

func (w *Event_VerifyRequest) Unwrap() *VerifyRequest {
	return w.VerifyRequest
}

func (w *Event_VerifyRequest) Pb() merkeltreepb.Event_Type {
	return &merkeltreepb.Event_VerifyRequest{VerifyRequest: (w.VerifyRequest).Pb()}
}

func (*Event_VerifyRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkeltreepb.Event_VerifyRequest]()}
}

type Event_VerifyResult struct {
	VerifyResult *VerifyResult
}

func (*Event_VerifyResult) isEvent_Type() {}

func (w *Event_VerifyResult) Unwrap() *VerifyResult {
	return w.VerifyResult
}

func (w *Event_VerifyResult) Pb() merkeltreepb.Event_Type {
	return &merkeltreepb.Event_VerifyResult{VerifyResult: (w.VerifyResult).Pb()}
}

func (*Event_VerifyResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkeltreepb.Event_VerifyResult]()}
}

func EventFromPb(pb *merkeltreepb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *merkeltreepb.Event {
	return &merkeltreepb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkeltreepb.Event]()}
}

type VerifyRequest struct {
	Chunk    []uint8
	RootHash []uint8
	Proof    [][]uint8
}

func VerifyRequestFromPb(pb *merkeltreepb.VerifyRequest) *VerifyRequest {
	return &VerifyRequest{
		Chunk:    pb.Chunk,
		RootHash: pb.RootHash,
		Proof:    pb.Proof,
	}
}

func (m *VerifyRequest) Pb() *merkeltreepb.VerifyRequest {
	return &merkeltreepb.VerifyRequest{
		Chunk:    m.Chunk,
		RootHash: m.RootHash,
		Proof:    m.Proof,
	}
}

func (*VerifyRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkeltreepb.VerifyRequest]()}
}

type VerifyResult struct {
	Success bool
}

func VerifyResultFromPb(pb *merkeltreepb.VerifyResult) *VerifyResult {
	return &VerifyResult{
		Success: pb.Success,
	}
}

func (m *VerifyResult) Pb() *merkeltreepb.VerifyResult {
	return &merkeltreepb.VerifyResult{
		Success: m.Success,
	}
}

func (*VerifyResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkeltreepb.VerifyResult]()}
}
