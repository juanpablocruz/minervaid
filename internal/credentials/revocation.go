package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type RevocationList struct {
	mu      sync.RWMutex
	store   string
	revoked map[string]struct{}
}

func NewRevocationList(path string) (*RevocationList, error) {
	rl := &RevocationList{store: path, revoked: make(map[string]struct{})}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(data) > 0 {
		var list []string
		if err := json.Unmarshal(data, &list); err != nil {
			return nil, err
		}
		for _, id := range list {
			rl.revoked[id] = struct{}{}
		}
	}
	return rl, nil
}

func (rl *RevocationList) Revoke(id string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if _, exists := rl.revoked[id]; exists {
		return fmt.Errorf("credential %s already revoked", id)
	}
	rl.revoked[id] = struct{}{}
	return rl.save()
}

func (rl *RevocationList) IsRevoked(id string) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	_, revoked := rl.revoked[id]
	return revoked
}

func (rl *RevocationList) List() []string {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	ids := make([]string, 0, len(rl.revoked))
	for id := range rl.revoked {
		ids = append(ids, id)
	}
	return ids
}

func (rl *RevocationList) save() error {
	// snapshot map directly (write lock held by Revoke)
	list := make([]string, 0, len(rl.revoked))
	for id := range rl.revoked {
		list = append(list, id)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(rl.store, data, 0644)
}
