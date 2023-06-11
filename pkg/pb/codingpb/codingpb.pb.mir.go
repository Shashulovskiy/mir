package codingpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_EncodeRequest)(nil)),
		reflect.TypeOf((*Event_EncodeResult)(nil)),
		reflect.TypeOf((*Event_DecodeRequest)(nil)),
		reflect.TypeOf((*Event_DecodeResult)(nil)),
		reflect.TypeOf((*Event_RebuildRequest)(nil)),
		reflect.TypeOf((*Event_RebuildResult)(nil)),
	}
}
