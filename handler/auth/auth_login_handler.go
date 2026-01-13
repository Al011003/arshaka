package handler

import (
	"net/http"

	req "backend/dto/request/auth"
	res "backend/dto/response/common"
	usecase "backend/usecase/auth"
	"backend/utils"

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

// Login godoc
// @Summary      User login
// @Description  Login menggunakan email dan password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      auth.LoginRequest  true  "Login credentials"
// @Success      200      {object}  response.BaseResponse  "Login berhasil"
// @Failure      400      {object}  response.BaseResponse  "Bad request"
// @Failure      401      {object}  response.BaseResponse  "Unauthorized"
// @Router       /auth/login [post]
func (h *LoginHandler) Login(c *gin.Context) {
	var request req.LoginRequest
	

	if err := request.BindAndValidate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
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

	utils.Success(c, resp, "login berhasil")
}