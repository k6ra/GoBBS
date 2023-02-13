package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"GoBBS/domain/service"
	"GoBBS/dto"
	"GoBBS/interface/handler/handlerctx"
	"GoBBS/interface/middleware"
	"GoBBS/interface/middleware/middlewarehelper"
	"GoBBS/usecase"
)

type userHandler struct {
	corsAllowOrigin  string
	corsAllowMethods []string
	corsAllowHeaders []string
	corsAllowMaxAge  int
	uc               usecase.User
	authMiddleware   middleware.Auth
}

// NewUserHandler ユーザーハンドラーを生成する
func NewUserHandler(
	usecase usecase.User,
	corsAllowOrigin string,
	corsAllowMethods []string,
	corsAllowHeaders []string,
	corsAllowMaxAge int) *userHandler {
	return &userHandler{
		corsAllowOrigin:  corsAllowOrigin,
		corsAllowMethods: corsAllowMethods,
		corsAllowHeaders: corsAllowHeaders,
		corsAllowMaxAge:  corsAllowMaxAge,
		uc:               usecase,
		authMiddleware:   middleware.NewAuth(usecase),
	}
}

// RegistUserHandler ハンドラー登録
func (h *userHandler) RegistHandlerFunc() {
	http.HandleFunc(
		"/users",
		middlewarehelper.Apply(
			handlerctx.NewAPIContext,
			h.new,
			middleware.NewCORS(
				h.corsAllowOrigin,
				h.corsAllowMethods,
				h.corsAllowHeaders,
				h.corsAllowMaxAge,
			).AddResponseHeader,
		),
	)

	http.HandleFunc(
		"/users/",
		middlewarehelper.Apply(
			handlerctx.NewAPIContext,
			h.edit,
			middleware.NewAuth(h.uc).VerifyAuth,
			middleware.NewPathParam("/users/:id").Parse,
		),
	)

	http.HandleFunc(
		"/login",
		middlewarehelper.Apply(
			handlerctx.NewAPIContext,
			h.auth,
		),
	)
}

// new 新規作成
func (h *userHandler) new(c handlerctx.APIContext) error {
	user, err := h.getUserFromReqBody(c.RequestBody())
	if err != nil {
		log.Printf("get user error : %v", err)
		c.WriteStatusCode(http.StatusBadRequest)
		return nil
	}

	switch c.RequestMethod() {
	case http.MethodPost:
		h.regist(c, user)
	default:
		c.WriteStatusCode(http.StatusMethodNotAllowed)
	}

	return nil
}

// edit 編集
func (h *userHandler) edit(c handlerctx.APIContext) error {
	user, err := h.getUserFromReqBody(c.RequestBody())
	if err != nil {
		log.Printf("get user error : %v", err)
		c.WriteStatusCode(http.StatusBadRequest)
		return nil
	}

	userID := c.PathParam()
	if userID == "" {
		c.WriteStatusCode(http.StatusBadRequest)
		return nil
	}
	user.ID = userID

	switch c.RequestMethod() {
	case http.MethodPut:
		h.update(c, user)
	case http.MethodDelete:
		h.delete(c, user)
	default:
		c.WriteStatusCode(http.StatusMethodNotAllowed)
	}

	return nil
}

// login ログインハンドラー
func (h *userHandler) auth(c handlerctx.APIContext) error {
	user, err := h.getUserFromReqBody(c.RequestBody())
	if err != nil {
		log.Printf("get user error : %v", err)
		c.WriteStatusCode(http.StatusBadRequest)
		return nil
	}

	switch c.RequestMethod() {
	case http.MethodPost:
		return h.login(c, user)
	default:
		c.WriteStatusCode(http.StatusMethodNotAllowed)
	}

	return nil
}

// regist ユーザー登録
func (h *userHandler) regist(c handlerctx.APIContext, user dto.User) {
	if err := h.uc.Regist(&user, time.Now()); err != nil {
		log.Printf("regist error : %v", err)
		if errors.Is(err, service.ErrUserAlreadyRegistered) {
			c.WriteStatusCode(http.StatusBadRequest)
			return
		}
		c.WriteStatusCode(http.StatusInternalServerError)
		return
	}

	c.WriteStatusCode(http.StatusOK)
}

// update ユーザー更新
func (h *userHandler) update(c handlerctx.APIContext, user dto.User) {
	if err := h.uc.Update(&user, time.Now()); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.WriteStatusCode(http.StatusBadRequest)
			return
		}
		log.Printf("update error: %v", err)
		c.WriteStatusCode(http.StatusInternalServerError)
		return
	}

	c.WriteStatusCode(http.StatusOK)
}

// delete ユーザー削除
func (h *userHandler) delete(c handlerctx.APIContext, user dto.User) {
	if err := h.uc.Delete(&user); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.WriteStatusCode(http.StatusBadRequest)
			return
		}

		log.Printf("delete error: %v", err)
		c.WriteStatusCode(http.StatusInternalServerError)
		return
	}

	c.WriteStatusCode(http.StatusOK)
}

// login ログイン
func (h *userHandler) login(c handlerctx.APIContext, user dto.User) error {
	token, err := h.uc.Authorize(user.Email, user.Password)
	if err != nil {
		log.Printf("login authorize error: %v", err)
		c.WriteStatusCode(http.StatusUnauthorized)
		return nil
	}

	jsonToken := dto.NewToken(token)
	return c.WriteResponseJSON(http.StatusOK, jsonToken)
}

// getUserFromReqBody リクエストボディからユーザー情報を取得する
func (h *userHandler) getUserFromReqBody(body io.ReadCloser) (dto.User, error) {
	var user dto.User
	err := json.NewDecoder(body).Decode(&user)

	return user, err
}
