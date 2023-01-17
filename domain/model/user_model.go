package model

// User ユーザー
type User struct {
	id       string
	name     string
	email    string
	password string
}

// NewUser ユーザーを生成する
func NewUser(id string, name string, email string, password string) *User {
	return &User{
		id:       id,
		name:     name,
		email:    email,
		password: password,
	}
}

// ID IDを返す
func (u *User) ID() string {
	return u.id
}

// Name 名前を返す
func (u *User) Name() string {
	return u.name
}

// Email メールアドレスを返す
func (u *User) Email() string {
	return u.email
}

// Password パスワードを返す
func (u *User) Password() string {
	return u.password
}

// VerifyPassword パスワードを検証する
func (u *User) VerifyPassword(password string) bool {
	return u.password == password
}
