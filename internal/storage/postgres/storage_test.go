package postgres_test

import (
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
)

func TestPostgresStorage_Start(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_Start(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Get(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_Get(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_List(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_List(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_List_between(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_List_between(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Exists(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_Exists(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Abandon(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_Abandon(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_SetValue(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_SetValue(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Delete(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	storage.TestStorage_Delete(t, s.store)

	s.teardown(t)
}
