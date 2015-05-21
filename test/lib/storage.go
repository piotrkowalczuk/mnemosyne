package lib

import (
	"testing"
	"time"

	"sync"

	"github.com/go-soa/mnemosyne/lib"
	"github.com/stretchr/testify/assert"
)

func testStorageNew(t *testing.T, storage lib.SessionStorage) {
	session, err := storage.New(lib.SessionData{
		"username": "test",
	})

	assert.NoError(t, err)
	assert.Len(t, session.ID, 128)
	assert.Equal(t, session.Data, lib.SessionData{
		"username": "test",
	})
}

func testStorageGet(t *testing.T, storage lib.SessionStorage) {
	new, err := storage.New(lib.SessionData{
		"username": "test",
	})
	assert.NoError(t, err)

	// Check for existing ID
	got, err := storage.Get(new.ID)
	assert.NoError(t, err)
	assert.Equal(t, new.ID, got.ID)
	assert.Equal(t, new.Data, got.Data)
	assert.True(t, new.ExpireAt.Equal(*got.ExpireAt))

	// Check for non existing ID
	got2, err2 := storage.Get("NOT EXISTS")
	assert.Error(t, err2)
	assert.EqualError(t, err2, lib.ErrSessionNotFound.Error())
	assert.Nil(t, got2)
}

func testStorageExists(t *testing.T, storage lib.SessionStorage) {
	new, err := storage.New(lib.SessionData{
		"username": "test",
	})
	assert.NoError(t, err)

	// Check for existing ID
	exists, err := storage.Exists(new.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check for non existing ID
	exists2, err2 := storage.Exists("NOT EXISTS")
	assert.NoError(t, err2)
	assert.False(t, exists2)
}

func testStorageAbandon(t *testing.T, storage lib.SessionStorage) {
	new, err := storage.New(lib.SessionData{
		"username": "test",
	})
	assert.NoError(t, err)

	// Check for existing ID
	err2 := storage.Abandon(new.ID)
	assert.NoError(t, err2)

	// Check for already abondond session
	err3 := storage.Abandon(new.ID)
	assert.EqualError(t, err3, lib.ErrSessionNotFound.Error())

	// Check for session that never exists
	err4 := storage.Abandon("NEVER EXISTS")
	assert.EqualError(t, err4, lib.ErrSessionNotFound.Error())
}

func testStorageSetData(t *testing.T, storage lib.SessionStorage) {
	new, err := storage.New(lib.SessionData{
		"username": "test",
	})
	assert.NoError(t, err)

	// Check for existing ID
	got, err2 := storage.SetData(lib.SessionDataEntry{
		ID:    new.ID,
		Key:   "email",
		Value: "fake@email.com",
	})
	assert.NoError(t, err2)
	assert.Equal(t, new.ID, got.ID)
	assert.Equal(t, 2, len(got.Data))
	assert.Equal(t, "fake@email.com", got.Data.Get("email"))
	assert.Equal(t, "test", got.Data.Get("username"))
	assert.NotNil(t, got.ExpireAt)

	// Check for overwritten field
	got2, err2 := storage.SetData(lib.SessionDataEntry{
		ID:    new.ID,
		Key:   "email",
		Value: "morefakethanbefore@email.com",
	})
	assert.NoError(t, err2)
	assert.Equal(t, new.ID, got2.ID)
	assert.Equal(t, 2, len(got2.Data))
	assert.Equal(t, "morefakethanbefore@email.com", got2.Data.Get("email"))
	assert.Equal(t, "test", got2.Data.Get("username"))
	assert.NotNil(t, got2.ExpireAt)

	// Check for non existing ID
	got3, err3 := storage.SetData(lib.SessionDataEntry{
		ID:    "NOT EXISTS",
		Key:   "email",
		Value: "fake@email.com",
	})
	assert.Error(t, err3, lib.ErrSessionNotFound.Error())
	assert.Nil(t, got3)

	wg := sync.WaitGroup{}
	// Check for concurent access
	concurent := func(t *testing.T, wg *sync.WaitGroup, key, value string) {
		defer wg.Done()

		// Check for overwritten field
		_, err := storage.SetData(lib.SessionDataEntry{
			ID:    new.ID,
			Key:   key,
			Value: value,
		})

		assert.NoError(t, err)
	}

	wg.Add(20)
	go concurent(t, &wg, "k1", "v1")
	go concurent(t, &wg, "k2", "v2")
	go concurent(t, &wg, "k3", "v3")
	go concurent(t, &wg, "k4", "v4")
	go concurent(t, &wg, "k5", "v5")
	go concurent(t, &wg, "k6", "v6")
	go concurent(t, &wg, "k7", "v7")
	go concurent(t, &wg, "k8", "v8")
	go concurent(t, &wg, "k9", "v9")
	go concurent(t, &wg, "k10", "v10")
	go concurent(t, &wg, "k11", "v11")
	go concurent(t, &wg, "k12", "v12")
	go concurent(t, &wg, "k13", "v13")
	go concurent(t, &wg, "k14", "v14")
	go concurent(t, &wg, "k15", "v15")
	go concurent(t, &wg, "k16", "v16")
	go concurent(t, &wg, "k17", "v17")
	go concurent(t, &wg, "k18", "v18")
	go concurent(t, &wg, "k19", "v19")
	go concurent(t, &wg, "k20", "v20")

	wg.Wait()

	got4, err4 := storage.Get(new.ID)
	assert.NoError(t, err4)
	assert.Equal(t, new.ID, got4.ID)
	assert.Equal(t, 22, len(got4.Data))
}

func testStorageCleanup(t *testing.T, storage lib.SessionStorage) {
	till := time.Now().Add(35 * time.Minute)

	affected, err := storage.Cleanup(&till)
	if assert.NoError(t, err) {
		assert.Equal(t, int64(4), affected)
	}
}
