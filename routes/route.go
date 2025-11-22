package routes

import (
	AuthHandler "backend/handler/auth"
	DataHandler "backend/handler/masterdata"

	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	registerHandler *AuthHandler.RegisterHandler,
	loginHandler *AuthHandler.LoginHandler,
	fakultasHandler *DataHandler.FakultasHandler,
    jurusanHandler *DataHandler.JurusanHandler,
	angkatanMapalaHandler *DataHandler.AngkatanMapalaHandler,
) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	auth.POST("/login", loginHandler.Login)
	mainRoute := r.Group("/api")
	mainRoute.Use(middleware.JWTAuth())
	
	userRoute := mainRoute.Group("/user")
	userRoute.Use(middleware.UserOnly())
	
	adminRoute := mainRoute.Group("/admin")
	adminRoute.Use(middleware.AdminOnly())
    adminRoute.POST("/register", registerHandler.RegisterUser)
	fakultas := adminRoute.Group("/fakultas")
    {
        fakultas.POST("/", fakultasHandler.CreateFakultas)
        fakultas.GET("/", fakultasHandler.GetAllFakultas)
        fakultas.PUT("/:id", fakultasHandler.UpdateFakultas)
        fakultas.DELETE("/:id", fakultasHandler.DeleteFakultas)
    }
	jurusan := adminRoute.Group("/jurusan")
    {
        jurusan.POST("/", jurusanHandler.CreateJurusan)
        jurusan.GET("/", jurusanHandler.GetAllJurusan)
        jurusan.PUT("/:id", jurusanHandler.UpdateJurusan)
        jurusan.DELETE("/:id", jurusanHandler.DeleteJurusan)
		jurusan.GET("/fakultas/:fakultas_id", jurusanHandler.GetJurusanByFakultas)
    }
	angkatan_mapala := adminRoute.Group("/angkatan-mapala")
    {
        angkatan_mapala.POST("/", angkatanMapalaHandler.Create)
        angkatan_mapala.GET("/", angkatanMapalaHandler.GetAll)
        angkatan_mapala.PUT("/:id", angkatanMapalaHandler.Update)
        angkatan_mapala.DELETE("/:id", angkatanMapalaHandler.Delete)
    }

	

	return r
}
