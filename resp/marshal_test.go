package resp

import "testing"

func TestMarshal(t *testing.T) {
	tests := []struct {
		name string
		val  Value
		want string
	}{
		{
			name: "Simple String",
			val:  Value{Type: TypeString, Str: "PONG"},
			want: "+PONG\r\n",
		},
		{
			name: "Integer zero",
			val:  Value{Type: TypeInt, Num: 0},
			want: ":0\r\n",
		},
		{
			name: "Bulk String with special chars",
			val:  Value{Type: TypeBulk, Str: "line1\nline2"},
			want: "$11\r\nline1\nline2\r\n",
		},
		{
			name: "Complex Array",
			val: Value{
				Type: TypeArray,
				Arr: []Value{
					{Type: TypeBulk, Str: "GET"},
					{Type: TypeBulk, Str: "key"},
				},
			},
			want: "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.val.Marshal())
			if got != tt.want {
				t.Errorf("Marshal() = %q, want %q", got, tt.want)
			}
		})
	}
}
