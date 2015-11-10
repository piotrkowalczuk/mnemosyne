package main

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
)

// SessionEntity ...
type SessionEntity struct {
	ID       *mnemosyne.ID `json:"id"`
	Data     Data          `json:"data"`
	ExpireAt *time.Time    `json:"expireAt"`
}

func (se *SessionEntity) ExpireAtTimestamp() *mnemosyne.Timestamp {
	return mnemosyne.TimeToTimestamp(*se.ExpireAt)
}

func newSessionEntityFromSession(session *mnemosyne.Session) (*SessionEntity, error) {
	expireAtTime := session.ExpireAtTime()

	return &SessionEntity{
		ID:       session.Id,
		Data:     session.Data,
		ExpireAt: &expireAtTime,
	}, nil
}

func newSessionFromSessionEntity(entity *SessionEntity) *mnemosyne.Session {
	expireAtTimestamp := entity.ExpireAtTimestamp()

	return &mnemosyne.Session{
		Id:       entity.ID,
		Data:     entity.Data,
		ExpireAt: expireAtTimestamp,
	}
}
