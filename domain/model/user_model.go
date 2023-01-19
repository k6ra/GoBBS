package model

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"

	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
)

type (
	// User ユーザーサービス
	// mockgen -source domain/model/user_model.go -destination mock/mock_model/user_model_mock.go
	User interface {
		ID() string
		Name() string
		Email() string
		Password() string
		Salt() string
		VerifyPassword(string) (bool, error)
		EncryptPassword() (string, string, error)
	}

	// user ユーザー
	user struct {
		id       string
		name     string
		email    string
		password string
		salt     string
	}
)

// NewUser ユーザーを生成する
func NewUser(id string, name string, email string, password string, salt string) User {
	return &user{
		id:       id,
		name:     name,
		email:    email,
		password: password,
		salt:     salt,
	}
}

// saltLen ソルトの長さ
const saltLen = 32

// ID IDを返す
func (u *user) ID() string {
	return u.id
}

// Name 名前を返す
func (u *user) Name() string {
	return u.name
}

// Email メールアドレスを返す
func (u *user) Email() string {
	return u.email
}

// Password パスワードを返す
func (u *user) Password() string {
	return u.password
}

// Salt ソルトを返す
func (u *user) Salt() string {
	return u.salt
}

// VerifyPassword パスワードを検証する
func (u *user) VerifyPassword(password string) (bool, error) {
	cryptPw, err := u.encryptPassword(password, u.salt)
	if err != nil {
		return false, errors.Wrap(err, "VerifyPassword error")
	}

	return u.password == cryptPw, nil
}

// EncryptPassword ソルトの生成およびパスワードを暗号化をし返す
func (u *user) EncryptPassword() (string, string, error) {
	salt, err := u.generateSalt()
	if err != nil {
		return "", "", errors.Wrap(err, "EncryptPassword error")
	}

	cryptPw, err := u.encryptPassword(u.password, salt)
	if err != nil {
		return "", "", errors.Wrap(err, "EncryptPassword error")
	}

	return cryptPw, salt, nil
}

// encryptPassword パスワードを暗号化する
func (u *user) encryptPassword(password string, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(password), []byte(salt), 1<<15, 8, 1, 32)
	if err != nil {
		return "", errors.Wrap(err, "encryptPassword error")
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}

// generateSalt ソルトを生成する
func (u *user) generateSalt() (string, error) {
	// ソルトの桁数分 asciiコード33(!)～126(~)の文字を生成する
	salt := make([]rune, saltLen)
	for i := 0; i < saltLen; i++ {
		val, err := rand.Int(rand.Reader, big.NewInt(94))
		if err != nil {
			return "", errors.Wrap(err, "generateSalt error")
		}
		salt[i] = rune(val.Int64() + 33)
	}

	return string(salt), nil
}
