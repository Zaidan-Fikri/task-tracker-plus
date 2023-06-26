package repository

import (
	"a21hc3NpZ25tZW50/model"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) *sessionsRepo {
	return &sessionsRepo{db}
}

func (s *sessionsRepo) AddSessions(session model.Session) error {
	result := s.db.Create(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *sessionsRepo) DeleteSession(token string) error {
	result := s.db.Delete(&model.Session{}, "token = ?", token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *sessionsRepo) UpdateSessions(session model.Session) error {
	result := s.db.Model(&model.Session{}).Where("email = ?", session.Email).Update("token", session.Token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *sessionsRepo) SessionAvailEmail(email string) (model.Session, error) {
	var session model.Session
	result := s.db.Where("email = ?", email).First(&session)
	if result.Error != nil {
		return model.Session{}, result.Error
	}
	return session, nil
}

func (s *sessionsRepo) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	result := s.db.Where("token = ?", token).First(&session)
	if result.Error != nil {
		return model.Session{}, result.Error
	}
	return session, nil
}

func (s *sessionsRepo) TokenValidity(token string) (model.Session, error) {
	session, err := s.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if s.TokenExpired(session) {
		err := s.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, err
	}

	return session, nil
}

func (s *sessionsRepo) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
