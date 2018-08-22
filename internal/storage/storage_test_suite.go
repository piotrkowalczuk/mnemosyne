package storage

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageStart(t *testing.T, s Storage) {
	t.Helper()

	subjectID := "subjectID"
	subjectClient := "subjectClient"
	bag := map[string]string{
		"username": "test",
	}
	session, err := s.Start(context.Background(), randomToken(t), "", subjectID, subjectClient, bag)

	if assert.NoError(t, err) {
		assert.Len(t, session.AccessToken, 128)
		assert.Equal(t, subjectID, session.SubjectId)
		assert.Equal(t, bag, session.Bag)
	}
}

func TestStorageGet(t *testing.T, s Storage) {
	t.Helper()

	ses, err := s.Start(context.Background(), randomToken(t), randomToken(t), "subjectID", "subjectClient", map[string]string{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	got, err := s.Get(context.Background(), ses.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, ses.AccessToken, got.AccessToken)
	assert.Equal(t, ses.RefreshToken, got.RefreshToken)
	assert.Equal(t, ses.Bag, got.Bag)
	if ses.ExpireAt.Seconds > got.ExpireAt.Seconds || (ses.ExpireAt.Seconds == got.ExpireAt.Seconds && ses.ExpireAt.Nanos > got.ExpireAt.Nanos) {
		t.Fatalf("after get expire at should be increased, got %s but expected %s", got.ExpireAt, ses.ExpireAt)
	}

	// Check for non existing Token
	got2, err2 := s.Get(context.Background(), "keyhash")
	assert.Error(t, err2)
	assert.Equal(t, err2, ErrSessionNotFound)
	assert.Nil(t, got2)
}

func TestStorageList(t *testing.T, s Storage) {
	nb := 10
	key := "index"
	sid := "subjectID"
	sc := "subjectClient"

	for i := 1; i <= nb; i++ {
		_, err := s.Start(context.Background(), randomToken(t), randomToken(t), sid, sc, map[string]string{key: strconv.FormatInt(int64(i), 10)})
		if err != nil {
			t.Fatalf("unexpected error on session start: %s", err.Error())
		}
	}

	sessions, err := s.List(context.Background(), 2, int64(nb), nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	_, err = s.List(context.Background(), 2, 0, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(sessions) != nb-2 {
		t.Fatalf("wrong number of sessions returned: expected %d but got %d", nb-2, len(sessions))
	}

	for i, s := range sessions {
		assert.NotEmpty(t, s.AccessToken)
		assert.NotEmpty(t, s.RefreshToken)
		assert.NotEmpty(t, s.ExpireAt)
		assert.Equal(t, s.SubjectId, sid)

		assert.Equal(t, s.Bag[key], strconv.FormatInt(int64(i+3), 10))
	}
}

func TestStorageListBetween(t *testing.T, s Storage) {
	nb := 10
	key := "index"
	sid := "subjectID"
	sc := "subjectClient"
	var (
		err      error
		from, to time.Time
	)

	for i := 1; i <= nb; i++ {
		res, err := s.Start(context.Background(), randomToken(t), "", sid, sc, map[string]string{key: strconv.FormatInt(int64(i), 10)})
		if err != nil {
			t.Fatalf("unexpected error on session start: %s", err.Error())
		}
		if i == 1 {
			if from, err = ptypes.Timestamp(res.ExpireAt); err != nil {
				t.Fatalf("timestamp conversion unexpected error: %s", err.Error())
			}
			from = from.Add(-1 * time.Second)
		}
		if i == nb {
			if to, err = ptypes.Timestamp(res.ExpireAt); err != nil {
				t.Fatalf("timestamp conversion unexpected error: %s", err.Error())
			}
			to = to.Add(1 * time.Second)
		}
	}

	sessions, err := s.List(context.Background(), 0, int64(nb), &from, &to)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if len(sessions) != nb {
		t.Fatalf("wrong number of sessions returned: expected %d but got %d", nb, len(sessions))
	}
}

func TestStorageExists(t *testing.T, s Storage) {
	ses, err := s.Start(context.Background(), randomToken(t), "", "subjectID", "subjectClient", map[string]string{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	exists, err := s.Exists(context.Background(), ses.AccessToken)
	require.NoError(t, err)
	assert.True(t, exists)

	// Check for non existing Token
	exists2, err2 := s.Exists(context.Background(), "keyhash")
	if assert.NoError(t, err2) {
		assert.False(t, exists2)
	}
}

func TestStorageAbandon(t *testing.T, s Storage) {
	ses, err := s.Start(context.Background(), randomToken(t), "", "subjectID", "subjectClient", map[string]string{
		"username": "test",
	})
	require.NoError(t, err)

	// Check for existing Token
	ok2, err2 := s.Abandon(context.Background(), ses.AccessToken)
	assert.True(t, ok2)
	require.NoError(t, err2)

	// Check for already abandoned session
	ok3, err3 := s.Abandon(context.Background(), ses.AccessToken)
	assert.False(t, ok3)
	assert.Equal(t, err3, ErrSessionNotFound)

	// Check for session that never exists
	ok4, err4 := s.Abandon(context.Background(), "keyhash")
	assert.False(t, ok4)
	assert.Equal(t, err4, ErrSessionNotFound)
}

func TestStorageSetValue(t *testing.T, s Storage) {
	ses, err := s.Start(context.Background(), randomToken(t), "", "subjectID", "subjectClient", map[string]string{
		"username": "test",
	})
	if err != nil {
		t.Fatalf("unexpected error on session start: %s", err.Error())
	}
	if s == nil {
		t.Fatal("storage is nil")
	}

	// Check for existing Token
	got, err2 := s.SetValue(context.Background(), ses.AccessToken, "email", "fake@email.com")
	require.NoError(t, err2)
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "fake@email.com", got["email"])
	assert.Equal(t, "test", got["username"])

	// Check for overwritten field
	bag2, err2 := s.SetValue(context.Background(), ses.AccessToken, "email", "morefakethanbefore@email.com")
	require.NoError(t, err2)
	assert.Equal(t, 2, len(bag2))
	assert.Equal(t, "morefakethanbefore@email.com", bag2["email"])
	assert.Equal(t, "test", bag2["username"])

	// Check for non existing Token
	bag3, err3 := s.SetValue(context.Background(), "keyhash", "email", "fake@email.com")
	assert.Equal(t, err3, ErrSessionNotFound)
	assert.Nil(t, bag3)

	wg := sync.WaitGroup{}
	// Check for concurrent access
	concurrent := func(t *testing.T, wg *sync.WaitGroup, key, value string) {
		defer wg.Done()

		// Check for overwritten field
		_, err := s.SetValue(context.Background(), ses.AccessToken, key, value)

		assert.NoError(t, err)
	}

	wg.Add(20)
	go concurrent(t, &wg, "k1", "v1")
	go concurrent(t, &wg, "k2", "v2")
	go concurrent(t, &wg, "k3", "v3")
	go concurrent(t, &wg, "k4", "v4")
	go concurrent(t, &wg, "k5", "v5")
	go concurrent(t, &wg, "k6", "v6")
	go concurrent(t, &wg, "k7", "v7")
	go concurrent(t, &wg, "k8", "v8")
	go concurrent(t, &wg, "k9", "v9")
	go concurrent(t, &wg, "k10", "v10")
	go concurrent(t, &wg, "k11", "v11")
	go concurrent(t, &wg, "k12", "v12")
	go concurrent(t, &wg, "k13", "v13")
	go concurrent(t, &wg, "k14", "v14")
	go concurrent(t, &wg, "k15", "v15")
	go concurrent(t, &wg, "k16", "v16")
	go concurrent(t, &wg, "k17", "v17")
	go concurrent(t, &wg, "k18", "v18")
	go concurrent(t, &wg, "k19", "v19")
	go concurrent(t, &wg, "k20", "v20")

	wg.Wait()

	got4, err4 := s.Get(context.Background(), ses.AccessToken)
	if assert.NoError(t, err4) {
		assert.Equal(t, ses.AccessToken, got4.AccessToken)
		assert.Equal(t, 22, len(got4.Bag))
	}
}

func TestStorageDelete(t *testing.T, s Storage) {
	nb := int64(10)
	key := "index"
	sid := "subjectID"
	sc := "subjectClient"

	for i := int64(1); i <= nb; i++ {
		_, err := s.Start(context.Background(), randomToken(t), "", sid, sc, map[string]string{key: strconv.FormatInt(i, 10)})
		if err != nil {
			t.Fatalf("unexpected error on session start: %s", err.Error())
		}
	}

	expiredAtFrom := time.Now().Add(-35 * time.Minute)
	expiredAtTo := time.Now().Add(35 * time.Minute)

	affected, err := s.Delete(context.Background(), "", "", "", &expiredAtFrom, &expiredAtTo)
	if assert.NoError(t, err) {
		assert.Equal(t, nb, affected)
	}

	_, err = s.Delete(context.Background(), "", "", "", nil, nil)
	assert.Error(t, err)

	data := []struct {
		subjectID     bool
		accessToken   bool
		refreshToken  bool
		expiredAtFrom bool
		expiredAtTo   bool
	}{
		{
			subjectID: true,
		},
		{
			accessToken: true,
		},
		{
			refreshToken: true,
		},
		{
			expiredAtFrom: true,
		},
		{
			expiredAtTo: true,
		},
		{
			accessToken:   true,
			expiredAtFrom: true,
		},
		{
			subjectID:     true,
			accessToken:   true,
			expiredAtFrom: true,
		},
		{
			accessToken:   true,
			refreshToken:  true,
			expiredAtFrom: true,
		},
		{
			accessToken:   true,
			expiredAtFrom: true,
			expiredAtTo:   true,
		},
		{
			accessToken:   true,
			refreshToken:  true,
			expiredAtFrom: true,
			expiredAtTo:   true,
		},
		{
			expiredAtFrom: true,
			expiredAtTo:   true,
		},
		{
			refreshToken:  true,
			expiredAtFrom: true,
			expiredAtTo:   true,
		},
		{
			accessToken: true,
			expiredAtTo: true,
		},
		{
			accessToken:  true,
			refreshToken: true,
			expiredAtTo:  true,
		},
	}

DataLoop:
	for _, args := range data {
		ses, err := s.Start(context.Background(), randomToken(t), randomToken(t), "subjectID", "subjectID", nil)
		require.NoError(t, err)

		if !assert.NoError(t, err) {
			continue DataLoop
		}

		var (
			subjectID, accessToken, refreshToken string
			expiredAtTo                          *time.Time
			expiredAtFrom                        *time.Time
		)

		if args.subjectID {
			subjectID = ses.SubjectId
		}
		if args.accessToken {
			accessToken = ses.AccessToken
		}
		if args.refreshToken {
			refreshToken = ses.RefreshToken
		}
		if args.expiredAtFrom {
			expireAtFrom, err := ptypes.Timestamp(ses.ExpireAt)
			if assert.NoError(t, err) {
				continue DataLoop
			}
			eaf := expireAtFrom.Add(-29 * time.Minute)
			expiredAtFrom = &eaf
		}
		if args.expiredAtTo {
			expireAtTo, err := ptypes.Timestamp(ses.ExpireAt)
			if assert.NoError(t, err) {
				continue DataLoop
			}
			eat := expireAtTo.Add(29 * time.Minute)
			expiredAtTo = &eat
		}

		affected, err = s.Delete(context.Background(), subjectID, accessToken, refreshToken, expiredAtFrom, expiredAtTo)
		if assert.NoError(t, err) {
			if assert.Equal(t, int64(1), affected, "one session should be removed for accessToken: %-5t, refreshToken: %-5t, ,expiredAtFrom: %-5t, expiredAtTo: %-5t", args.accessToken, args.refreshToken, args.expiredAtFrom, args.expiredAtTo) {
				t.Logf("as expected session can be deleted with arguments accessToken: %-5t, refreshToken: %-5t, expiredAtFrom: %-5t, expiredAtTo: %-5t", args.accessToken, args.refreshToken, args.expiredAtFrom, args.expiredAtTo)
			}
		}

		affected, err = s.Delete(context.Background(), subjectID, accessToken, refreshToken, expiredAtFrom, expiredAtTo)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(0), affected)
		}
	}
}

func randomToken(t *testing.T) string {
	at, err := mnemosyne.RandomAccessToken()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	return at
}
