package codingpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_EncodeRequest) Unwrap() *EncodeRequest {
	return w.EncodeRequest
}

func (w *Event_EncodeResult) Unwrap() *EncodeResult {
	return w.EncodeResult
}

func (w *Event_DecodeRequest) Unwrap() *DecodeRequest {
	return w.DecodeRequest
}

func (w *Event_DecodeResult) Unwrap() *DecodeResult {
	return w.DecodeResult
}

func (w *Event_RebuildRequest) Unwrap() *RebuildRequest {
	return w.RebuildRequest
}

func (w *Event_RebuildResult) Unwrap() *RebuildResult {
	return w.RebuildResult
}
