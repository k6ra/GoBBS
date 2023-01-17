package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/usecase"
)

type userHandler struct {
	uc usecase.User
}

// NewUserHandler ユーザーハンドラーを生成する
func NewUserHandler(usecase usecase.User) *userHandler {
	return &userHandler{uc: usecase}
}

// RegistUserHandler ハンドラー登録
func (h *userHandler) RegistUserHandlerFunc() {
	http.HandleFunc("/user", h.userHandler)
}

// usewrHandler ユーザーハンドラー
func (h *userHandler) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.regist(w, r)
	case http.MethodPut:
		h.update(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	}
}

// regist ユーザー登録
func (h *userHandler) regist(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("regist decode error : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
func (h *userHandler) update(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("update decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

func (h *userHandler) delete(w http.ResponseWriter, r *http.Request) {
	var user dto.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("delete decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
