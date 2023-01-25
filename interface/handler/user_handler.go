package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/interface/middleware"
	"GoBBS/usecase"
)

type userHandler struct {
	uc             usecase.User
	authMiddleware middleware.Auth
	jsonMarshal    func(any) ([]byte, error)
}

// NewUserHandler ユーザーハンドラーを生成する
func NewUserHandler(usecase usecase.User) *userHandler {
	return &userHandler{
		uc:             usecase,
		authMiddleware: middleware.NewAuth(usecase),
		jsonMarshal:    json.Marshal,
	}
}

// RegistUserHandler ハンドラー登録
func (h *userHandler) RegistHandlerFunc() {
	http.HandleFunc("/users", h.new)
	http.HandleFunc("/users/", h.edit)
	http.HandleFunc("/login", h.auth)
}

// new 新規作成
func (h *userHandler) new(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromReqBody(r.Body)
	if err != nil {
		log.Printf("get user error : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.regist(w, user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// edit 編集
func (h *userHandler) edit(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromReqBody(r.Body)
	if err != nil {
		log.Printf("get user error : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := h.getUserIDFromPathParam(r.URL.Path)
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.ID = userID

	switch r.Method {
	case http.MethodPut:
		if !h.authMiddleware.VerifyAuth(r.Header, user.ID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.update(w, user)
	case http.MethodDelete:
		if !h.authMiddleware.VerifyAuth(r.Header, user.ID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.delete(w, user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// login ログインハンドラー
func (h *userHandler) auth(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromReqBody(r.Body)
	if err != nil {
		log.Printf("get user error : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.login(w, user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// regist ユーザー登録
func (h *userHandler) regist(w http.ResponseWriter, user dto.User) {
	if err := h.uc.Regist(&user, time.Now()); err != nil {
		log.Printf("regist error : %v", err)
		if errors.Is(err, service.ErrUserAlreadyRegistered) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// update ユーザー更新
func (h *userHandler) update(w http.ResponseWriter, user dto.User) {
	if err := h.uc.Update(&user, time.Now()); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("update error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// delete ユーザー削除
func (h *userHandler) delete(w http.ResponseWriter, user dto.User) {
	if err := h.uc.Delete(&user); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("delete error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// login ログイン
func (h *userHandler) login(w http.ResponseWriter, user dto.User) {
	token, err := h.uc.Authorize(user.Email, user.Password)
	if err != nil {
		log.Printf("login authorize error: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonToken, err := h.jsonMarshal(dto.NewToken(token))
	if err != nil {
		log.Printf("login token marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonToken)
}

// getUserFromReqBody リクエストボディからユーザー情報を取得する
func (h *userHandler) getUserFromReqBody(body io.ReadCloser) (dto.User, error) {
	var user dto.User
	err := json.NewDecoder(body).Decode(&user)

	return user, err
}

// getUserIDFromPathParam パスからユーザーIDを取得する
func (h *userHandler) getUserIDFromPathParam(path string) string {
	subPath := strings.Split(path, "/")
	if len(subPath) != 3 {
		return ""
	}

	return subPath[2]
}
