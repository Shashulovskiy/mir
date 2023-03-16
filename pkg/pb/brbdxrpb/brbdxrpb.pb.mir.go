package brbdxrpb

import (
	reflect "reflect"
)

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_StartMessage)(nil)),
		reflect.TypeOf((*Message_EchoMessage)(nil)),
		reflect.TypeOf((*Message_ReadyMessage)(nil)),
	}
}
