package main

import (
	"testing"
	"time"

	"strconv"
	"sync"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	notExistsToken = &mnemosyne.Token{
		Hash: "NOT EXISTS",
	}
)

func testStorage_Create(t *testing.T, s Storage) {
	session, err := s.Create(Data{
		"username": "test",
	})

	if assert.NoError(t, err) {
		assert.Len(t, session.Token.Hash, 128)
		assert.Equal(t, session.Data, map[string]string{
			"username": "test",
		})
	}
}

func testStorage_Get(t *testing.T, s Storage) {
	new, err := s.Create(Data{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	got, err := s.Get(new.Token)
	require.NoError(t, err)
	assert.Equal(t, new.Token, got.Token)
	assert.Equal(t, new.Data, got.Data)
	assert.Equal(t, new.ExpireAt, got.ExpireAt)

	// Check for non existing Token
	got2, err2 := s.Get(notExistsToken)
	assert.Error(t, err2)
	assert.EqualError(t, err2, errSessionNotFound.Error())
	assert.Nil(t, got2)
}

func testStorage_List(t *testing.T, s Storage) {
	nb := 10
	key := "index"

	for i := 1; i <= nb; i++ {
		_, err := s.Create(Data{key: strconv.FormatInt(int64(i), 10)})
		require.NoError(t, err)
	}

	sessions, err := s.List(2, int64(nb), nil, nil)
	if assert.NoError(t, err) {
		assert.Len(t, sessions, nb)
		for i, s := range sessions {
			assert.NotEmpty(t, s.Token)
			assert.NotEmpty(t, s.ExpireAt)
			assert.Equal(t, s.Value(key), strconv.FormatInt(int64(i+1), 10))
		}
	}
}

func testStorage_Exists(t *testing.T, s Storage) {
	new, err := s.Create(Data{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	exists, err := s.Exists(new.Token)
	require.NoError(t, err)
	assert.True(t, exists)

	// Check for non existing Token
	exists2, err2 := s.Exists(notExistsToken)
	if assert.NoError(t, err2) {
		assert.False(t, exists2)
	}
}

func testStorage_Abandon(t *testing.T, s Storage) {
	new, err := s.Create(Data{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	ok2, err2 := s.Abandon(new.Token)
	assert.True(t, ok2)
	require.NoError(t, err2)

	// Check for already abondond session
	ok3, err3 := s.Abandon(new.Token)
	assert.False(t, ok3)
	assert.EqualError(t, err3, errSessionNotFound.Error())

	// Check for session that never exists
	ok4, err4 := s.Abandon(notExistsToken)
	assert.False(t, ok4)
	assert.EqualError(t, err4, errSessionNotFound.Error())
}

func testStorage_SetData(t *testing.T, s Storage) {
	new, err := s.Create(Data{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	got, err2 := s.SetData(new.Token, "email", "fake@email.com")
	require.NoError(t, err2)
	assert.Equal(t, new.Token, got.Token)
	assert.Equal(t, 2, len(got.Data))
	assert.Equal(t, "fake@email.com", got.Value("email"))
	assert.Equal(t, "test", got.Value("username"))
	assert.NotNil(t, got.ExpireAt)

	// Check for overwritten field
	got2, err2 := s.SetData(new.Token, "email", "morefakethanbefore@email.com")
	require.NoError(t, err2)
	assert.Equal(t, new.Token, got2.Token)
	assert.Equal(t, 2, len(got2.Data))
	assert.Equal(t, "morefakethanbefore@email.com", got2.Value("email"))
	assert.Equal(t, "test", got2.Value("username"))
	assert.NotNil(t, got2.ExpireAt)

	// Check for non existing Token
	got3, err3 := s.SetData(notExistsToken, "email", "fake@email.com")
	require.Error(t, err3, errSessionNotFound.Error())
	assert.Nil(t, got3)

	wg := sync.WaitGroup{}
	// Check for concurent access
	concurent := func(t *testing.T, wg *sync.WaitGroup, key, value string) {
		defer wg.Done()

		// Check for overwritten field
		_, err := s.SetData(new.Token, key, value)

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

	got4, err4 := s.Get(new.Token)
	if assert.NoError(t, err4) {
		assert.Equal(t, new.Token, got4.Token)
		assert.Equal(t, 22, len(got4.Data))
	}
}

func testStorage_Delete(t *testing.T, s Storage) {
	expiredAtTo := time.Now().Add(35 * time.Minute)

	affected, err := s.Delete(nil, nil, &expiredAtTo)
	if assert.NoError(t, err) {
		assert.Equal(t, int64(14), affected)
	}

	data := []struct {
		id            bool
		expiredAtFrom bool
		expiredAtTo   bool
	}{
		{
			id: true,
		},
		{
			expiredAtFrom: true,
		},
		{
			expiredAtTo: true,
		},
		{
			id:            true,
			expiredAtFrom: true,
		},
		{
			id:            true,
			expiredAtFrom: true,
			expiredAtTo:   true,
		},
		{
			expiredAtFrom: true,
			expiredAtTo:   true,
		},
		{
			id:          true,
			expiredAtTo: true,
		},
	}

DataLoop:
	for _, args := range data {
		new, err := s.Create(nil)
		require.NoError(t, err)

		if !assert.NoError(t, err) {
			continue DataLoop
		}

		var id *mnemosyne.Token
		var expiredAtTo *time.Time
		var expiredAtFrom *time.Time

		if args.id {
			id = new.Token
		}

		if args.expiredAtFrom {
			eaf := new.ExpireAtTime().Add(-29 * time.Minute)
			expiredAtFrom = &eaf
		}
		if args.expiredAtTo {
			eat := new.ExpireAtTime().Add(29 * time.Minute)
			expiredAtTo = &eat
		}

		affected, err = s.Delete(id, expiredAtFrom, expiredAtTo)
		if assert.NoError(t, err) {
			if assert.Equal(t, int64(1), affected, "one session should be removed for id: %-5t, expiredAtFrom: %-5t, expiredAtTo: %-5t", args.id, args.expiredAtFrom, args.expiredAtTo) {
				t.Logf("as expected session can be deleted with arguments id: %-5t, expiredAtFrom: %-5t, expiredAtTo: %-5t", args.id, args.expiredAtFrom, args.expiredAtTo)
			}
		}

		affected, err = s.Delete(id, expiredAtFrom, expiredAtTo)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(0), affected)
		}
	}
}
