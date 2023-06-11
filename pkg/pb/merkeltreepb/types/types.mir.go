package merkletreepbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	merkletreepb "github.com/filecoin-project/mir/pkg/pb/merkletreepb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() merkletreepb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb merkletreepb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *merkletreepb.Event_VerifyRequest:
		return &Event_VerifyRequest{VerifyRequest: VerifyRequestFromPb(pb.VerifyRequest)}
	case *merkletreepb.Event_VerifyResult:
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

func (w *Event_VerifyRequest) Pb() merkletreepb.Event_Type {
	return &merkletreepb.Event_VerifyRequest{VerifyRequest: (w.VerifyRequest).Pb()}
}

func (*Event_VerifyRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkletreepb.Event_VerifyRequest]()}
}

type Event_VerifyResult struct {
	VerifyResult *VerifyResult
}

func (*Event_VerifyResult) isEvent_Type() {}

func (w *Event_VerifyResult) Unwrap() *VerifyResult {
	return w.VerifyResult
}

func (w *Event_VerifyResult) Pb() merkletreepb.Event_Type {
	return &merkletreepb.Event_VerifyResult{VerifyResult: (w.VerifyResult).Pb()}
}

func (*Event_VerifyResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkletreepb.Event_VerifyResult]()}
}

func EventFromPb(pb *merkletreepb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *merkletreepb.Event {
	return &merkletreepb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkletreepb.Event]()}
}

type VerifyRequest struct {
	Chunk    []uint8
	RootHash []uint8
	Proof    [][]uint8
}

func VerifyRequestFromPb(pb *merkletreepb.VerifyRequest) *VerifyRequest {
	return &VerifyRequest{
		Chunk:    pb.Chunk,
		RootHash: pb.RootHash,
		Proof:    pb.Proof,
	}
}

func (m *VerifyRequest) Pb() *merkletreepb.VerifyRequest {
	return &merkletreepb.VerifyRequest{
		Chunk:    m.Chunk,
		RootHash: m.RootHash,
		Proof:    m.Proof,
	}
}

func (*VerifyRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkletreepb.VerifyRequest]()}
}

type VerifyResult struct {
	Success bool
}

func VerifyResultFromPb(pb *merkletreepb.VerifyResult) *VerifyResult {
	return &VerifyResult{
		Success: pb.Success,
	}
}

func (m *VerifyResult) Pb() *merkletreepb.VerifyResult {
	return &merkletreepb.VerifyResult{
		Success: m.Success,
	}
}

func (*VerifyResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*merkletreepb.VerifyResult]()}
}
