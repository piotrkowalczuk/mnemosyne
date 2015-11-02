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

func newSessionEntityFromSession(session *mnemosyne.Session) (*SessionEntity, error) {
	et, err := time.Parse(time.RFC3339, session.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &SessionEntity{
		ID:       session.Id,
		Data:     session.Data,
		ExpireAt: &et,
	}, nil
}

func newSessionFromSessionEntity(entity *SessionEntity) *mnemosyne.Session {
	return &mnemosyne.Session{
		Id:       entity.ID,
		Data:     entity.Data,
		ExpireAt: entity.ExpireAt.Format(time.RFC3339),
	}
}
