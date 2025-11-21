package routes

import (
	handler "backend/handler/auth"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	registerHandler *handler.RegisterHandler,
	loginHandler *handler.LoginHandler,
) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	auth.POST("/login", loginHandler.Login)
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth())   // ← WAJIB
    admin.POST("/register", registerHandler.RegisterUser)

	

	return r
}
