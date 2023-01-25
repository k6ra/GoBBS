package config

import "os"

// Env 環境変数
type Env struct {
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	SecretKey  string
}

// 環境変数キャッシュ
var env *Env

// GetEnv 環境変数を取得する
func GetEnv() *Env {
	if env != nil {
		return env
	}

	env = &Env{}
	env.DBHost = os.Getenv("DB_HOST")
	env.DBName = os.Getenv("DB_NAME")
	env.DBUser = os.Getenv("DB_USER")
	env.DBPassword = os.Getenv("DB_PASSWORD")
	env.SecretKey = os.Getenv("SECRET_KEY")

	return env
}
