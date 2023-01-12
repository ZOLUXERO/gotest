package infrastructure

import (
	"fmt"

	"github.com/ZOLUXERO/gotest/core"
)

type UserService struct {
	Repo    core.UserRepository
	Emitter EventEmitter
}

func (s *UserService) SubtractPoints(userID string, points int, reason string) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("usuario no encontrado")
	}
	if user.Points < points {
		return fmt.Errorf("no tiene suficientes puntos")
	}
	user.Points -= points
	err = s.Repo.Save(user)
	if err != nil {
		return err
	}
	event := &core.PointsEvent{
		UserID:    user.ID,
		Points:    -points,
		Operation: "redeem",
		Reason:    reason,
	}
	return s.Emitter.Emit(event)
}

func (s *UserService) AddPoints(userID string, points int, reason string) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("usuario no encontrado")
	}
	user.Points += points
	err = s.Repo.Save(user)
	if err != nil {
		return err
	}
	event := &core.PointsEvent{
		UserID:    user.ID,
		Points:    points,
		Operation: "accumulate",
		Reason:    reason,
	}
	return s.Emitter.Emit(event)
}

func (s *UserService) GetPoints(userID string) (int, error) {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, fmt.Errorf("usuario no encontrado")
	}
	return user.Points, nil
}
