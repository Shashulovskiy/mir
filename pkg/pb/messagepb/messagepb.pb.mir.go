package messagepb

import (
	reflect "reflect"
)

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_Iss)(nil)),
		reflect.TypeOf((*Message_Bcb)(nil)),
		reflect.TypeOf((*Message_MultisigCollector)(nil)),
		reflect.TypeOf((*Message_Pingpong)(nil)),
		reflect.TypeOf((*Message_Checkpoint)(nil)),
		reflect.TypeOf((*Message_SbMessage)(nil)),
		reflect.TypeOf((*Message_Brb)(nil)),
		reflect.TypeOf((*Message_Brbdxr)(nil)),
		reflect.TypeOf((*Message_Brbct)(nil)),
	}
}
