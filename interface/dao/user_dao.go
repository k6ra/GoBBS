package dao

import (
	"database/sql"
	"time"

	"GoBBS/domain/model"
	"GoBBS/domain/repository"
	"GoBBS/dto"
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
func (u *UserDAO) FindByEmail(email string) (*model.User, error) {
	rows, err := u.tx.Query("select id, name, email, password from user where email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user dto.User
	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		return user.MapUserModel(), nil
	}
	return nil, nil
}

// Regist ユーザーを登録する
func (u *UserDAO) Regist(user *model.User, now time.Time) error {
	stmt, err := u.tx.Prepare(`
		insert into user (email, name, password, created_at, updated_at)
		values(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(
		user.Email(),
		user.Name(),
		user.Password(),
		now,
		now,
	); err != nil {
		return err
	}

	return nil
}

// Update ユーザーを更新する
func (u *UserDAO) Update(user *model.User, now time.Time) error {
	rows, err := u.tx.Query("select id, name, email, password from user where id = ? for update", user.ID())
	if err != nil {
		return err
	}
	defer rows.Close()

	stmt, err := u.tx.Prepare("update user set name = ?, password = ?, updated_at = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		user.Name(),
		user.Password(),
		now,
		user.ID(),
	); err != nil {
		return err
	}

	return nil
}

// Delete ユーザーを削除する
func (u *UserDAO) Delete(user *model.User) error {
	stmt, err := u.tx.Prepare("delete from user where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		user.ID(),
	); err != nil {
		return err
	}

	return nil
}
