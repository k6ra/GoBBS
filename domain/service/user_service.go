package service

import (
	"errors"
	"time"

	"GoBBS/domain/model"
	"GoBBS/domain/repository"
)

type (
	// User ユーザーサービス
	// mockgen -source domain/service/user_service.go -destination mock/mock_service/user_service_mock.go
	User interface {
		Authorize(email string, password string) (*model.User, error)
		IsDuplicate(email string) (bool, error)
		Regist(user *model.User, now time.Time) error
		Update(user *model.User, now time.Time) error
		Delete(user *model.User) error
	}

	// UserFactory ユーザーサービスファクトリー
	UserFactory interface {
		NewUserService(repo repository.User) User
	}

	userService struct {
		repo repository.User
	}

	userServiceFactory struct{}
)

var _ User = (*userService)(nil)

var (
	ErrAuthorizeFail         = errors.New("authorize failed")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyRegistered = errors.New("user already registered")
)

// NewUserServiceFactory ユーザーサービスファクトリーを生成する
func NewUserServiceFactory() *userServiceFactory {
	return &userServiceFactory{}
}

// NewUserService ユーザーサービスを生成する
func (f *userServiceFactory) NewUserService(repo repository.User) User {
	return &userService{repo: repo}
}

// Authorize ユーザーを認証する
func (s *userService) Authorize(email string, password string) (*model.User, error) {
	if user, err := s.repo.FindByEmail(email); err != nil {
		return nil, err
	} else if !user.VerifyPassword(password) {
		return nil, ErrAuthorizeFail
	} else {
		return user, nil
	}
}

// IsDuplicate 与えられたメールアドレスが登録済みか判定する
func (s *userService) IsDuplicate(email string) (bool, error) {
	if user, err := s.repo.FindByEmail(email); err != nil {
		return false, err
	} else {
		return user != nil, nil
	}
}

// Regist ユーザーを登録する
func (s *userService) Regist(user *model.User, now time.Time) error {
	if duplicate, err := s.IsDuplicate(user.Email()); err != nil {
		return err
	} else if duplicate {
		return ErrUserAlreadyRegistered
	}

	return s.repo.Regist(user, now)
}

// Update ユーザーを更新する
func (s *userService) Update(user *model.User, now time.Time) error {
	findUser, err := s.repo.FindByEmail(user.Email())
	if err != nil {
		return err
	}
	if findUser == nil {
		return ErrUserNotFound
	}

	margedUser := model.NewUser(
		findUser.ID(),
		user.Name(),
		findUser.Email(),
		user.Password(),
	)

	return s.repo.Update(margedUser, now)
}

// Delete ユーザーを削除する
func (s *userService) Delete(user *model.User) error {
	findUser, err := s.repo.FindByEmail(user.Email())
	if err != nil {
		return err
	}
	if findUser == nil {
		return ErrUserNotFound
	}

	margedUser := model.NewUser(
		findUser.ID(),
		user.Name(),
		findUser.Email(),
		user.Password(),
	)

	return s.repo.Delete(margedUser)
}
