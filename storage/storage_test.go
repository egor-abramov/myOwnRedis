package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	// Инициализируем хранилище
	s := &Storage{Data: make(map[string]string)}

	s.Set("key1", "value1")
	val, ok := s.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected value1, got %s", val)
	}

	s.Set("key1", "new_value")
	val, _ = s.Get("key1")
	if val != "new_value" {
		t.Errorf("expected new_value, got %s", val)
	}

	_, ok = s.Get("missing")
	if ok {
		t.Error("expected ok=false for missing key")
	}

	deleted := s.Delete("key1")
	if !deleted {
		t.Error("expected Delete to return true for existing key")
	}
	_, ok = s.Get("key1")
	if ok {
		t.Error("key should be deleted")
	}
}
