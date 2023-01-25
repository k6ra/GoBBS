package dao

import (
	"database/sql"
	"time"

	"GoBBS/domain/model"
	"GoBBS/domain/repository"
	"GoBBS/dto"

	"github.com/pkg/errors"
)

// UserDAO ユーザーDAO
type UserDAO struct {
	tx *sql.Tx
}

var _ repository.User = (*UserDAO)(nil)

// NewUserDAO ユーザーDAOを生成する
func NewUserDAO(tx *sql.Tx) *UserDAO {
	return &UserDAO{
		tx: tx,
	}
}

// FindByEmail メールアドレスを指定してユーザーを取得する
func (u *UserDAO) FindByEmail(email string) (model.User, error) {
	rows, err := u.tx.Query("select id, name, email, password, salt from user where email = ?", email)
	if err != nil {
		return nil, errors.Wrap(err, "FindByEmail error")
	}
	defer rows.Close()

	var user dto.User
	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Salt); err != nil {
			return nil, errors.Wrap(err, "FindByEmail error")
		}
		return user.MapUserModel(), nil
	}
	return nil, nil
}

// Regist ユーザーを登録する
func (u *UserDAO) Regist(user model.User, now time.Time) error {
	stmt, err := u.tx.Prepare(`
		insert into user (email, name, password, salt, created_at, updated_at)
		values(?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return errors.Wrap(err, "Regist error")
	}
	defer stmt.Close()

	cryptPw, salt, err := user.EncryptPassword()
	if err != nil {
		return errors.Wrap(err, "Regist error")
	}
	if _, err = stmt.Exec(
		user.Email(),
		user.Name(),
		cryptPw,
		salt,
		now,
		now,
	); err != nil {
		return errors.Wrap(err, "Regist error")
	}

	return nil
}

// Update ユーザーを更新する
func (u *UserDAO) Update(user model.User, now time.Time) error {
	rows, err := u.tx.Query("select id, name, email, password from user where id = ? for update", user.ID())
	if err != nil {
		return errors.Wrap(err, "Update error")
	}
	rows.Close()

	stmt, err := u.tx.Prepare("update user set name = ?, password = ?, updated_at = ? where id = ?")
	if err != nil {
		return errors.Wrap(err, "Update error")
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		user.Name(),
		user.Password(),
		now,
		user.ID(),
	); err != nil {
		return errors.Wrap(err, "Update error")
	}

	return nil
}

// Delete ユーザーを削除する
func (u *UserDAO) Delete(user model.User) error {
	stmt, err := u.tx.Prepare("delete from user where id = ?")
	if err != nil {
		return errors.Wrap(err, "Delete error")
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		user.ID(),
	); err != nil {
		return errors.Wrap(err, "Delete error")
	}

	return nil
}
