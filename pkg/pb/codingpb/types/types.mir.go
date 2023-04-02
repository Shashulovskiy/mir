package codingpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	codingpb "github.com/filecoin-project/mir/pkg/pb/codingpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() codingpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb codingpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *codingpb.Event_EncodeRequest:
		return &Event_EncodeRequest{EncodeRequest: EncodeRequestFromPb(pb.EncodeRequest)}
	case *codingpb.Event_EncodeResult:
		return &Event_EncodeResult{EncodeResult: EncodeResultFromPb(pb.EncodeResult)}
	case *codingpb.Event_DecodeRequest:
		return &Event_DecodeRequest{DecodeRequest: DecodeRequestFromPb(pb.DecodeRequest)}
	case *codingpb.Event_DecodeResult:
		return &Event_DecodeResult{DecodeResult: DecodeResultFromPb(pb.DecodeResult)}
	case *codingpb.Event_RebuildRequest:
		return &Event_RebuildRequest{RebuildRequest: RebuildRequestFromPb(pb.RebuildRequest)}
	case *codingpb.Event_RebuildResult:
		return &Event_RebuildResult{RebuildResult: RebuildResultFromPb(pb.RebuildResult)}
	}
	return nil
}

type Event_EncodeRequest struct {
	EncodeRequest *EncodeRequest
}

func (*Event_EncodeRequest) isEvent_Type() {}

func (w *Event_EncodeRequest) Unwrap() *EncodeRequest {
	return w.EncodeRequest
}

func (w *Event_EncodeRequest) Pb() codingpb.Event_Type {
	return &codingpb.Event_EncodeRequest{EncodeRequest: (w.EncodeRequest).Pb()}
}

func (*Event_EncodeRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event_EncodeRequest]()}
}

type Event_EncodeResult struct {
	EncodeResult *EncodeResult
}

func (*Event_EncodeResult) isEvent_Type() {}

func (w *Event_EncodeResult) Unwrap() *EncodeResult {
	return w.EncodeResult
}

func (w *Event_EncodeResult) Pb() codingpb.Event_Type {
	return &codingpb.Event_EncodeResult{EncodeResult: (w.EncodeResult).Pb()}
}

func (*Event_EncodeResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event_EncodeResult]()}
}

type Event_DecodeRequest struct {
	DecodeRequest *DecodeRequest
}

func (*Event_DecodeRequest) isEvent_Type() {}

func (w *Event_DecodeRequest) Unwrap() *DecodeRequest {
	return w.DecodeRequest
}

func (w *Event_DecodeRequest) Pb() codingpb.Event_Type {
	return &codingpb.Event_DecodeRequest{DecodeRequest: (w.DecodeRequest).Pb()}
}

func (*Event_DecodeRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event_DecodeRequest]()}
}

type Event_DecodeResult struct {
	DecodeResult *DecodeResult
}

func (*Event_DecodeResult) isEvent_Type() {}

func (w *Event_DecodeResult) Unwrap() *DecodeResult {
	return w.DecodeResult
}

func (w *Event_DecodeResult) Pb() codingpb.Event_Type {
	return &codingpb.Event_DecodeResult{DecodeResult: (w.DecodeResult).Pb()}
}

func (*Event_DecodeResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event_DecodeResult]()}
}

type Event_RebuildRequest struct {
	RebuildRequest *RebuildRequest
}

func (*Event_RebuildRequest) isEvent_Type() {}

func (w *Event_RebuildRequest) Unwrap() *RebuildRequest {
	return w.RebuildRequest
}

func (w *Event_RebuildRequest) Pb() codingpb.Event_Type {
	return &codingpb.Event_RebuildRequest{RebuildRequest: (w.RebuildRequest).Pb()}
}

func (*Event_RebuildRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event_RebuildRequest]()}
}

type Event_RebuildResult struct {
	RebuildResult *RebuildResult
}

func (*Event_RebuildResult) isEvent_Type() {}

func (w *Event_RebuildResult) Unwrap() *RebuildResult {
	return w.RebuildResult
}

func (w *Event_RebuildResult) Pb() codingpb.Event_Type {
	return &codingpb.Event_RebuildResult{RebuildResult: (w.RebuildResult).Pb()}
}

func (*Event_RebuildResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event_RebuildResult]()}
}

func EventFromPb(pb *codingpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *codingpb.Event {
	return &codingpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.Event]()}
}

type EncodeRequest struct {
	TotalShards int64
	DataShards  int64
	PaddedData  []uint8
	Origin      *codingpb.EncodeOrigin
}

func EncodeRequestFromPb(pb *codingpb.EncodeRequest) *EncodeRequest {
	return &EncodeRequest{
		TotalShards: pb.TotalShards,
		DataShards:  pb.DataShards,
		PaddedData:  pb.PaddedData,
		Origin:      pb.Origin,
	}
}

func (m *EncodeRequest) Pb() *codingpb.EncodeRequest {
	return &codingpb.EncodeRequest{
		TotalShards: m.TotalShards,
		DataShards:  m.DataShards,
		PaddedData:  m.PaddedData,
		Origin:      m.Origin,
	}
}

func (*EncodeRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.EncodeRequest]()}
}

type EncodeResult struct {
	Encoded [][]uint8
	Origin  *codingpb.EncodeOrigin
}

func EncodeResultFromPb(pb *codingpb.EncodeResult) *EncodeResult {
	return &EncodeResult{
		Encoded: pb.Encoded,
		Origin:  pb.Origin,
	}
}

func (m *EncodeResult) Pb() *codingpb.EncodeResult {
	return &codingpb.EncodeResult{
		Encoded: m.Encoded,
		Origin:  m.Origin,
	}
}

func (*EncodeResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.EncodeResult]()}
}

type DecodeRequest struct {
	TotalShards int64
	DataShards  int64
	Shares      []*codingpb.Share
	Origin      *codingpb.DecodeOrigin
}

func DecodeRequestFromPb(pb *codingpb.DecodeRequest) *DecodeRequest {
	return &DecodeRequest{
		TotalShards: pb.TotalShards,
		DataShards:  pb.DataShards,
		Shares:      pb.Shares,
		Origin:      pb.Origin,
	}
}

func (m *DecodeRequest) Pb() *codingpb.DecodeRequest {
	return &codingpb.DecodeRequest{
		TotalShards: m.TotalShards,
		DataShards:  m.DataShards,
		Shares:      m.Shares,
		Origin:      m.Origin,
	}
}

func (*DecodeRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.DecodeRequest]()}
}

type DecodeResult struct {
	Success bool
	Decoded []uint8
	Origin  *codingpb.DecodeOrigin
}

func DecodeResultFromPb(pb *codingpb.DecodeResult) *DecodeResult {
	return &DecodeResult{
		Success: pb.Success,
		Decoded: pb.Decoded,
		Origin:  pb.Origin,
	}
}

func (m *DecodeResult) Pb() *codingpb.DecodeResult {
	return &codingpb.DecodeResult{
		Success: m.Success,
		Decoded: m.Decoded,
		Origin:  m.Origin,
	}
}

func (*DecodeResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.DecodeResult]()}
}

type RebuildRequest struct {
	TotalShards int64
	DataShards  int64
	Shares      []*codingpb.Share
	Origin      *codingpb.RebuildOrigin
}

func RebuildRequestFromPb(pb *codingpb.RebuildRequest) *RebuildRequest {
	return &RebuildRequest{
		TotalShards: pb.TotalShards,
		DataShards:  pb.DataShards,
		Shares:      pb.Shares,
		Origin:      pb.Origin,
	}
}

func (m *RebuildRequest) Pb() *codingpb.RebuildRequest {
	return &codingpb.RebuildRequest{
		TotalShards: m.TotalShards,
		DataShards:  m.DataShards,
		Shares:      m.Shares,
		Origin:      m.Origin,
	}
}

func (*RebuildRequest) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.RebuildRequest]()}
}

type RebuildResult struct {
	Success bool
	Decoded []uint8
	Origin  *codingpb.RebuildOrigin
}

func RebuildResultFromPb(pb *codingpb.RebuildResult) *RebuildResult {
	return &RebuildResult{
		Success: pb.Success,
		Decoded: pb.Decoded,
		Origin:  pb.Origin,
	}
}

func (m *RebuildResult) Pb() *codingpb.RebuildResult {
	return &codingpb.RebuildResult{
		Success: m.Success,
		Decoded: m.Decoded,
		Origin:  m.Origin,
	}
}

func (*RebuildResult) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*codingpb.RebuildResult]()}
}
