package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Storage struct {
	mu         sync.RWMutex
	Data       map[string]string `json:"data"`
	Expiration map[string]int64  `json:"expiration"`
	path       string
}

func NewStorage(path string) *Storage {
	s := &Storage{Data: make(map[string]string),
		Expiration: make(map[string]int64),
		path:       path,
	}
	go s.saveDump()
	return s
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expired, ok := s.Expiration[key]
	if !ok || expired > time.Now().UnixNano() {
		val, exists := s.Data[key]
		return val, exists
	}
	return "", false
}

func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	s.Data[key] = value
	s.mu.Unlock()
}

func (s *Storage) SetTTL(key string, expiration int64) {
	s.mu.Lock()
	s.Expiration[key] = expiration
	s.mu.Unlock()
}

func (s *Storage) Delete(key string) bool {
	s.mu.Lock()
	_, ok := s.Data[key]
	delete(s.Data, key)
	delete(s.Expiration, key)
	s.mu.Unlock()
	return ok
}

func (s *Storage) Load() error {
	file, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(file, s)
}

func (s *Storage) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	tempPath := s.path + ".tmp"
	err = os.WriteFile(tempPath, data, 0644)
	if err != nil {
		return err
	}
	return os.Rename(tempPath, s.path)
}

func (s *Storage) saveDump() {
	ticker := time.NewTicker(3 * time.Minute)
	for range ticker.C {
		<-ticker.C
		s.Save()
	}
}
