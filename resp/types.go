package resp

import "strconv"

const (
	stringPrefix = '+'
	errorPrefix  = '-'
	intPrefix    = ':'
	arrayPrefix  = '*'
	bulkPrefix   = '$'

	TypeString = "string"
	TypeInt    = "int"
	TypeArray  = "array"
	TypeBulk   = "bulk"
	TypeError  = "error"
	TypeNull   = "null"
)

type Value struct {
	Type string
	Str  string
	Num  int
	Arr  []Value
}

func (v Value) Marshal() []byte {
	switch v.Type {
	case TypeString:
		return v.marshalString()
	case TypeInt:
		return v.marshalInt()
	case TypeError:
		return v.marshalError()
	case TypeNull:
		return v.marshalNull()
	case TypeArray:
		return v.marshalArray()
	case TypeBulk:
		return v.marshalBulk()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, stringPrefix)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalInt() []byte {
	var bytes []byte
	bytes = append(bytes, intPrefix)
	bytes = append(bytes, strconv.Itoa(v.Num)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, errorPrefix)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, bulkPrefix)
	bytes = append(bytes, strconv.Itoa(len(v.Str))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, arrayPrefix)
	bytes = append(bytes, strconv.Itoa(len(v.Arr))...)
	bytes = append(bytes, '\r', '\n')
	for _, val := range v.Arr {
		bytes = append(bytes, val.Marshal()...)
	}
	return bytes
}
