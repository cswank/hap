package hap

import (
	"errors"
	"testing"
)

type failingStore struct {
	Store
	getErr error
	setErr error
}

func (st failingStore) Get(key string) ([]byte, error) {
	if st.getErr != nil {
		return nil, st.getErr
	}
	return st.Store.Get(key)
}

func (st failingStore) Set(key string, value []byte) error {
	if st.setErr != nil {
		return st.setErr
	}
	return st.Store.Set(key, value)
}

func TestMigrateMissingSchema(t *testing.T) {
	st := &storer{NewMemStore()}
	if err := migrate(st); err != nil {
		t.Fatal(err)
	}

	if s, err := st.GetString("schema"); err != nil || s != "1" {
		t.Fatalf("schema = %q, %v", s, err)
	}
}

func TestMigrateGetError(t *testing.T) {
	errRead := errors.New("read failed")
	st := &storer{failingStore{Store: NewMemStore(), getErr: errRead}}
	if err := migrate(st); !errors.Is(err, errRead) {
		t.Fatalf("expected %v, got %v", errRead, err)
	}
}

func TestMigrateSetError(t *testing.T) {
	errWrite := errors.New("write failed")
	st := &storer{failingStore{Store: NewMemStore(), setErr: errWrite}}
	if err := migrate(st); !errors.Is(err, errWrite) {
		t.Fatalf("expected %v, got %v", errWrite, err)
	}
}
