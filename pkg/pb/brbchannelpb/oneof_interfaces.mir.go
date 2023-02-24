package brbchannelpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_Request) Unwrap() *BroadcastRequest {
	return w.Request
}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_AnyMessage) Unwrap() *AnyMessage {
	return w.AnyMessage
}
