package postgres_test

import (
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
)

func TestPostgresStorage_Start(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageStart(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Get(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageGet(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_List(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageList(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_List_between(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageListBetween(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Exists(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageExists(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Abandon(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageAbandon(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_SetValue(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageSetValue(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Delete(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorageDelete(t, s.store)

	s.teardown(t)
}
