package infrastructure

import (
	"testing"

	"github.com/ZOLUXERO/gotest/core"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
	users map[string]*core.User
}

func (m *MockUserRepository) GetByID(id string) (*core.User, error) {
	args := m.Called(id)
	return args.Get(0).(*core.User), args.Error(1)
}

func (m *MockUserRepository) Save(user *core.User) error {
	args := m.Called(user)
	m.users[user.ID] = user
	return args.Error(0)
}

func (m *MockUserRepository) SubtractPoints(user *core.User, points int64) error {
	args := m.Called(user, points)
	return args.Error(0)
}

func (m *MockUserRepository) AddPoints(user *core.User, points int64) error {
	args := m.Called(user, points)
	return args.Error(0)
}

type MockKafkaEventEmitter struct {
	mock.Mock
}

func (m *MockKafkaEventEmitter) Emit(event *core.PointsEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestUserService_AddPoints(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockRepo.On("GetByID", "user-123").Return(&core.User{ID: "user-123", Points: 100}, nil)
	mockRepo.On("Save", &core.User{ID: "user-123", Points: 150}).Return(nil)

	mockKafkaEmitter := &MockKafkaEventEmitter{}
	mockKafkaEmitter.On("Emit", &kafka.Message{}).Return(nil)

	service := &UserService{
		Repo:    mockRepo,
		Emitter: mockKafkaEmitter,
	}

	err := service.AddPoints("user-123", 50, "bonus")
	assert.Nil(t, err)
	assert.Equal(t, 150, mockRepo.users["user-123"].Points)
	mockRepo.AssertExpectations(t)
	mockKafkaEmitter.AssertExpectations(t)
}

func TestUserService_SubtractPoints(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockRepo.On("GetByID", "user-123").Return(&core.User{ID: "user-123", Points: 100}, nil)
	mockRepo.On("Save", &core.User{ID: "user-123", Points: 50}).Return(nil)

	service := &UserService{
		Repo: mockRepo,
	}

	err := service.SubtractPoints("user-123", 50, "test")
	assert.Nil(t, err)
	assert.Equal(t, 50, mockRepo.users["user-123"].Points)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetPoints(t *testing.T) {
	mockRepo := &MockUserRepository{}
	mockRepo.On("GetByID", "user-123").Return(&core.User{ID: "user-123", Points: 100}, nil)

	service := &UserService{
		Repo: mockRepo,
	}

	points, err := service.GetPoints("user-123")
	assert.Nil(t, err)
	assert.Equal(t, 100, points)
	mockRepo.AssertExpectations(t)
}
