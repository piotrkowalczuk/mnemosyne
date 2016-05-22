package mnemosyned

import "testing"

func TestPostgresStorage_Start(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_Start(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Get(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_Get(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_List(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_List(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_List_between(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_List_between(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Exists(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_Exists(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Abandon(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_Abandon(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_SetValue(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_SetValue(t, s.store)

	s.teardown(t)
}

func TestPostgresStorage_Delete(t *testing.T) {
	s := &postgresSuite{}
	s.setup(t)

	testStorage_Delete(t, s.store)

	s.teardown(t)
}
