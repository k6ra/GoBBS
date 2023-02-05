package dao

import (
	"GoBBS/domain/model"
	"GoBBS/domain/repository"
	"GoBBS/mock/mock_model"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestNewUserDAO(t *testing.T) {
	type args struct {
		tx *sql.Tx
	}
	tests := []struct {
		name string
		args args
		want *UserDAO
	}{
		{
			name: "正常ケース",
			args: args{
				tx: &sql.Tx{},
			},
			want: &UserDAO{
				tx: &sql.Tx{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserDAO(tt.args.tx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserDAO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDAO_FindByEmailSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectQuery("select id, name, email, password, salt from user where email = ?").
		WithArgs("email").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "password", "salt"}).
				AddRow("1", "example 1", "email@email.com", "examplepas", "salt")).
		RowsWillBeClosed()

	dao := NewUserDAO(tx)
	got, err := dao.FindByEmail("email")
	if err != nil {
		t.Errorf("予期せぬエラー(error: %s)", err)
	}

	want := model.NewUser("1", "example 1", "email@email.com", "examplepas", "salt")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_FindByEmailNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectQuery("select id, name, email, password, salt from user where email = ?").
		WithArgs("email").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "password"})).
		RowsWillBeClosed()

	dao := NewUserDAO(tx)
	got, err := dao.FindByEmail("email")
	if err != repository.ErrUserNotFound {
		t.Errorf("予期せぬエラー(error: %s)", err)
	}

	if got != nil {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, nil)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_FindByEmailScanFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s", err)
	}

	mock.ExpectQuery("select id, name, email, password, salt from user where email = ?").
		WithArgs("email").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(nil, "example 1", "email@email.com", "examplepas")).
		RowsWillBeClosed()

	dao := NewUserDAO(tx)
	got, err := dao.FindByEmail("email")
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if got != nil {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, nil)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}

}

func TestUserDAO_FindByEmailQueryFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s", err)
	}

	mock.ExpectQuery("select id, name, email, password, salt from user where email = ?").
		WithArgs("email").
		WillReturnError(errors.New("ng"))

	dao := NewUserDAO(tx)
	got, err := dao.FindByEmail("email")
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if got != nil {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, nil)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_RegistSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	now := time.Now()
	mock.ExpectPrepare(`
		insert into user (email, name, password, salt, created_at, updated_at)
		values(?, ?, ?, ?, ?, ?)
	`).
		WillBeClosed()

	mock.ExpectExec("insert into user (email, name, password, salt, created_at, updated_at) values(?, ?, ?, ?, ?, ?)").
		WithArgs("email@email.com", "example 1", "examplepas", "examplesalt", now, now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	gomock.InOrder(
		mockUser.EXPECT().EncryptPassword().Return("examplepas", "examplesalt", nil),
		mockUser.EXPECT().Email().Return("email@email.com"),
		mockUser.EXPECT().Name().Return("example 1"),
	)
	err = dao.Regist(mockUser, now)
	if err != nil {
		t.Errorf("予期せぬエラー(error: %s)", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_RegistInsertFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	now := time.Now()
	mock.ExpectPrepare(`
		insert into user (email, name, password, salt, created_at, updated_at)
		values(?, ?, ?, ?, ?, ?)
	`).
		WillBeClosed()

	mock.ExpectExec("insert into user (email, name, password, salt, created_at, updated_at) values(?, ?, ?, ?, ?, ?)").
		WithArgs("email@email.com", "example 1", "examplepas", "examplesalt", now, now).
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	gomock.InOrder(
		mockUser.EXPECT().EncryptPassword().Return("examplepas", "examplesalt", nil),
		mockUser.EXPECT().Email().Return("email@email.com"),
		mockUser.EXPECT().Name().Return("example 1"),
	)

	err = dao.Regist(mockUser, now)
	if err == nil {
		t.Error("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_RegistPrepareFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	now := time.Now()
	mock.ExpectPrepare(`
		insert into user (email, name, password, salt, created_at, updated_at)
		values(?, ?, ?, ?, ?, ?)
	`).
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)

	err = dao.Regist(mockUser, now)
	if err == nil {
		t.Error("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_UpdateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectQuery("select id, name, email, password from user where id = ? for update").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow("1", "example 1", "email@email.com", "examplepas")).
		RowsWillBeClosed()

	mock.ExpectPrepare(`update user set name = ?, password = ?, updated_at = ? where id = ?`).
		WillBeClosed()

	now := time.Now()
	mock.ExpectExec("update user set name = ?, password = ?, updated_at = ? where id = ?").
		WithArgs("example 1", "examplepas", now, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	gomock.InOrder(
		mockUser.EXPECT().ID().Return("1"),
		mockUser.EXPECT().Name().Return("example 1"),
		mockUser.EXPECT().Password().Return("examplepas"),
		mockUser.EXPECT().ID().Return("1"),
	)

	err = dao.Update(mockUser, now)
	if err != nil {
		t.Errorf("予期せぬエラー(error: %s)", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_UpdateFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectQuery("select id, name, email, password from user where id = ? for update").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow("1", "example 1", "email@email.com", "examplepas")).
		RowsWillBeClosed()

	mock.ExpectPrepare(`update user set name = ?, password = ?, updated_at = ? where id = ?`).
		WillBeClosed()

	now := time.Now()
	mock.ExpectExec("update user set name = ?, password = ?, updated_at = ? where id = ?").
		WithArgs("example 1", "examplepas", now, "1").
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	gomock.InOrder(
		mockUser.EXPECT().ID().Return("1"),
		mockUser.EXPECT().Name().Return("example 1"),
		mockUser.EXPECT().Password().Return("examplepas"),
		mockUser.EXPECT().ID().Return("1"),
	)

	err = dao.Update(mockUser, now)
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_UpdatePrepareFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectQuery("select id, name, email, password from user where id = ? for update").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow("1", "example 1", "email@email.com", "examplepas")).
		RowsWillBeClosed()

	mock.ExpectPrepare(`update user set name = ?, password = ?, updated_at = ? where id = ?`).
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	mockUser.EXPECT().ID().Return("1")

	now := time.Now()
	err = dao.Update(mockUser, now)
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_UpdateQueryFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectQuery("select id, name, email, password from user where id = ? for update").
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	mockUser.EXPECT().ID().Return("1")

	now := time.Now()
	err = dao.Update(mockUser, now)
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}
func TestUserDAO_DeleteSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectPrepare(`delete from user where id = ?`).
		WillBeClosed()

	mock.ExpectExec("delete from user where id = ?").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	mockUser.EXPECT().ID().Return("1")

	err = dao.Delete(mockUser)
	if err != nil {
		t.Errorf("予期せぬエラー(error: %s)", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_DeleteFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectPrepare(`delete from user where id = ?`).
		WillBeClosed()

	mock.ExpectExec("delete from user where id = ?").
		WithArgs("1").
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)
	mockUser.EXPECT().ID().Return("1")

	err = dao.Delete(mockUser)
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestUserDAO_DeletePrepareFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("txの生成に失敗(error: %s)", err)
	}

	mock.ExpectPrepare(`delete from user where id = ?`).
		WillReturnError(errors.New("ng"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dao := NewUserDAO(tx)
	mockUser := mock_model.NewMockUser(ctrl)

	err = dao.Delete(mockUser)
	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}
