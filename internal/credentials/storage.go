package credentials

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type InMemoryStore struct {
	mu          sync.RWMutex
	credentials map[string]*Credential
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{credentials: make(map[string]*Credential)}
}

func (s *InMemoryStore) Save(cred *Credential) error {
	s.mu.Lock()

	defer s.mu.Unlock()
	if cred == nil || cred.ID == "" {
		return errors.New("credential or ID is empty")
	}
	s.credentials[cred.ID] = cred
	return nil
}

type ErrNotFound struct{ ID string }

func (e ErrNotFound) Error() string { return "credential not found: " + e.ID }

func (s *InMemoryStore) Get(id string) (*Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cred, ok := s.credentials[id]
	if !ok {
		return nil, ErrNotFound{ID: id}
	}
	return cred, nil
}

type FileStore struct {
	Dir string
}

func (fs *FileStore) Save(cred *Credential) error {
	if cred == nil || cred.ID == "" {
		return errors.New("credential or ID is empty")
	}

	data, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		return err
	}

	filename := fs.Dir + "/" + cred.ID + ".json"
	return os.WriteFile(filename, data, 0644)
}

func (fs *FileStore) Get(id string) (*Credential, error) {
	filename := fs.Dir + "/" + id + ".json"
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cred Credential
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}
