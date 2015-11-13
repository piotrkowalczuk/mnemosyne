package main

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
)

// SessionEntity ...
type SessionEntity struct {
	Token    *mnemosyne.Token `json:"token"`
	Data     Data             `json:"data"`
	ExpireAt *time.Time       `json:"expireAt"`
}

func (se *SessionEntity) ExpireAtTimestamp() *mnemosyne.Timestamp {
	return mnemosyne.TimeToTimestamp(*se.ExpireAt)
}

func newSessionEntityFromSession(session *mnemosyne.Session) (*SessionEntity, error) {
	expireAtTime := session.ExpireAtTime()

	return &SessionEntity{
		Token:    session.Token,
		Data:     session.Data,
		ExpireAt: &expireAtTime,
	}, nil
}

func newSessionFromSessionEntity(entity *SessionEntity) *mnemosyne.Session {
	expireAtTimestamp := entity.ExpireAtTimestamp()

	return &mnemosyne.Session{
		Token:    entity.Token,
		Data:     entity.Data,
		ExpireAt: expireAtTimestamp,
	}
}
