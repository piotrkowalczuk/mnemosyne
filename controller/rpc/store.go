package rpc

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-soa/mnemosyne/lib"
)

type Store struct {
	Logger  *logrus.Logger
	Storage lib.SessionStorage
}

func (s *Store) Get(id lib.SessionID, session *lib.Session) error {
	ses, err := s.Storage.Get(id)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	*session = *ses

	return err
}

func (s *Store) New(data lib.SessionData, session *lib.Session) error {
	ses, err := s.Storage.New(data)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	*session = *ses

	s.Logger.WithFields(logrus.Fields{
		"id":        ses.ID,
		"expire_at": ses.ExpireAt,
	}).Debug("New session has been created.")
	return err
}

func (s *Store) Exists(id lib.SessionID, exists *bool) error {
	ex, err := s.Storage.Exists(id)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	*exists = ex

	s.Logger.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Session existance has been checked.")
	return err
}

func (s *Store) Abandon(id lib.SessionID, ok *bool) error {
	err := s.Storage.Abandon(id)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	*ok = true

	s.Logger.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Session has been abandoned.")
	return nil
}

func (s *Store) SetData(entry lib.SessionDataEntry, session *lib.Session) error {
	ses, err := s.Storage.SetData(entry)
	if err != nil {
		s.Logger.Error(err)
		return err
	}

	*session = *ses

	s.Logger.WithFields(logrus.Fields{
		"id":    entry.ID,
		"key":   entry.Key,
		"value": entry.Value,
	}).Debug("Session data has been set.")
	return err
}
