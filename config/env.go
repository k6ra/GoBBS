package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// env 環境変数
type env struct {
	dbHost           string
	dbName           string
	dbUser           string
	dbPassword       string
	jwtSecretKey     string
	corsAllowOrigin  string
	corsAllowMethods []string
	corsAllowHeaders []string
	corsMaxAge       int
}

// 環境変数キャッシュ
var envCache *env

// GetEnv 環境変数を取得する
func GetEnv() (*env, error) {
	if envCache != nil {
		return envCache, nil
	}

	envCache = &env{}
	envCache.dbHost = os.Getenv("DB_HOST")
	envCache.dbName = os.Getenv("DB_NAME")
	envCache.dbUser = os.Getenv("DB_USER")
	envCache.dbPassword = os.Getenv("DB_PASSWORD")
	envCache.jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	envCache.corsAllowOrigin = os.Getenv("CORS_ALLOW_ORIGIN")
	envCache.corsAllowMethods = strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ",")
	envCache.corsAllowHeaders = strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ",")

	if corsMaxAge, err := strconv.Atoi(os.Getenv("CORS_MAX_AGE")); err != nil {
		return nil, errors.Wrap(err, "GetEnv CORS_MAX_AGE error")
	} else {
		envCache.corsMaxAge = corsMaxAge
	}

	return envCache, nil
}

// DBHost DBホスト名を返す
func (e *env) DBHost() string {
	return e.dbHost
}

// DBName DB名を返す
func (e *env) DBName() string {
	return e.dbName
}

// DBUser DBユーザーを返す
func (e *env) DBUser() string {
	return e.dbUser
}

// DBPassword DBユーザーのパスワードを返す
func (e *env) DBPassword() string {
	return e.dbPassword
}

// JWTSecretKey JWTシークレットキーを返す
func (e *env) JWTSecretKey() string {
	return e.jwtSecretKey
}

// CORSAllowOrigin CORS許可オリジンを返す
func (e *env) CORSAllowOrigin() string {
	return e.corsAllowOrigin
}

// CORSAllowMethods CORS許可メソッドを返す
func (e *env) CORSAllowMethods() []string {
	return e.corsAllowMethods
}

// CORSAllowHeaders CORS許可ヘッダを返す
func (e *env) CORSAllowHeaders() []string {
	return e.corsAllowHeaders
}

// CORSMaxAge CORSのキャッシュ時間(秒)を返す
func (e *env) CORSMaxAge() int {
	return e.corsMaxAge
}
