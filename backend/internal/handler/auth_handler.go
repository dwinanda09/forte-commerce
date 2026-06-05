package handler

import (
	"github.com/dwinanda09/forte-commerce/internal/usecase"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, 400, "ERR-HDL-001", "Invalid request body")
	}

	token, err := h.uc.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return util.Fail(c, 401, "ERR-HDL-002", "Invalid credentials")
	}

	return util.OK(c, LoginResponse{Token: token})
}
