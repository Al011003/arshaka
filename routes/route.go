package routes

import (
	AdminHandler "backend/handler/admin"
	AngkatanMapalaHandler "backend/handler/angkatan_mapala"
	AuthHandler "backend/handler/auth"
	DeviceHandler "backend/handler/device_token"
	DataHandler "backend/handler/masterdata"
	superAdminHandler "backend/handler/superadmin"
	UserHandler "backend/handler/user"

	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	registerHandler *AuthHandler.RegisterHandler,
	loginHandler *AuthHandler.LoginHandler,
	userResetPasswordHandler *AuthHandler.PasswordResetHandler,
	changePasswordHandler *AuthHandler.UpdatePasswordHandler,
	fakultasHandler *DataHandler.FakultasHandler,
    jurusanHandler *DataHandler.JurusanHandler,
	angkatanMapalaHandler *AngkatanMapalaHandler.AngkatanMapalaHandler,
	userUpdateHandler *UserHandler.UserUpdateHandler,
	userProfileHandelr *UserHandler.UserProfileHandler,
	userChangeEmailHandler *UserHandler.UpdateEmailHandler,
	userPhotoPicHandler *UserHandler.UserPhotoHandler,


	adminHandler *AdminHandler.AdminUpdateHandler,
	superAdminUpdateHandler *superAdminHandler.SuperAdminSelfUpdateHandler,
	superAdminProfileHandler *superAdminHandler.SuperAdminProfileHandler,
	superAdminGetAllUserHandler *superAdminHandler.SuperAdminUserHandler,
	superAdminGetUserHandler *superAdminHandler.SuperAdminGetUserHandler,
	superAdminAccResetHandler *superAdminHandler.SuperAdminResetPasswordHandler,
	superAdminDeleteUserHandler *superAdminHandler.SuperAdminDeleteUserHandler,
	adminGetAllUserHandler *AdminHandler.AdminUserHandler,
	adminProfileHandler *AdminHandler.AdminProfileHandler,
	adminGetDetailUserHandler *AdminHandler.AdminGetUserHandler,
	adminResetPassHandler *AdminHandler.AdminForgotPasswordHandler,
	deviceTokenhandler *DeviceHandler.DeviceTokenHandler,
	

) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	auth.POST("/login", loginHandler.Login)
	auth.POST("/forgot-password", userResetPasswordHandler.RequestOTP)
	auth.POST("/forgot-password/change", userResetPasswordHandler.ResetPassword)
	mainRoute := r.Group("/api")
	mainRoute.Use(middleware.JWTAuth())
	
	userRoute := mainRoute.Group("/user")
	userRoute.Use(middleware.UserOnly())
	userUpdateRoute := userRoute.Group("/update")
		userUpdateRoute.PUT("/", userUpdateHandler.UpdateSelf)
		userPhotoRoute := userUpdateRoute.Group("/photo")
			userPhotoRoute.POST("/upload", userPhotoPicHandler.UpdatePhoto)
			userPhotoRoute.DELETE("/delete", userPhotoPicHandler.DeletePhoto)




	userRoute.GET("/profile", userProfileHandelr.GetProfile)
	userRoute.POST("/device-token", deviceTokenhandler.Save)
	userRoute.POST("/password", changePasswordHandler.UpdatePassword)
	userRoute.POST("/email", userChangeEmailHandler.UpdateEmail)
	
    
	
	adminRoute := mainRoute.Group("/admin")
	adminRoute.Use(middleware.AdminOnly())
		adminRoute.GET("/profile", adminProfileHandler.GetProfile)
    	adminRoute.POST("/register-user", registerHandler.RegisterUser)
		adminRoute.POST("/password", changePasswordHandler.UpdatePassword)
	
		resetPassword := adminRoute.Group("/reset-password")
			resetPassword.POST("/", adminResetPassHandler.RequestForgotPassword)
			resetPassword.POST("/delete", adminResetPassHandler.CancelForgotPassword)

		fakultas := adminRoute.Group("/fakultas")
			fakultas.POST("/", fakultasHandler.CreateFakultas)
			fakultas.GET("/", fakultasHandler.GetAllFakultas)
			fakultas.PUT("/:id", fakultasHandler.UpdateFakultas)
			fakultas.DELETE("/:id", fakultasHandler.DeleteFakultas)
    
		jurusan := adminRoute.Group("/jurusan")
			jurusan.POST("/", jurusanHandler.CreateJurusan)
			jurusan.GET("/", jurusanHandler.GetAllJurusan)
			jurusan.PUT("/:id", jurusanHandler.UpdateJurusan)
			jurusan.DELETE("/:id", jurusanHandler.DeleteJurusan)
			jurusan.GET("/fakultas/:fakultas_id", jurusanHandler.GetJurusanByFakultas)
			
		angkatan_mapala := adminRoute.Group("/angkatan-mapala")
			angkatan_mapala.POST("/", angkatanMapalaHandler.Create)
			angkatan_mapala.GET("/", angkatanMapalaHandler.GetAll)
			angkatan_mapala.PUT("/:id", angkatanMapalaHandler.Update)
			angkatan_mapala.DELETE("/:id", angkatanMapalaHandler.Delete)
    
		update := adminRoute.Group("/update")
			update.PUT("/", adminHandler.AdminUpdateSelf)

		user := adminRoute.Group("/user")
			user.PUT("/update/:id", adminHandler.AdminUpdateUser)
			user.GET("/", adminGetAllUserHandler.GetUsers)
			user.GET("/:id", adminGetDetailUserHandler.GetDetailUser)
		

	superAdminRoute := mainRoute.Group("/super-admin")
	superAdminRoute.Use(middleware.SuperAdminOnly())
	superAdminRoute.POST("/register-admin", registerHandler.RegisterAdmin)
	superAdminRoute.PUT("/update", superAdminUpdateHandler.SuperAdminUpdateSelf)
	superAdminRoute.GET("/profile", superAdminProfileHandler.GetProfile)
	superAdminRoute.POST("/password", changePasswordHandler.UpdatePassword)


	userDetailRoute := superAdminRoute.Group("/user")
	userDetailRoute.GET("/", superAdminGetAllUserHandler.GetUsers)
	userDetailRoute.GET("/:id", superAdminGetUserHandler.GetDetailUser)
	userDetailRoute.DELETE("/delete/:id", superAdminDeleteUserHandler.DeleteUser)


	resetPassRoute := superAdminRoute.Group("/reset-password")
			resetPassRoute.GET("/", superAdminAccResetHandler.GetAllRequests)
			resetPassRoute.POST("/approve/:resetID", superAdminAccResetHandler.ApproveReset)
	

	

	return r
}
