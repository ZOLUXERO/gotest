package infrastructure

import "fmt"

type UserService struct {
	repo    UserRepository
	emitter KafkaEventEmitter
}

func (s *UserService) RedeemPoints(userID string, points int, reason string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	if user.Points < points {
		return fmt.Errorf("not enough points")
	}
	user.Points -= points
	err = s.repo.Save(user)
	if err != nil {
		return err
	}
	event := &core.PointsEvent{
		UserID:    user.ID,
		Points:    -points,
		Operation: "redeem",
		Reason:    reason,
	}
	return s.emitter.Emit(event)
}

func (s *UserService) AccumulatePoints(userID string, points int, reason string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.Points += points
	err = s.repo.Save(user)
	if err != nil {
		return err
	}
	event := &core.PointsEvent{
		UserID:    user.ID,
		Points:    points,
		Operation: "accumulate",
		Reason:    reason,
	}
	return s.emitter.Emit(event)
}

func (s *UserService) GetPoints(userID string) (int, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, fmt.Errorf("user not found")
	}
	return user.Points, nil
}
