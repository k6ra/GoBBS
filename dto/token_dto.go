package dto

// Token トークン
type Token struct {
	Token string `json:"token"`
}

// NewToken トークンを生成する
func NewToken(token string) *Token {
	return &Token{
		Token: token,
	}
}
