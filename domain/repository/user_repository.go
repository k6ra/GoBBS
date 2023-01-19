package repository

import (
	"GoBBS/domain/model"
	"time"
)

// User ユーザーリポジトリ
// mockgen -source domain/repository/user_repository.go -destination mock/mock_repository/user_repository_mock.go
type User interface {
	FindByEmail(email string) (model.User, error)
	Regist(user model.User, now time.Time) error
	Update(user model.User, now time.Time) error
	Delete(user model.User) error
}
