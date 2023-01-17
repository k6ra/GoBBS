package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExecWithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	got, err := ExecWithTx(db, func(tx *sql.Tx) (string, error) {
		return "ok", nil
	})

	if err != nil {
		t.Errorf("予期せぬエラー(error: %s)", err)
	}

	want := "ok"
	if got != want {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestExecWithTx_BeginFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("ng"))

	got, err := ExecWithTx(db, func(tx *sql.Tx) (string, error) {
		return "ok", nil
	})

	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	want := ""
	if got != want {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestExecWithTx_FuncFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	got, err := ExecWithTx(db, func(tx *sql.Tx) (string, error) {
		return "", errors.New("ng")
	})

	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	want := ""
	if got != want {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestExecWithTx_FuncFailRollbackFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("ng"))

	got, err := ExecWithTx(db, func(tx *sql.Tx) (string, error) {
		return "", errors.New("ng")
	})

	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	want := ""
	if got != want {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}

func TestExecWithTx_CommitFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmockの生成に失敗(error: %s)", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("ng"))

	got, err := ExecWithTx(db, func(tx *sql.Tx) (string, error) {
		return "ok", nil
	})
	fmt.Printf("%#v\n", mock)

	if err == nil {
		t.Errorf("予期せぬ正常終了")
	}

	want := ""
	if got != want {
		t.Errorf("戻り値不一致 got: %#v want: %#v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("予期せぬDB操作(error: %s)", err)
	}
}
