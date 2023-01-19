package dto

import (
	"GoBBS/domain/model"
)

// User ユーザー
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string `json:"-"`
}

// NewUser ユーザーモデルを元にDTOユーザーを生成する
func NewUser(user model.User) *User {
	return &User{
		ID:       user.ID(),
		Name:     user.Name(),
		Email:    user.Email(),
		Password: user.Password(),
		Salt:     user.Salt(),
	}
}

// MapUserModel DTOユーザーの情報を元にユーザーモデルを生成する
func (u *User) MapUserModel() model.User {
	return model.NewUser(u.ID, u.Name, u.Email, u.Password, u.Salt)
}
