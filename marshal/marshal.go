package marshal

import (
	"errors"
)

type MarshalledObject struct {
	MajorVersion byte
	MinorVersion byte
	data         []byte
}

type marshalledObjectType byte

var TypeMismatch = errors.New("gorails/marshal: an attempt to implicitly typecast the marshalled object")

const (
	TYPE_UNKNOWN marshalledObjectType = 0
	TYPE_NIL     marshalledObjectType = 1
	TYPE_BOOLEAN marshalledObjectType = 2
	TYPE_INTEGER marshalledObjectType = 3
	TYPE_FLOAT   marshalledObjectType = 4
	TYPE_ARRAY   marshalledObjectType = 5
	TYPE_MAP     marshalledObjectType = 6
)

func CreateMarshalledObject(serialized_data []byte) *MarshalledObject {
	return &(MarshalledObject{serialized_data[0], serialized_data[1], serialized_data[2:]})
}

func assertType(obj *MarshalledObject, expected_type marshalledObjectType) (err error) {
	if obj.GetType() != expected_type {
		err = TypeMismatch
	}

	return
}

func (obj *MarshalledObject) GetType() (object_type marshalledObjectType) {
	switch obj.data[0] {
	case '0':
		object_type = TYPE_NIL
	case 'T', 'F':
		object_type = TYPE_BOOLEAN
	case 'i':
		object_type = TYPE_INTEGER
	case 'f':
		object_type = TYPE_FLOAT
	case '[':
		object_type = TYPE_ARRAY
	case '{':
		object_type = TYPE_MAP
	default:
		object_type = TYPE_UNKNOWN
	}

	return
}

func (obj *MarshalledObject) GetAsBoolean() (value bool, err error) {
	err = assertType(obj, TYPE_BOOLEAN)
	if err == nil {
		value = obj.data[0] == 'T'
	}

	return
}

func (obj *MarshalledObject) GetAsInteger() (value int, err error) {
	err = assertType(obj, TYPE_INTEGER)
	if err != nil {
		return
	}

	if obj.data[1] > 0x05 && obj.data[1] < 0xfb {
		value = int(obj.data[1])

		if value > 0x7f {
			value = -(0xff ^ value + 1) + 5
		} else {
			value -= 5
		}
	} else if obj.data[1] <= 0x05 {
		value = 0
		i := obj.data[1]

		for ; i > 0; i-- {
			value = value<<8 + int(obj.data[i+1])
		}
	} else {
		value = 0
		i := 0xff - obj.data[1] + 1

		for ; i > 0; i-- {
			value = value<<8 + (0xff - int(obj.data[i+1]))
		}

		value = -(value + 1)
	}

	return
}
