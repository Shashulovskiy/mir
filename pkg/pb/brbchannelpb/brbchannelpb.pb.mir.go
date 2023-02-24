package brbchannelpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_Request)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
	}
}

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_AnyMessage)(nil)),
	}
}
