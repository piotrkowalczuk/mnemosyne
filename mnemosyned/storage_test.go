package mnemosyned

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testStorage_Start(t *testing.T, s Storage) {
	subjectID := "subjectID"
	bag := map[string]string{
		"username": "test",
	}
	session, err := s.Start(subjectID, bag)

	if assert.NoError(t, err) {
		assert.Len(t, session.AccessToken.Hash, 128)
		assert.Equal(t, subjectID, session.SubjectId)
		assert.Equal(t, bag, session.Bag)
	}
}

func testStorage_Get(t *testing.T, s Storage) {
	ses, err := s.Start("subjectID", map[string]string{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	got, err := s.Get(ses.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, ses.AccessToken, got.AccessToken)
	assert.Equal(t, ses.Bag, got.Bag)
	assert.Equal(t, ses.ExpireAt, got.ExpireAt)

	// Check for non existing Token
	got2, err2 := s.Get(&mnemosyne.AccessToken{Key: []byte("key"), Hash: []byte("hash")})
	assert.Error(t, err2)
	assert.EqualError(t, err2, SessionNotFound.Error())
	assert.Nil(t, got2)
}

func testStorage_List(t *testing.T, s Storage) {
	nb := 10
	key := "index"
	sid := "subjectID"

	for i := 1; i <= nb; i++ {
		_, err := s.Start(sid, map[string]string{key: strconv.FormatInt(int64(i), 10)})
		if err != nil {
			t.Fatalf("unexpected error on session start: %s", err.Error())
		}
	}

	sessions, err := s.List(2, int64(nb), nil, nil)
	if len(sessions) != nb-2 {
		t.Fatalf("wrong number of sessions returned: expected %d but got %d", nb, len(sessions))
	}

	if assert.NoError(t, err) {
		assert.Len(t, sessions, nb-2)
		for i, s := range sessions {
			assert.NotEmpty(t, s.AccessToken)
			assert.NotEmpty(t, s.ExpireAt)
			assert.Equal(t, s.SubjectId, sid)

			assert.Equal(t, s.Bag[key], strconv.FormatInt(int64(i+3), 10))
		}
	}
}

func testStorage_Exists(t *testing.T, s Storage) {
	new, err := s.Start("subjectID", map[string]string{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	exists, err := s.Exists(new.AccessToken)
	require.NoError(t, err)
	assert.True(t, exists)

	// Check for non existing Token
	exists2, err2 := s.Exists(&mnemosyne.AccessToken{Key: []byte("key"), Hash: []byte("hash")})
	if assert.NoError(t, err2) {
		assert.False(t, exists2)
	}
}

func testStorage_Abandon(t *testing.T, s Storage) {
	new, err := s.Start("subjectID", map[string]string{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	ok2, err2 := s.Abandon(new.AccessToken)
	assert.True(t, ok2)
	require.NoError(t, err2)

	// Check for already abandoned session
	ok3, err3 := s.Abandon(new.AccessToken)
	assert.False(t, ok3)
	assert.EqualError(t, err3, SessionNotFound.Error())

	// Check for session that never exists
	ok4, err4 := s.Abandon(&mnemosyne.AccessToken{Key: []byte("key"), Hash: []byte("hash")})
	assert.False(t, ok4)
	assert.EqualError(t, err4, SessionNotFound.Error())
}

func testStorage_SetValue(t *testing.T, s Storage) {
	new, err := s.Start("subjectID", map[string]string{
		"username": "test",
	})
	if err != nil {
		t.Fatalf("unexpected error on session start: %s", err.Error())
	}
	if s == nil {
		t.Fatalf("storage is nil")
	}

	// Check for existing Token
	got, err2 := s.SetValue(new.AccessToken, "email", "fake@email.com")
	require.NoError(t, err2)
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "fake@email.com", got["email"])
	assert.Equal(t, "test", got["username"])

	// Check for overwritten field
	bag2, err2 := s.SetValue(new.AccessToken, "email", "morefakethanbefore@email.com")
	require.NoError(t, err2)
	assert.Equal(t, 2, len(bag2))
	assert.Equal(t, "morefakethanbefore@email.com", bag2["email"])
	assert.Equal(t, "test", bag2["username"])

	// Check for non existing Token
	bag3, err3 := s.SetValue(&mnemosyne.AccessToken{Key: []byte("key"), Hash: []byte("hash")}, "email", "fake@email.com")
	require.Error(t, err3, SessionNotFound.Error())
	assert.Nil(t, bag3)

	wg := sync.WaitGroup{}
	// Check for concurent access
	concurent := func(t *testing.T, wg *sync.WaitGroup, key, value string) {
		defer wg.Done()

		// Check for overwritten field
		_, err := s.SetValue(new.AccessToken, key, value)

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

	got4, err4 := s.Get(new.AccessToken)
	if assert.NoError(t, err4) {
		assert.Equal(t, new.AccessToken, got4.AccessToken)
		assert.Equal(t, 22, len(got4.Bag))
	}
}

func testStorage_Delete(t *testing.T, s Storage) {
	nb := int64(10)
	key := "index"
	sid := "subjectID"

	for i := int64(1); i <= nb; i++ {
		_, err := s.Start(sid, map[string]string{key: strconv.FormatInt(i, 10)})
		if err != nil {
			t.Fatalf("unexpected error on session start: %s", err.Error())
		}
	}

	expiredAtTo := time.Now().Add(35 * time.Minute)

	affected, err := s.Delete(nil, nil, &expiredAtTo)
	if assert.NoError(t, err) {
		assert.Equal(t, nb, affected)
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
		new, err := s.Start("subjectID", nil)
		require.NoError(t, err)

		if !assert.NoError(t, err) {
			continue DataLoop
		}

		var (
			id            *mnemosyne.AccessToken
			expiredAtTo   *time.Time
			expiredAtFrom *time.Time
		)

		if args.id {
			id = new.AccessToken
		}

		if args.expiredAtFrom {
			expireAtFrom, err := ptypes.Timestamp(new.ExpireAt)
			if assert.NoError(t, err) {
				continue DataLoop
			}
			eaf := expireAtFrom.Add(-29 * time.Minute)
			expiredAtFrom = &eaf
		}
		if args.expiredAtTo {
			expireAtTo, err := ptypes.Timestamp(new.ExpireAt)
			if assert.NoError(t, err) {
				continue DataLoop
			}
			eat := expireAtTo.Add(29 * time.Minute)
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
