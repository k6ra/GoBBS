package usecase

import (
	"database/sql"
	"time"

	"GoBBS/domain/model"
	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/interface/dao"
)

// User ユーザーユースケース
// mockgen -source usecase/user_usecase.go -destination mock/mock_usecase/user_usecase_mock.go
type User interface {
	Regist(*dto.User, time.Time) error
	Update(*dto.User, time.Time) error
	Authorize(email string, password string) (*model.User, error)
	Delete(*dto.User) error
}

type userUseCase struct {
	db                 *sql.DB
	userServiceFactory service.UserFactory
}

var _ User = (*userUseCase)(nil)

// NewUserUseCase ユーザーユースケースを生成する
func NewUserUseCase(db *sql.DB, f service.UserFactory) *userUseCase {
	return &userUseCase{
		db:                 db,
		userServiceFactory: f,
	}
}

// Regist ユーザーを登録する
func (uc *userUseCase) Regist(user *dto.User, now time.Time) error {
	_, err := dao.ExecWithTx(
		uc.db,
		func(tx *sql.Tx) (any, error) {
			return nil, uc.userServiceFactory.NewUserService(dao.NewUserDAO(tx)).Regist(user.MapUserModel(), now)
		},
	)

	return err
}

// Authorize 認証する
func (uc *userUseCase) Authorize(email string, password string) (*model.User, error) {
	user, err := dao.ExecWithTx(
		uc.db,
		func(tx *sql.Tx) (*model.User, error) {
			return uc.userServiceFactory.NewUserService(dao.NewUserDAO(tx)).Authorize(email, password)
		},
	)

	return user, err
}

// Update 更新する
func (uc *userUseCase) Update(user *dto.User, now time.Time) error {
	_, err := dao.ExecWithTx(
		uc.db,
		func(tx *sql.Tx) (any, error) {
			return nil, uc.userServiceFactory.NewUserService(dao.NewUserDAO(tx)).Update(user.MapUserModel(), now)
		},
	)

	return err
}

// Delete 削除する
func (uc *userUseCase) Delete(user *dto.User) error {
	_, err := dao.ExecWithTx(
		uc.db,
		func(tx *sql.Tx) (any, error) {
			return nil, uc.userServiceFactory.NewUserService(dao.NewUserDAO(tx)).Delete(user.MapUserModel())
		},
	)

	return err
}
