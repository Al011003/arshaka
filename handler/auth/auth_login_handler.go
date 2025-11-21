package handler

import (
	"net/http"

	req "backend/dto/request/auth"
	res "backend/dto/response"
	usecase "backend/usecase/auth"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	authUsecase usecase.LoginUsecase
}

func NewLoginHandler(a usecase.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		authUsecase: a,
	}
}

// LOGIN HANDLER
func (h *LoginHandler) Login(c *gin.Context) {
	var request req.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status:  "error",
			Message: "payload tidak valid",
		})
		return
	}

	resp, err := h.authUsecase.Login(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, res.BaseResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status:  "success",
		Message: "login berhasil",
		Data:    resp,
	})
}
