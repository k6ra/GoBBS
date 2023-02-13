package config

import (
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name    string
		init    func()
		want    *env
		wantErr bool
	}{
		{
			name: "正常ケース",
			init: func() {
				envCache = nil
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("MYSQL_ROOT_PASSWORD", "root")
				t.Setenv("MYSQL_DATABASE", "bbs")
				t.Setenv("DB_PORT", "3306")
				t.Setenv("DB_NAME", "bbs")
				t.Setenv("DB_USER", "user")
				t.Setenv("DB_PASSWORD", "password")
				t.Setenv("JWT_SECRET_KEY", "secretkey")
				t.Setenv("CORS_ALLOW_ORIGIN", "http://localhost")
				t.Setenv("CORS_ALLOW_METHODS", "POST")
				t.Setenv("CORS_ALLOW_HEADERS", "*")
				t.Setenv("CORS_MAX_AGE", "7200")
			},
			want: &env{
				dbHost:           "localhost",
				dbName:           "bbs",
				dbUser:           "user",
				dbPassword:       "password",
				jwtSecretKey:     "secretkey",
				corsAllowOrigin:  "http://localhost",
				corsAllowMethods: []string{"POST"},
				corsAllowHeaders: []string{"*"},
				corsMaxAge:       7200,
			},
			wantErr: false,
		},
		{
			name: "異常ケース(CORSキャッシュ時間数値以外)",
			init: func() {
				envCache = nil
				t.Setenv("DB_HOST", "localhost")
				t.Setenv("MYSQL_ROOT_PASSWORD", "root")
				t.Setenv("MYSQL_DATABASE", "bbs")
				t.Setenv("DB_PORT", "3306")
				t.Setenv("DB_NAME", "bbs")
				t.Setenv("DB_USER", "user")
				t.Setenv("DB_PASSWORD", "password")
				t.Setenv("JWT_SECRET_KEY", "secretkey")
				t.Setenv("CORS_ALLOW_ORIGIN", "http://localhost")
				t.Setenv("CORS_ALLOW_METHODS", "POST")
				t.Setenv("CORS_ALLOW_HEADERS", "*")
				t.Setenv("CORS_MAX_AGE", "ng")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "正常ケース(キャッシュ返却)",
			init: func() {
				envCache = &env{
					dbHost:           "localhost",
					dbName:           "bbs",
					dbUser:           "user",
					dbPassword:       "password",
					jwtSecretKey:     "secretkey",
					corsAllowOrigin:  "http://localhost",
					corsAllowMethods: []string{"POST"},
					corsAllowHeaders: []string{"*"},
					corsMaxAge:       7200,
				}
				t.Setenv("DB_HOST", "")
				t.Setenv("MYSQL_ROOT_PASSWORD", "")
				t.Setenv("MYSQL_DATABASE", "")
				t.Setenv("DB_PORT", "")
				t.Setenv("DB_NAME", "")
				t.Setenv("DB_USER", "")
				t.Setenv("DB_PASSWORD", "")
				t.Setenv("JWT_SECRET_KEY", "")
				t.Setenv("CORS_ALLOW_ORIGIN", "")
				t.Setenv("CORS_ALLOW_METHODS", "")
				t.Setenv("CORS_ALLOW_HEADERS", "")
				t.Setenv("CORS_MAX_AGE", "")
			},
			want: &env{
				dbHost:           "localhost",
				dbName:           "bbs",
				dbUser:           "user",
				dbPassword:       "password",
				jwtSecretKey:     "secretkey",
				corsAllowOrigin:  "http://localhost",
				corsAllowMethods: []string{"POST"},
				corsAllowHeaders: []string{"*"},
				corsMaxAge:       7200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init()
			got, err := GetEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_DBHost(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want string
	}{
		{
			name: "正常ケース",
			e: &env{
				dbHost: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.DBHost(); got != tt.want {
				t.Errorf("env.DBHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_DBName(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want string
	}{
		{
			name: "正常ケース",
			e: &env{
				dbName: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.DBName(); got != tt.want {
				t.Errorf("env.DBName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_DBUser(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want string
	}{
		{
			name: "正常ケース",
			e: &env{
				dbUser: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.DBUser(); got != tt.want {
				t.Errorf("env.DBUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_DBPassword(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want string
	}{
		{
			name: "正常ケース",
			e: &env{
				dbPassword: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.DBPassword(); got != tt.want {
				t.Errorf("env.DBPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_JWTSecretKey(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want string
	}{
		{
			name: "正常ケース",
			e: &env{
				jwtSecretKey: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.JWTSecretKey(); got != tt.want {
				t.Errorf("env.JWTSecretKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_CORSAllowOrigin(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want string
	}{
		{
			name: "正常ケース",
			e: &env{
				corsAllowOrigin: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.CORSAllowOrigin(); got != tt.want {
				t.Errorf("env.CORSAllowOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_CORSAllowMethods(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want []string
	}{
		{
			name: "正常ケース",
			e: &env{
				corsAllowMethods: []string{"1", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.CORSAllowMethods(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("env.CORSAllowMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_CORSAllowHeaders(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want []string
	}{
		{
			name: "正常ケース",
			e: &env{
				corsAllowHeaders: []string{"1", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.CORSAllowHeaders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("env.CORSAllowHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_CORSMaxAge(t *testing.T) {
	tests := []struct {
		name string
		e    *env
		want int
	}{
		{
			name: "正常ケース",
			e: &env{
				corsMaxAge: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.CORSMaxAge(); got != tt.want {
				t.Errorf("env.CORSMaxAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
