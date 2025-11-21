package handler

import (
	"fmt"
	"net/http"

	req "backend/dto/request/auth"
	usecase "backend/usecase/auth"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
    authUsecase usecase.RegisterUsecase
}

func NewRegisterHandler(a usecase.RegisterUsecase) *RegisterHandler {
    return &RegisterHandler{
        authUsecase: a,
    }
}

// Register user/admin
func (h *RegisterHandler) RegisterUser(c *gin.Context) {
    var request req.RegisterRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    // Ambil role creator dari JWT middleware
    creatorRole := c.GetString("role") // pastikan middleware set "role" di context
    fmt.Println("ROLE DARI TOKEN:", creatorRole)

    resp, err := h.authUsecase.RegisterUser(request, creatorRole)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "user berhasil dibuat",
        "data":    resp,
    })
}
