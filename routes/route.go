package routes

import (
	AdminHandler "backend/handler/admin"
	AngkatanMapalaHandler "backend/handler/angkatan_mapala"
	AuthHandler "backend/handler/auth"
	BarangHandler "backend/handler/barang"
	CartHandler "backend/handler/cart"
	DeviceHandler "backend/handler/device_token"
	LoanHandler "backend/handler/loan"
	DataHandler "backend/handler/masterdata"
	superAdminHandler "backend/handler/superadmin"
	UserHandler "backend/handler/user"

	"backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	barangCrudHandler *BarangHandler.BarangHandler,
	barangPhotoHandler *BarangHandler.BarangPhotoHandler,
	barangCheckHandler * BarangHandler.BarangAvailabilityHandler,
	
	cartHandler *CartHandler.CartHandler,

	loanUserHandler *LoanHandler.LoanUserHandler,
	loanUserGetLoanHandler *LoanHandler.LoanUserGetHandler,

) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	AllowCredentials: false,
}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	auth.POST("/login", loginHandler.Login)
	auth.POST("/request-otp", userResetPasswordHandler.RequestOTP)
	auth.POST("/verify-reset-otp", userResetPasswordHandler.VerifyResetOTP)
	auth.POST("/reset-password", userResetPasswordHandler.ResetPassword)
		adminResetPassword := auth.Group("/admin-reset")
			adminResetPassword.POST("/", adminResetPassHandler.RequestForgotPassword)
	mainRoute := r.Group("/api")
	mainRoute.Use(middleware.JWTAuth())
	

	//userr
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

	barangU := userRoute.Group("/barang")
		barangU.GET("/", barangCrudHandler.GetAll)
		barangU.GET("/:id", barangCrudHandler.GetByID)
		barangU.GET("/availabilitycheck/:barang_id", barangCheckHandler.GetCalendar)

	cart := userRoute.Group("/cart")
		cart.GET("", cartHandler.GetMyCart)              // Get my cart
		cart.GET("/count", cartHandler.GetCartItemCount) // Get cart item count
		cart.POST("", cartHandler.AddToCart)             // Add to cart
		cart.PUT("/:id", cartHandler.UpdateCartItem)     // Update cart item
		cart.DELETE("/:id", cartHandler.RemoveFromCart)  // Remove from cart
		cart.DELETE("", cartHandler.ClearCart)

	userLoan := userRoute.Group("/loan")
		userLoan.POST("", loanUserHandler.Create)
		userLoan.PUT("/:loan_code", loanUserHandler.UpdateHeader)
		userLoan.PUT("/:loan_code/items", loanUserHandler.UpdateItems)
		userLoan.DELETE("/:loan_code", loanUserHandler.Cancel)

		userLoan.GET("", loanUserHandler.GetMyLoans)
		userLoan.GET("/:loan_code", loanUserHandler.GetDetail)    
	
	
    
	
	adminRoute := mainRoute.Group("/admin")
	adminRoute.Use(middleware.AdminOnly())
		adminRoute.GET("/profile", adminProfileHandler.GetProfile)
    	adminRoute.POST("/register-user", registerHandler.RegisterUser)
		adminRoute.POST("/password", changePasswordHandler.UpdatePassword)
	


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
		
		barang := adminRoute.Group("/barang")
			// CRUD Barang
			barang.POST("", barangCrudHandler.Create)           // POST /admin/barang
			barang.GET("", barangCrudHandler.GetAll)            // GET /admin/barang
			barang.GET("/:id", barangCrudHandler.GetByID)       // GET /admin/barang/:id
			barang.PUT("/:id", barangCrudHandler.Update)        // PUT /admin/barang/:id
			barang.DELETE("/:id", barangCrudHandler.Delete)     // DELETE /admin/barang/:id
			// Photo Management
			barang.POST("/:id/photo", barangPhotoHandler.UpdatePhoto)     // POST /admin/barang/:id/photo
			barang.DELETE("/:id/photo", barangPhotoHandler.DeletePhoto)   // DELETE /admin/barang/:id/photo

		

	superAdminRoute := mainRoute.Group("/super-admin")
	superAdminRoute.Use(middleware.SuperAdminOnly())
	superAdminRoute.POST("/register-user", registerHandler.RegisterUser)
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
			resetPassRoute.POST("/cancel/:resetID", superAdminAccResetHandler.CancelReset)
	

	return r
}
