package marshal

import (
	"testing"
)

func TestCreateMarshalledObject(t *testing.T) {
	m := CreateMarshalledObject([]byte{4, 8, 1})

	if m.MajorVersion != 4 {
		t.Errorf("CreateMarshalledObject created an object with Marshal major version set to %d instead of 4", m.MajorVersion)
	}

	if m.MinorVersion != 8 {
		t.Errorf("CreateMarshalledObject created an object with Marshal minor version set to %d instead of 8", m.MinorVersion)
	}
}

type getTypeTestCase struct {
	Data        []byte
	Expectation marshalledObjectType
}

func TestGetType(t *testing.T) {
	marshalledObjectTypeNames := []string{"unknown", "nil", "boolean", "integer", "float", "array", "map"}

	tests := []getTypeTestCase{
		// Nil
		{[]byte{4, 8, 48}, TYPE_NIL},
		// Booleans
		{[]byte{4, 8, 70}, TYPE_BOOLEAN}, // false
		{[]byte{4, 8, 84}, TYPE_BOOLEAN}, // true
		// Integers
		{[]byte{4, 8, 105, 0}, TYPE_INTEGER},                 // 0
		{[]byte{4, 8, 105, 6}, TYPE_INTEGER},                 // 1
		{[]byte{4, 8, 105, 250}, TYPE_INTEGER},               // -1
		{[]byte{4, 8, 105, 3, 64, 226, 1}, TYPE_INTEGER},     // 123456
		{[]byte{4, 8, 105, 253, 192, 29, 254}, TYPE_INTEGER}, // -123456
		// Floats
		{[]byte{4, 8, 102, 6, 48}, TYPE_FLOAT},                               // 0.0
		{[]byte{4, 8, 102, 8, 49, 46, 53}, TYPE_FLOAT},                       // 1.5
		{[]byte{4, 8, 102, 9, 45, 49, 46, 53}, TYPE_FLOAT},                   // -1.5
		{[]byte{4, 8, 102, 12, 49, 46, 50, 53, 101, 51, 48}, TYPE_FLOAT},     // 1.25e30
		{[]byte{4, 8, 102, 13, 49, 46, 50, 53, 101, 45, 51, 48}, TYPE_FLOAT}, // 1.25e-30
		// Arrays
		{[]byte{4, 8, 91, 0}, TYPE_ARRAY},                                             // []
		{[]byte{4, 8, 91, 6, 73, 34, 8, 102, 111, 111, 6, 58, 6, 69, 84}, TYPE_ARRAY}, // ["foo"]
		// Maps (Ruby hashes)
		{[]byte{4, 8, 123, 0}, TYPE_MAP},                                                                 // {}
		{[]byte{4, 8, 123, 6, 58, 8, 102, 111, 111, 73, 34, 8, 98, 97, 114, 6, 58, 6, 69, 84}, TYPE_MAP}, // {foo: "bar"}
	}

	for _, testCase := range tests {
		object_type := CreateMarshalledObject(testCase.Data).GetType()
		if object_type != testCase.Expectation {
			t.Errorf("GetType() returned '%s' instead of '%s'", marshalledObjectTypeNames[int(object_type)], marshalledObjectTypeNames[testCase.Expectation])
		}
	}
}

type getAsBooleanTestCase struct {
	Data        []byte
	Expectation bool
}

func TestGetAsBoolean(t *testing.T) {
	tests := []getAsBooleanTestCase{
		{[]byte{4, 8, 70}, false},
		{[]byte{4, 8, 84}, true},
	}

	value, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsBoolean() // should return an error
	if err == nil {
		t.Error("GetAsBoolean() returned no error when attempted to typecast nil to boolean")
	}

	for _, testCase := range tests {
		value, err = CreateMarshalledObject(testCase.Data).GetAsBoolean()

		if err != nil {
			t.Errorf("GetAsBoolean() returned an error: '%s' for %t", err.Error(), testCase.Expectation)
		}

		if value != testCase.Expectation {
			t.Errorf("GetAsBoolean() returned '%t' instead of '%t'", value, testCase.Expectation)
		}
	}
}

type getAsIntegerTestCase struct {
	Data        []byte
	Expectation int
}

func TestGetAsInteger(t *testing.T) {
	tests := []getAsIntegerTestCase{
		{[]byte{4, 8, 0x69, 0x00}, 0},
		{[]byte{4, 8, 0x69, 0x06}, 1},
		{[]byte{4, 8, 0x69, 0x7f}, 122},
		{[]byte{4, 8, 0x69, 0x01, 0x7b}, 123},
		{[]byte{4, 8, 0x69, 0x02, 0x00, 0x01}, 256},
		{[]byte{4, 8, 0x69, 0x04, 0xff, 0xff, 0xff, 0x3f}, (2 << 29) - 1},
		{[]byte{4, 8, 0x69, 0xfa}, -1},
		{[]byte{4, 8, 0x69, 0xff, 0x84}, -124},
		{[]byte{4, 8, 0x69, 0xfe, 0xff, 0xfe}, -257},
		{[]byte{4, 8, 0x69, 0xfc, 0x00, 0x00, 0x00, 0xc0}, -(2 << 29)},
	}

	value, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsInteger() // should return an error
	if err == nil {
		t.Error("GetAsInteger() returned no error when attempted to typecast nil to boolean")
	}

	for _, testCase := range tests {
		value, err = CreateMarshalledObject(testCase.Data).GetAsInteger()

		if err != nil {
			t.Errorf("GetAsInteger() returned an error: '%s' for %d", err.Error(), testCase.Expectation)
		}

		if value != testCase.Expectation {
			t.Errorf("GetAsInteger() returned '%d' instead of '%d'", value, testCase.Expectation)
		}
	}
}
