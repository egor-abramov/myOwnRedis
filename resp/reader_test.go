package resp

import (
	"reflect"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expect    Value
		expectErr bool
	}{
		{
			name:   "Bulk string normal",
			input:  "$5\r\nhello\r\n",
			expect: Value{Type: TypeBulk, Str: "hello"},
		},
		{
			name:   "Bulk string empty",
			input:  "$0\r\n\r\n",
			expect: Value{Type: TypeBulk, Str: ""},
		},
		{
			name:   "Bulk string Null",
			input:  "$-1\r\n",
			expect: Value{Type: TypeBulk, Str: ""},
		},

		{
			name:   "Empty Array",
			input:  "*0\r\n",
			expect: Value{Type: TypeArray, Arr: []Value{}},
		},
		{
			name:  "Simple Array of Bulk strings",
			input: "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n",
			expect: Value{
				Type: TypeArray,
				Arr: []Value{
					{Type: TypeBulk, Str: "ECHO"},
					{Type: TypeBulk, Str: "hello"},
				},
			},
		},
		{
			name:  "Nested Arrays (Array of Arrays)",
			input: "*2\r\n*1\r\n$1\r\n1\r\n*1\r\n$1\r\n2\r\n",
			expect: Value{
				Type: TypeArray,
				Arr: []Value{
					{Type: TypeArray, Arr: []Value{{Type: TypeBulk, Str: "1"}}},
					{Type: TypeArray, Arr: []Value{{Type: TypeBulk, Str: "2"}}},
				},
			},
		},

		{
			name:      "Invalid type Prefix",
			input:     "!unknown\r\n",
			expectErr: true,
		},
		{
			name:      "Truncated Bulk string",
			input:     "$5\r\nhel",
			expectErr: true,
		},
		{
			name:      "Invalid Integer format",
			input:     ":abc\r\n",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewResp(strings.NewReader(tt.input))
			got, err := r.Read()

			if (err != nil) != tt.expectErr {
				t.Errorf("Read() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("Read() got = %#v, expect %#v", got, tt.expect)
			}
		})
	}
}
