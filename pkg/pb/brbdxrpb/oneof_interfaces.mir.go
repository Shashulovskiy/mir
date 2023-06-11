package brbdxrpb

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_StartMessage) Unwrap() *StartMessage {
	return w.StartMessage
}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_ReadyMessage) Unwrap() *ReadyMessage {
	return w.ReadyMessage
}
