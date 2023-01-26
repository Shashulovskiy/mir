package merkeltreepb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_VerifyRequest)(nil)),
		reflect.TypeOf((*Event_VerifyResult)(nil)),
	}
}
