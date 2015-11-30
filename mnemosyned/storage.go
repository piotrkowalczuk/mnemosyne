package main

import (
	"errors"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/stretchr/testify/mock"
)

const (
	storageEngineInMemory = "in_memory"
	storageEnginePostgres = "postgres"
	storageEngineRedis    = "redis"
)

var (
	errSessionNotFound = errors.New("mnemosyned: session not found")
)

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error

	Start(string, map[string]string) (*mnemosyne.Session, error)
	Abandon(*mnemosyne.Token) (bool, error)
	Get(*mnemosyne.Token) (*mnemosyne.Session, error)
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosyne.Session, error)
	Exists(*mnemosyne.Token) (bool, error)
	Delete(*mnemosyne.Token, *time.Time, *time.Time) (int64, error)

	SetValue(*mnemosyne.Token, string, string) (map[string]string, error)
	//	DeleteValue(*mnemosyne.Token, string) (*mnemosyne.Session, error)
	//	Clear(*mnemosyne.Token) (*mnemosyne.Session, error)
}

type storageMock struct {
	mock.Mock
}

// Start implements Storage interface.
func (sm *storageMock) Start(subjectID string, bag map[string]string) (*mnemosyne.Session, error) {
	args := sm.Called(subjectID, bag)

	ses, ok := args.Get(0).(*mnemosyne.Session)
	if !ok {
		return nil, args.Error(1)
	}
	return ses, args.Error(1)
}

// Ąbandon implements Storage interface.
func (sm *storageMock) Abandon(token *mnemosyne.Token) (bool, error) {
	args := sm.Called(token)

	return args.Bool(0), args.Error(1)
}

// Get implements Storage interface.
func (sm *storageMock) Get(token *mnemosyne.Token) (*mnemosyne.Session, error) {
	args := sm.Called(token)

	ses, ok := args.Get(0).(*mnemosyne.Session)
	if !ok {
		return nil, args.Error(1)
	}
	return ses, args.Error(1)
}

// List implements Storage interface.
func (sm *storageMock) List(offset, limit int64, expireAtFrom, expireAtTo *time.Time) ([]*mnemosyne.Session, error) {
	args := sm.Called(offset, limit, expireAtFrom, expireAtTo)

	ses, ok := args.Get(0).([]*mnemosyne.Session)
	if !ok {
		return nil, args.Error(1)
	}
	return ses, args.Error(1)
}

// Exists implements Storage interface.
func (sm *storageMock) Exists(token *mnemosyne.Token) (bool, error) {
	args := sm.Called(token)

	return args.Bool(0), args.Error(1)
}

// Delete implements Storage interface.
func (sm *storageMock) Delete(token *mnemosyne.Token, expireAtFrom, expireAtTo *time.Time) (int64, error) {
	args := sm.Called(token, expireAtFrom, expireAtTo)

	return args.Get(0).(int64), args.Error(1)
}

// SetValue implements Storage interface.
func (sm *storageMock) SetValue(token *mnemosyne.Token, key, value string) (map[string]string, error) {
	args := sm.Called(token, key, value)

	return args.Get(0).(map[string]string), args.Error(1)
}

// Setup implements Storage
func (sm *storageMock) Setup() error {
	return sm.Called().Error(0)
}

// Teardown implements Storage
func (sm *storageMock) TearDown() error {
	return sm.Called().Error(0)
}
