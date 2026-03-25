package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}

func (r *Resp) Read() (v Value, err error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return
	}
	v.Type = string(typ)
	switch typ {
	case arrayPrefix:
		return r.readArray()
	case bulkPrefix:
		return r.readBulk()
	default:
		return v, fmt.Errorf("unrecognized resp type: %s", v.Type)
	}
}

func (r *Resp) readLine() (string, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) >= 2 {
		line = line[:len(line)-2]
	}
	return line, nil
}

func (r *Resp) readInt() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(line)
}

func (r *Resp) readBulk() (Value, error) {
	n, err := r.readInt()
	if err != nil {
		return Value{}, err
	}

	if n < 0 {
		return Value{Type: TypeBulk}, nil
	}

	bytes := make([]byte, n)
	_, err = r.reader.Read(bytes)
	if err != nil {
		return Value{}, err
	}
	_, err = r.readLine()
	if err != nil {
		return Value{}, err
	}
	v := Value{Type: TypeBulk}
	v.Str = string(bytes)
	return v, nil
}

func (r *Resp) readArray() (Value, error) {
	n, err := r.readInt()
	if err != nil {
		return Value{}, err
	}
	v := Value{Type: TypeArray, Arr: []Value{}}
	for i := 0; i < n; i++ {
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		v.Arr = append(v.Arr, val)
	}
	return v, nil
}
