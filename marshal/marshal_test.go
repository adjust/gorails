package marshal

import (
	"testing"
	"reflect"
)

func TestCreateMarshalledObject(t *testing.T) {
	m := CreateMarshalledObject([]byte{4, 8, 1})

	if m.MajorVersion != 4 {
		t.Errorf("CreateMarshalledObject created an object with Marshal major version set to %v instead of 4", m.MajorVersion)
	}

	if m.MinorVersion != 8 {
		t.Errorf("CreateMarshalledObject created an object with Marshal minor version set to %v instead of 8", m.MinorVersion)
	}
}

type getTypeTestCase struct {
	Data        []byte
	Expectation marshalledObjectType
}

func TestGetType(t *testing.T) {
	marshalledObjectTypeNames := []string{"unknown", "nil", "bool", "integer", "float", "string", "array", "map"}

	tests := []getTypeTestCase{
		// Nil
		{[]byte{4, 8, 48}, TYPE_NIL},
		// Booleans
		{[]byte{4, 8, 70}, TYPE_BOOL}, // false
		{[]byte{4, 8, 84}, TYPE_BOOL}, // true
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
		// Strings
		{[]byte{4, 8, 73, 34, 0, 6, 58, 6, 69, 84}, TYPE_STRING},                                                           // ''
		{[]byte{4, 8, 58, 10, 104, 101, 108, 108, 111}, TYPE_STRING},                                                       // :hello
		{[]byte{4, 8, 73, 34, 17, 72, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100, 6, 58, 6, 69, 84}, TYPE_STRING}, // 'Hello, world'
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
			t.Errorf("GetType() returned '%v' instead of '%v'", marshalledObjectTypeNames[int(object_type)], marshalledObjectTypeNames[testCase.Expectation])
		}
	}
}

type getAsBoolTestCase struct {
	Data        []byte
	Expectation bool
}

func TestGetAsBool(t *testing.T) {
	tests := []getAsBoolTestCase{
		{[]byte{4, 8, 70}, false},
		{[]byte{4, 8, 84}, true},
	}

	value, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsBool() // should return an error
	if err == nil {
		t.Error("GetAsBool() returned no error when attempted to typecast nil to boolean")
	}

	for _, testCase := range tests {
		value, err = CreateMarshalledObject(testCase.Data).GetAsBool()

		if err != nil {
			t.Errorf("GetAsBool() returned an error: '%v' for %v", err.Error(), testCase.Expectation)
		}

		if value != testCase.Expectation {
			t.Errorf("GetAsBool() returned '%v' instead of '%v'", value, testCase.Expectation)
		}
	}
}

type getAsIntegerTestCase struct {
	Data        []byte
	Expectation int64
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
		t.Error("GetAsInteger() returned no error when attempted to typecast nil to int")
	}

	for _, testCase := range tests {
		value, err = CreateMarshalledObject(testCase.Data).GetAsInteger()

		if err != nil {
			t.Errorf("GetAsInteger() returned an error: '%v' for %v", err.Error(), testCase.Expectation)
		}

		if value != testCase.Expectation {
			t.Errorf("GetAsInteger() returned '%v' instead of '%v'", value, testCase.Expectation)
		}
	}
}

type getAsFloatTestCase struct {
	Data        []byte
	Expectation float64
}

func TestGetAsFloat(t *testing.T) {
	tests := []getAsFloatTestCase{
		{[]byte{4, 8, 102, 6, 48}, 0.0},
		{[]byte{4, 8, 102, 13, 49, 46, 52, 51, 101, 45, 49, 48}, 1.43e-10},
		{[]byte{4, 8, 102, 13, 49, 46, 52, 51, 101, 45, 49, 48}, 1.43E-10},
		{[]byte{4, 8, 102, 10, 48, 46, 49, 50, 53}, 0.125},
		{[]byte{4, 8, 102, 10, 49, 50, 46, 53, 54}, 12.56},
		{[]byte{4, 8, 102, 12, 49, 46, 52, 51, 101, 49, 48}, 1.43e+10},
		{[]byte{4, 8, 102, 12, 49, 46, 52, 51, 101, 49, 48}, 1.43E+10},
		{[]byte{4, 8, 102, 14, 45, 49, 46, 52, 51, 101, 45, 49, 48}, -1.43e-10},
		{[]byte{4, 8, 102, 14, 45, 49, 46, 52, 51, 101, 45, 49, 48}, -1.43E-10},
		{[]byte{4, 8, 102, 11, 45, 48, 46, 49, 50, 53}, -0.125},
		{[]byte{4, 8, 102, 11, 45, 49, 50, 46, 53, 54}, -12.56},
		{[]byte{4, 8, 102, 13, 45, 49, 46, 52, 51, 101, 49, 48}, -1.43e+10},
		{[]byte{4, 8, 102, 13, 45, 49, 46, 52, 51, 101, 49, 48}, -1.43E+10},
	}

	value, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsFloat() // should return an error
	if err == nil {
		t.Error("GetAsFloat() returned no error when attempted to typecast nil to float")
	}

	for _, testCase := range tests {
		value, err = CreateMarshalledObject(testCase.Data).GetAsFloat()

		if err != nil {
			t.Errorf("GetAsFloat() returned an error: '%v' for %v", err.Error(), testCase.Expectation)
		}

		if value != testCase.Expectation {
			t.Errorf("GetAsFloat() returned '%v' instead of '%v'", value, testCase.Expectation)
		}
	}
}

type getAsStringTestCase struct {
	Data        []byte
	Expectation string
}

func TestGetAsString(t *testing.T) {
	tests := []getAsStringTestCase{
		{[]byte{4, 8, 73, 34, 0, 6, 58, 6, 69, 84}, ""},                                                                       // ''
		{[]byte{4, 8, 58, 10, 104, 101, 108, 108, 111}, "hello"},                                                              // :hello
		{[]byte{4, 8, 73, 34, 17, 72, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100, 6, 58, 6, 69, 84}, "Hello, world"}, // 'Hello, world'
	}

	value, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsString() // should return an error
	if err == nil {
		t.Error("GetAsString() returned no error when attempted to typecast nil to string")
	}

	for _, testCase := range tests {
		value, err = CreateMarshalledObject(testCase.Data).GetAsString()

		if err != nil {
			t.Errorf("GetAsString() returned an error: '%v' for %v", err.Error(), testCase.Expectation)
		}

		if value != testCase.Expectation {
			t.Errorf("GetAsString() returned '%v' instead of '%v'", value, testCase.Expectation)
		}
	}
}

type getAsArrayOfIntsTestCase struct {
	Data        []byte
	Expectation []int64
}

type getAsArrayOfStringsTestCase struct {
	Data        []byte
	Expectation []string
}

func TestGetAsArray(t *testing.T) {
	int_tests := []getAsArrayOfIntsTestCase{
		{[]byte{4, 8, 91, 0}, []int64{}},
		{[]byte{4, 8, 91, 10, 105, 255, 0, 105, 250, 105, 0, 105, 6, 105, 2, 0, 1}, []int64{-256, -1, 0, 1, 256}},
	}

	_, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsArray() // should return an error
	if err == nil {
		t.Error("GetAsArray() returned no error when attempted to typecast nil to array")
	}

	for _, testCase := range int_tests {
		value, err := CreateMarshalledObject(testCase.Data).GetAsArray()

		if err != nil {
			t.Errorf("GetAsArray() returned an error: '%v' for %v", err.Error(), testCase.Expectation)
		}

		if len(value) != len(testCase.Expectation) {
			t.Errorf("GetAsArray() returned an array with length %d for %v", len(value), testCase.Expectation)
		} else {
			for i, v := range value {
				value, err := v.GetAsInteger()

				if err != nil {
					t.Errorf("GetAsArray() returned an error '%v' for element #%d (%d) of %v", err.Error(), i, testCase.Expectation[i], testCase.Expectation)
				}

				if value != testCase.Expectation[i] {
					t.Errorf("GetAsArray() returned '%v' instead of '%v'", value, testCase.Expectation)
				}
			}
		}
	}

	string_tests := []getAsArrayOfStringsTestCase{
		{[]byte{4, 8, 91, 6, 73, 34, 8, 102, 111, 111, 6, 58, 6, 69, 84}, []string{"foo"}}, // ["foo"]
		{[]byte{4, 8, 91, 6, 58, 8, 98, 97, 114}, []string{"bar"}}, // [:bar]
		{[]byte{4, 8, 91, 8, 73, 34, 8, 102, 111, 111, 6, 58, 6, 69, 84, 73, 34, 8, 98, 97, 114, 6, 59, 0, 84, 58, 8, 98, 97, 122}, []string{"foo", "bar", "baz"}}, // ["foo", "bar", :baz]
		{[]byte{4, 8, 91, 8, 73, 34, 8, 102, 111, 111, 6, 58, 6, 69, 84, 73, 34, 8, 98, 97, 114, 6, 58, 13, 101, 110, 99, 111, 100, 105, 110, 103, 34, 14, 83, 104, 105, 102, 116, 95, 74, 73, 83, 58, 8, 98, 97, 122}, []string{"foo", "bar", "baz"}}, // ["foo", "bar".force_encoding("SHIFT_JIS"), :baz]
		{[]byte{4, 8, 91, 7, 73, 34, 6, 120, 6, 58, 6, 69, 84, 64, 6}, []string{"x", "x"}},
	}

	for _, testCase := range string_tests {
		value, err := CreateMarshalledObject(testCase.Data).GetAsArray()

		if err != nil {
			t.Errorf("GetAsArray() returned an error: '%v' for %v", err.Error(), testCase.Expectation)
		}

		if len(value) != len(testCase.Expectation) {
			t.Errorf("GetAsArray() returned an array with length %d for %v", len(value), testCase.Expectation)
		} else {
			for i, v := range value {
				value, err := v.GetAsString()

				if err != nil {
					t.Errorf("GetAsArray() returned an error '%v' for element #%d (%v %d) of %v", err.Error(), i, testCase.Expectation[i], v.GetType(), testCase.Expectation)
				}

				if value != testCase.Expectation[i] {
					t.Errorf("GetAsArray() returned '%v' instead of '%v'", value, testCase.Expectation)
				}
			}
		}
	}
}

type getAsMapOfIntsTestCase struct {
	Data        []byte
	Expectation map[string]int64
}

type getAsMapOfStringsTestCase struct {
	Data        []byte
	Expectation map[string]string
}

type getAsMapOfMapsTestCase struct {
	Data        []byte
	Expectation map[string]map[string]int64
}

func TestGetAsMap(t *testing.T) {
	int_tests := []getAsMapOfIntsTestCase{
		{
			[]byte{4, 8, 123, 0},
			map[string]int64{},
		},
		{
			[]byte{4, 8, 123, 12, 73, 34, 6, 48, 6, 58, 6, 69, 84, 105, 0, 105, 6, 105, 6, 105, 250, 105, 250, 48, 105, 255, 0, 73, 34, 8, 102, 111, 111, 6, 59, 0, 84, 105, 2, 0, 1, 73, 34, 8, 98, 97, 114, 6, 58, 13, 101, 110, 99, 111, 100, 105, 110, 103, 34, 14, 83, 104, 105, 102, 116, 95, 74, 73, 83, 105, 2, 188, 2, 58, 8, 98, 97, 122, 105, 254, 68, 253},
			map[string]int64{
				"0":     0,
				"1":     1,
				"-1":    -1,
				"<nil>": -256,
				"foo":   256,
				"bar":   700,
				"baz":   -700,
			},
		},
	}

	_, err := CreateMarshalledObject([]byte{4, 8, 48}).GetAsMap() // should return an error
	if err == nil {
		t.Error("GetAsMap() returned no error when attempted to typecast nil to map")
	}

	for _, testCase := range int_tests {
		value, err := CreateMarshalledObject(testCase.Data).GetAsMap()

		m := make(map[string]int64)
		for k, v := range value {
			m[k], err = v.GetAsInteger()

			if err != nil {
				t.Errorf("GetAsMap() returned an error while parsing %s", k)
			}
		}

		if ! reflect.DeepEqual(m, testCase.Expectation) {
			t.Errorf("%v is not equal %v", m, testCase.Expectation)
		}
	}

	string_tests := []getAsMapOfStringsTestCase{
		{
			[]byte{4, 8, 123, 12, 73, 34, 6, 48, 6, 58, 6, 69, 84, 73, 34, 6, 48, 6, 59, 0, 84, 105, 6, 73, 34, 6, 49, 6, 59, 0, 84, 105, 250, 73, 34, 0, 6, 59, 0, 84, 48, 73, 34, 8, 102, 111, 111, 6, 59, 0, 84, 73, 34, 8, 102, 111, 111, 6, 59, 0, 84, 73, 34, 8, 98, 97, 114, 6, 58, 13, 101, 110, 99, 111, 100, 105, 110, 103, 34, 14, 83, 104, 105, 102, 116, 95, 74, 73, 83, 73, 34, 8, 98, 97, 114, 6, 59, 0, 84, 58, 8, 98, 97, 122, 59, 7, 73, 34, 6, 48, 6, 59, 0, 84},
			map[string]string{
				"0":     "0",   // "0" => "0"
				"1":     "1",   // 1 => "1"
				"-1":    "",    // -1 => ""
				"<nil>": "foo", // nil => "foo"
				"foo":   "bar", // "foo" => "bar".force_encoding("SHIFT_JIS")
				"bar":   "baz", // "bar".force_encoding("SHIFT_JIS") => :baz
				"baz":   "0",   // :baz => "0"
			},
		},
		{
			[]byte{4, 8, 123, 8, 58, 6, 97, 73, 34, 6, 120, 6, 58, 6, 69, 84, 58, 6, 98, 64, 6, 58, 6, 99, 64, 6},
			map[string]string{
				"a":  "x",
				"b":  "x",
				"c":  "x",
			},
		},
	}

	for _, testCase := range string_tests {
		value, err := CreateMarshalledObject(testCase.Data).GetAsMap()

		m := make(map[string]string)
		for k, v := range value {
			m[k], err = v.GetAsString()

			if err != nil {
				t.Errorf("GetAsMap() returned an error while parsing %s %d: %s", k, v.GetType(), err.Error())
			}
		}

		if ! reflect.DeepEqual(m, testCase.Expectation) {
			t.Errorf("%v is not equal %v", m, testCase.Expectation)
		}
	}

	map_tests := []getAsMapOfMapsTestCase{
		{
			[]byte{4, 8, 123, 8, 58, 6, 97, 123, 6, 73, 34, 6, 120, 6, 58, 6, 69, 84, 105, 6, 58, 6, 98, 64, 6, 58, 6, 99, 64, 6},
			map[string]map[string]int64{
				"a":  map[string]int64{"x": 1},
				"b":  map[string]int64{"x": 1},
				"c":  map[string]int64{"x": 1},
			},
		},
	}

	for _, testCase := range map_tests {
		value, _ := CreateMarshalledObject(testCase.Data).GetAsMap()

		m := make(map[string]map[string]int64)
		for k, v := range value {
			vv, err := v.GetAsMap()

			if err != nil {
				t.Errorf("GetAsMap() returned an error while parsing %v", v)
			}

			m2 := make(map[string]int64)
			for k2, v2 := range vv {
				m2[k2], err = v2.GetAsInteger()

				if err != nil {
					t.Errorf("GetAsInteger() returned an error while parsing %v", v2)
				}
			}

			m[k] = m2

			if err != nil {
				t.Errorf("GetAsMap() returned an error while parsing %s %d: %s", k, v.GetType(), err.Error())
			}
		}

		if ! reflect.DeepEqual(m, testCase.Expectation) {
			t.Errorf("%v is not equal %v", m, testCase.Expectation)
		}
	}
}
