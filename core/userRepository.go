package core

type UserRepository interface {
	GetByID(id string) (*User, error)
	Save(user *User) error
	SubtractPoints(user *User, points int64) error
	AddPoints(user *User, points int64) error
}
