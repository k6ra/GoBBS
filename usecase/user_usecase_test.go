package usecase

import (
	"GoBBS/domain/model"
	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/mock/mock_service"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestNewUserUseCase(t *testing.T) {
	type args struct {
		db *sql.DB
		f  service.UserFactory
	}
	tests := []struct {
		name string
		args args
		want *userUseCase
	}{
		{
			name: "正常ケース",
			args: args{
				db: &sql.DB{},
				f:  &mock_service.MockUserFactory{},
			},
			want: &userUseCase{
				db:                 &sql.DB{},
				userServiceFactory: &mock_service.MockUserFactory{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserUseCase(tt.args.db, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userUseCase_Regist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user *dto.User
		now  time.Time
	}
	tests := []struct {
		name    string
		uc      *userUseCase
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成失敗(error: %v)", err)
					}
					mock.ExpectBegin()
					mock.ExpectCommit()
					return db
				}(),
				userServiceFactory: func() *mock_service.MockUserFactory {
					svc := mock_service.NewMockUser(ctrl)
					svc.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(nil)

					mock := mock_service.NewMockUserFactory(ctrl)
					mock.EXPECT().NewUserService(gomock.Any()).Return(svc)
					return mock
				}(),
			},
			args: args{
				user: &dto.User{},
				now:  time.Now(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成に失敗(error: %v)", err)
					}
					mock.ExpectBegin().WillReturnError(errors.New("ng"))
					return db
				}(),
			},
			args: args{
				user: &dto.User{},
				now:  time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.uc.Regist(tt.args.user, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("userUseCase.Regist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userUseCase_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		uc      *userUseCase
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "正常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成失敗(error: %v)", err)
					}
					mock.ExpectBegin()
					mock.ExpectCommit()
					return db
				}(),
				userServiceFactory: func() *mock_service.MockUserFactory {
					svc := mock_service.NewMockUser(ctrl)
					svc.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(
						model.NewUser("id", "name", "email", "password"),
						nil,
					)

					mock := mock_service.NewMockUserFactory(ctrl)
					mock.EXPECT().NewUserService(gomock.Any()).Return(svc)
					return mock
				}(),
			},
			args: args{
				email:    "email",
				password: "password",
			},
			want:    model.NewUser("id", "name", "email", "password"),
			wantErr: false,
		},
		{
			name: "異常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成に失敗(error: %v)", err)
					}
					mock.ExpectBegin().WillReturnError(errors.New("ng"))
					return db
				}(),
			},
			args: args{
				email:    "email",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.uc.Authorize(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUseCase.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userUseCase.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user *dto.User
		now  time.Time
	}
	tests := []struct {
		name    string
		uc      *userUseCase
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成失敗(error: %v)", err)
					}
					mock.ExpectBegin()
					mock.ExpectCommit()
					return db
				}(),
				userServiceFactory: func() *mock_service.MockUserFactory {
					svc := mock_service.NewMockUser(ctrl)
					svc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

					mock := mock_service.NewMockUserFactory(ctrl)
					mock.EXPECT().NewUserService(gomock.Any()).Return(svc)
					return mock
				}(),
			},
			args: args{
				user: &dto.User{},
				now:  time.Now(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成に失敗(error: %v)", err)
					}
					mock.ExpectBegin().WillReturnError(errors.New("ng"))
					return db
				}(),
			},
			args: args{
				user: &dto.User{},
				now:  time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.uc.Update(tt.args.user, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("userUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user *dto.User
	}
	tests := []struct {
		name    string
		uc      *userUseCase
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成失敗(error: %v)", err)
					}
					mock.ExpectBegin()
					mock.ExpectCommit()
					return db
				}(),
				userServiceFactory: func() *mock_service.MockUserFactory {
					svc := mock_service.NewMockUser(ctrl)
					svc.EXPECT().Delete(gomock.Any()).Return(nil)

					mock := mock_service.NewMockUserFactory(ctrl)
					mock.EXPECT().NewUserService(gomock.Any()).Return(svc)
					return mock
				}(),
			},
			args: args{
				user: &dto.User{},
			},
			wantErr: false,
		},
		{
			name: "異常ケース",
			uc: &userUseCase{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("sqlmockの生成に失敗(error: %v)", err)
					}
					mock.ExpectBegin().WillReturnError(errors.New("ng"))
					return db
				}(),
			},
			args: args{
				user: &dto.User{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.uc.Delete(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
