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
var UnknownObjectType = errors.New("gorails/marshal: unknown marshalled object type")

const (
  TYPE_UNKNOWN marshalledObjectType = 0
  TYPE_MAP     marshalledObjectType = 1
  TYPE_ARRAY   marshalledObjectType = 2
  TYPE_BOOLEAN marshalledObjectType = 3
  TYPE_INTEGER marshalledObjectType = 4
  TYPE_FLOAT   marshalledObjectType = 5
  TYPE_NIL     marshalledObjectType = 6
)

func CreateMarshalledObject(serialized_data []byte) (*MarshalledObject) {
  return &(MarshalledObject{serialized_data[0], serialized_data[1], serialized_data[2:]})
}

func (obj *MarshalledObject) GetType() (object_type marshalledObjectType, err error) {
  object_type = TYPE_UNKNOWN

  switch (obj.data[0]) {
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
    err = UnknownObjectType
  }

  return
}
