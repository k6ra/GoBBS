package service

import (
	"GoBBS/domain/model"
	"GoBBS/domain/repository"
	"GoBBS/mock/mock_model"
	"GoBBS/mock/mock_repository"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		repo repository.User
	}
	tests := []struct {
		name string
		f    *userServiceFactory
		args args
		want *userService
	}{
		{
			name: "正常ケース",
			args: args{
				repo: mock_repository.NewMockUser(ctrl),
			},
			want: &userService{
				repo: mock_repository.NewMockUser(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.NewUserService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		want    model.User
		wantErr bool
	}{
		{
			name: "正常ケース",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().VerifyPassword(gomock.Any()).Return(true, nil)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(mockUser, nil)
					return mock
				}(),
			},
			args: args{
				email:    "email",
				password: "password",
			},
			want: func() model.User {
				mockUser := mock_model.NewMockUser(ctrl)
				return mockUser
			}(),
			wantErr: false,
		},
		{
			name: "異常ケース(パスワード不一致)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().VerifyPassword(gomock.Any()).Return(false, nil)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(mockUser, nil)
					return mock
				}(),
			},
			args: args{
				email:    "email",
				password: "ng",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常ケース(パスワード検証エラー)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().VerifyPassword(gomock.Any()).Return(false, errors.New("ng"))
					mock.EXPECT().FindByEmail(gomock.Any()).Return(mockUser, nil)
					return mock
				}(),
			},
			args: args{
				email:    "email",
				password: "ng",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常ケース(ユーザー取得失敗)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(
						nil,
						errors.New("test"),
					)
					return mock
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
			got, err := tt.s.Authorize(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_IsDuplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "正常ケース(メールアドレス登録済み)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mockUser := mock_model.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(mockUser, nil)
					return mock
				}(),
			},
			args: args{
				email: "email",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "正常ケース(メールアドレス未登録)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, repository.ErrUserNotFound)
					return mock
				}(),
			},
			args: args{
				email: "email",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "異常ケース(ユーザー取得失敗)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				email: "email",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.IsDuplicate(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.IsDuplicate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userService.IsDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Regist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user model.User
		now  time.Time
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					gomock.InOrder(
						mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, repository.ErrUserNotFound),
						mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(nil),
					)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(登録チェックエラー)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "異常ケース(登録済みエラー)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(model.NewUser("id", "name", "email", "password", "salt"), nil)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "異常ケース(登録エラー)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					gomock.InOrder(
						mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, repository.ErrUserNotFound),
						mock.EXPECT().Regist(gomock.Any(), gomock.Any()).Return(errors.New("test")),
					)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Regist(tt.args.user, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("userService.Regist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user model.User
		now  time.Time
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mockUser := mock_model.NewMockUser(ctrl)
					gomock.InOrder(
						mockUser.EXPECT().ID().Return("findID"),
						mockUser.EXPECT().Email().Return("findEmail"),
						mockUser.EXPECT().Salt().Return("findSalt"),
					)
					gomock.InOrder(
						mock.EXPECT().FindByEmail(gomock.Any()).Return(mockUser, nil),
						mock.EXPECT().Update(
							model.NewUser("findID", "name", "findEmail", "password", "findSalt"),
							gomock.Any(),
						).Return(nil),
					)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					gomock.InOrder(
						mockUser.EXPECT().Email().Return("email"),
						mockUser.EXPECT().Name().Return("name"),
						mockUser.EXPECT().Password().Return("password"),
					)
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(ユーザーが未登録)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, nil)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "異常ケース(ユーザー取得失敗)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(
						nil,
						errors.New("test"),
					)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
				now: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.user, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("userService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		wantErr bool
	}{
		{
			name: "正常ケース",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mockUser := mock_model.NewMockUser(ctrl)
					gomock.InOrder(
						mockUser.EXPECT().ID().Return("findID"),
						mockUser.EXPECT().Email().Return("findEmail"),
						mockUser.EXPECT().Salt().Return("findSalt"),
					)
					gomock.InOrder(
						mock.EXPECT().FindByEmail(gomock.Any()).Return(mockUser, nil),
						mock.EXPECT().Delete(model.NewUser("findID", "name", "findEmail", "password", "findSalt")).Return(nil),
					)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					gomock.InOrder(
						mockUser.EXPECT().Email().Return("email"),
						mockUser.EXPECT().Name().Return("name"),
						mockUser.EXPECT().Password().Return("password"),
					)
					return mockUser
				}(),
			},
			wantErr: false,
		},
		{
			name: "異常ケース(ユーザー未登録)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, nil)
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
			},
			wantErr: true,
		},
		{
			name: "異常ケース(ユーザー取得失敗)",
			s: &userService{
				repo: func() *mock_repository.MockUser {
					mock := mock_repository.NewMockUser(ctrl)
					mock.EXPECT().FindByEmail(gomock.Any()).Return(nil, errors.New("ng"))
					return mock
				}(),
			},
			args: args{
				user: func() model.User {
					mockUser := mock_model.NewMockUser(ctrl)
					mockUser.EXPECT().Email().Return("email")
					return mockUser
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewUserServiceFactory(t *testing.T) {
	tests := []struct {
		name string
		want *userServiceFactory
	}{
		{
			name: "正常ケース",
			want: &userServiceFactory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserServiceFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserServiceFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}
