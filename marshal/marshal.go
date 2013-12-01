package marshal

import (
	"errors"
	"strconv"
)

type MarshalledObject struct {
	MajorVersion byte
	MinorVersion byte
	data         []byte
}

type marshalledObjectType byte

var TypeMismatch = errors.New("gorails/marshal: an attempt to implicitly typecast the marshalled object")
var IncompleteData = errors.New("gorails/marshal: incomplete data")

const (
	TYPE_UNKNOWN marshalledObjectType = 0
	TYPE_NIL     marshalledObjectType = 1
	TYPE_BOOL    marshalledObjectType = 2
	TYPE_INTEGER marshalledObjectType = 3
	TYPE_FLOAT   marshalledObjectType = 4
	TYPE_STRING  marshalledObjectType = 5
	TYPE_ARRAY   marshalledObjectType = 6
	TYPE_MAP     marshalledObjectType = 7
)

func CreateMarshalledObject(serialized_data []byte) *MarshalledObject {
	return &(MarshalledObject{serialized_data[0], serialized_data[1], serialized_data[2:]})
}

func (obj *MarshalledObject) GetType() marshalledObjectType {
	if len(obj.data) == 0 {
		return TYPE_UNKNOWN
	}

	switch obj.data[0] {
	case '0':
		return TYPE_NIL
	case 'T', 'F':
		return TYPE_BOOL
	case 'i':
		return TYPE_INTEGER
	case 'f':
		return TYPE_FLOAT
	case ':':
		return TYPE_STRING
	case 'I':
		if len(obj.data) > 1 && obj.data[1] == '"' {
			return TYPE_STRING
		}
	case '[':
		return TYPE_ARRAY
	case '{':
		return TYPE_MAP
	}

	return TYPE_UNKNOWN
}

func (obj *MarshalledObject) GetAsBool() (value bool, err error) {
	err = assertType(obj, TYPE_BOOL)
	if err != nil {
		return
	}

	value, _ = parseBool(obj.data)

	return
}

func (obj *MarshalledObject) GetAsInteger() (value int, err error) {
	err = assertType(obj, TYPE_INTEGER)
	if err != nil {
		return
	}

	value, _ = parseInt(obj.data[1:])

	return
}

func (obj *MarshalledObject) GetAsFloat() (value float64, err error) {
	err = assertType(obj, TYPE_FLOAT)
	if err != nil {
		return
	}

	str, _ := parseString(obj.data[1:])
	value, err = strconv.ParseFloat(str, 64)

	return
}

func (obj *MarshalledObject) GetAsString() (value string, err error) {
	err = assertType(obj, TYPE_STRING)
	if err != nil {
		return
	}

	if obj.data[0] == ':' {
		value, _ = parseString(obj.data[1:])
	} else {
		value, _ = parseString(obj.data[2:])
	}

	return
}

func assertType(obj *MarshalledObject, expected_type marshalledObjectType) (err error) {
	if obj.GetType() != expected_type {
		err = TypeMismatch
	}

	return
}

func parseBool(data []byte) (bool, int) {
	return data[0] == 'T', 1
}

func parseInt(data []byte) (int, int) {
	if data[0] > 0x05 && data[0] < 0xfb {
		value := int(data[0])

		if value > 0x7f {
			return -(0xff ^ value + 1) + 5, 1
		} else {
			return value - 5, 1
		}
	} else if data[0] <= 0x05 {
		value := 0
		i := data[0]

		for ; i > 0; i-- {
			value = value<<8 + int(data[i])
		}

		return value, int(data[0])
	} else {
		value := 0
		i := 0xff - data[0] + 1

		for ; i > 0; i-- {
			value = value<<8 + (0xff - int(data[i]))
		}

		return -(value + 1), int(0xff - data[0] + 1)
	}
}

func parseString(data []byte) (string, int) {
	length, size := parseInt(data)
	value := string(data[size : length+size])

	if len(data) > length+size+1 && data[length+size+1] == ':' {
		enc_symbol, enc_size := parseString(data[length+size+2:])

		if enc_symbol == "E" {
			_, enc_name_size := parseBool(data[length+size+enc_size+1:])
			enc_size += enc_name_size
		} else {
			_, enc_name_size := parseString(data[length+size+enc_size+3:])
			enc_size += enc_name_size
		}

		size += enc_size
	}

	return value, length + size
}
