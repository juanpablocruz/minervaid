package credentials

import (
	"os"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	cred := NewCredential("id1", "issuer1", map[string]interface{}{"foo": "bar"})
	if err := store.Save(cred); err != nil {
		t.Fatalf("Error guardando en memoria: %v", err)
	}
	got, err := store.Get("id1")
	if err != nil {
		t.Fatalf("Error recuperando de memoria: %v", err)
	}
	if got.ID != "id1" {
		t.Errorf("ID recuperado = %s; se esperaba id1", got.ID)
	}
}

func TestFileStore(t *testing.T) {
	dir, err := os.MkdirTemp("", "credstore")
	if err != nil {
		t.Fatalf("Error creando temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	fs := &FileStore{Dir: dir}
	cred := NewCredential("id2", "issuer2", map[string]interface{}{"baz": "qux"})
	if err := fs.Save(cred); err != nil {
		t.Fatalf("Error guardando en fichero: %v", err)
	}
	got, err := fs.Get("id2")
	if err != nil {
		t.Fatalf("Error recuperando de fichero: %v", err)
	}
	if got.ID != "id2" {
		t.Errorf("ID recuperado = %s; se esperaba id2", got.ID)
	}
}
