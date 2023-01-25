package usecase

import (
	"database/sql"
	"time"

	"GoBBS/domain/model"
	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/interface/dao"
	"GoBBS/interface/security"
)

// User ユーザーユースケース
// mockgen -source usecase/user_usecase.go -destination mock/mock_usecase/user_usecase_mock.go
type User interface {
	Regist(*dto.User, time.Time) error
	Update(*dto.User, time.Time) error
	Authorize(email string, password string) (string, error)
	VerifyAuthorization(token string, userID string) bool
	Delete(*dto.User) error
}

type userUseCase struct {
	db                 *sql.DB
	userServiceFactory service.UserFactory
	token              security.Token
}

var _ User = (*userUseCase)(nil)

// NewUserUseCase ユーザーユースケースを生成する
func NewUserUseCase(db *sql.DB, f service.UserFactory, t security.Token) *userUseCase {
	return &userUseCase{
		db:                 db,
		userServiceFactory: f,
		token:              t,
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
func (uc *userUseCase) Authorize(email string, password string) (string, error) {
	user, err := dao.ExecWithTx(
		uc.db,
		func(tx *sql.Tx) (model.User, error) {
			return uc.userServiceFactory.NewUserService(dao.NewUserDAO(tx)).Authorize(email, password)
		},
	)

	if err != nil {
		return "", err
	}

	token, err := uc.token.Generate(user.ID(), time.Now())
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyAuthorization トークンを検証する
func (uc *userUseCase) VerifyAuthorization(token string, userID string) bool {
	return uc.token.Verify(token, userID)
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
