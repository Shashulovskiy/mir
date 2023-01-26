package merkeltreepb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_VerifyRequest) Unwrap() *VerifyRequest {
	return w.VerifyRequest
}

func (w *Event_VerifyResult) Unwrap() *VerifyResult {
	return w.VerifyResult
}
