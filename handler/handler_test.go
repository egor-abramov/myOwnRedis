package handler

import (
	"myOwnRedis/resp"
	"testing"
)

type MockStorage struct {
	data map[string]string
}

func (m *MockStorage) SetTTL(string, int64) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStorage) Get(key string) (string, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *MockStorage) Set(key, value string) {
	m.data[key] = value
}

func (m *MockStorage) Delete(key string) bool {
	_, ok := m.data[key]
	delete(m.data, key)
	return ok
}

func TestHandler_Handle(t *testing.T) {
	storage := &MockStorage{data: make(map[string]string)}
	h := NewHandler(storage)

	tests := []struct {
		name     string
		input    resp.Value
		expected resp.Value
	}{
		{
			name: "PING without args",
			input: resp.Value{
				Type: resp.TypeArray,
				Arr:  []resp.Value{{Type: resp.TypeString, Str: "PING"}},
			},
			expected: resp.Value{Type: resp.TypeString, Str: "PONG"},
		},
		{
			name: "PING with arg",
			input: resp.Value{
				Type: resp.TypeArray,
				Arr:  []resp.Value{{Type: resp.TypeString, Str: "PING"}, {Type: resp.TypeString, Str: "hello"}},
			},
			expected: resp.Value{Type: resp.TypeString, Str: "hello"},
		},
		{
			name: "SET command",
			input: resp.Value{
				Type: resp.TypeArray,
				Arr: []resp.Value{
					{Type: resp.TypeString, Str: "SET"},
					{Type: resp.TypeString, Str: "key1"},
					{Type: resp.TypeString, Str: "val1"},
				},
			},
			expected: resp.Value{Type: resp.TypeString, Str: "OK"},
		},
		{
			name: "GET existing key",
			input: resp.Value{
				Type: resp.TypeArray,
				Arr: []resp.Value{
					{Type: resp.TypeString, Str: "GET"},
					{Type: resp.TypeString, Str: "key1"},
				},
			},
			expected: resp.Value{Type: resp.TypeString, Str: "val1"},
		},
		{
			name: "DELETE existing key",
			input: resp.Value{
				Type: resp.TypeArray,
				Arr: []resp.Value{
					{Type: resp.TypeString, Str: "DEL"},
					{Type: resp.TypeString, Str: "key1"},
				},
			},
			expected: resp.Value{Type: resp.TypeString, Str: "1"},
		},
		{
			name: "Unknown command",
			input: resp.Value{
				Type: resp.TypeArray,
				Arr:  []resp.Value{{Type: resp.TypeString, Str: "UNKNOWN"}},
			},
			expected: resp.Value{Type: resp.TypeError, Str: "ERR unknown command"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := h.Handle(tt.input)
			if res.Type != tt.expected.Type || res.Str != tt.expected.Str {
				t.Errorf("expected %v, got %v", tt.expected, res)
			}
		})
	}
}
