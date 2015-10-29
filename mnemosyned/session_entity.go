package main

import (
	"time"

	"github.com/piotrkowalczuk/mnemosyne/shared"
)

// SessionEntity ...
type SessionEntity struct {
	ID       *shared.ID `json:"id"`
	Data     Data       `json:"data"`
	ExpireAt *time.Time `json:"expireAt"`
}

func newSessionEntityFromSession(session *shared.Session) (*SessionEntity, error) {
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

func newSessionFromSessionEntity(entity *SessionEntity) *shared.Session {
	return &shared.Session{
		Id:       entity.ID,
		Data:     entity.Data,
		ExpireAt: entity.ExpireAt.Format(time.RFC3339),
	}
}
