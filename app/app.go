package app

import (
	service "backend/email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/config"
	adminhandler "backend/handler/admin"
	angkatanmapalahandler "backend/handler/angkatan_mapala"
	authhandler "backend/handler/auth"
	baranghandler "backend/handler/barang"
	devicehandler "backend/handler/device_token"
	datahandler "backend/handler/masterdata"
	superadminhandler "backend/handler/superadmin"
	userhandler "backend/handler/user"
	"backend/model"
	"backend/repo"
	"backend/routes"
	superAdminResetPassUC "backend/usecase/super_admin/reset_password"

	adminProfileUC "backend/usecase/admin/profile"
	adminResetPasswordUC "backend/usecase/admin/reset_password"
	adminUpdateUC "backend/usecase/admin/update"

	angkatanMapalaUC "backend/usecase/angkatan_mapala"
	authUC "backend/usecase/auth"
	deviceUC "backend/usecase/device_token"
	dataUC "backend/usecase/masterdata"

	barangUC "backend/usecase/barang"
	superAdminDeleteUC "backend/usecase/super_admin/delete"
	superAdminGPUC "backend/usecase/super_admin/getprofile"
	superAdminUpdateUC "backend/usecase/super_admin/update"
	userProfileUC "backend/usecase/user/profile"
	userUpdateUC "backend/usecase/user/update"
	"backend/utils"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func NewApp() (*App, error) {
	// Init logger
	logFile := utils.InitLogger()
	defer logFile.Close()

	utils.Log.Info("App starting...")

	// Init DB
	db := config.InitDB()

	// Auto migrate models
	if err := db.AutoMigrate(&model.User{}, &model.Fakultas{}, &model.Jurusan{}, &model.AngkatanMapala{}, &model.DeviceToken{}, &model.PasswordResetRequest{}, &model.PasswordReset{}, &model.Barang{}); err != nil {
		return nil, err
	}

	emailService := service.NewSMTPEmailService()

	// Init Repos
	authRepo := repo.NewAuthRepository(db)
	jurusanRepo := repo.NewJurusanRepo(db)
	fakultasRepo := repo.NewFakultasRepo(db)
	angaktanMapalaRepo := repo.NewAngkatanMapalaRepo(db)
	userRepo := repo.NewUserRepo(db)
	deviceRepo := repo.NewDeviceTokenRepo(db)
	PasswordResetRepo := repo.NewPasswordResetRepo(db)
	PasswordResetRepository := repo.NewPasswordResetRepository(db)
	barangRepo := repo.NewBarangRepository(db)

	// Init Usecases
	//autj
	registerUC := authUC.NewRegisterUsecase(authRepo, fakultasRepo, jurusanRepo, angaktanMapalaRepo)
	loginUC := authUC.NewLoginUsecase(authRepo)
	userResetUC := authUC.NewPasswordResetUsecase(userRepo, PasswordResetRepository, emailService)
	changePassUC := authUC.NewUpdatePasswordUsecase(userRepo)
	//masterdata
	fakultasUC := dataUC.NewFakultasUsecase(fakultasRepo)
	jurusanUC := dataUC.NewJurusanUsecase(jurusanRepo, fakultasRepo)
	angkatanMapalaUC := angkatanMapalaUC.NewAngkatanMapalaUsecase(angaktanMapalaRepo)
	//user
	userSelfUpdatefUC := userUpdateUC.NewUserSelfUsecase(userRepo)
	userGetProfileUc := userProfileUC.NewUserProfileUsecase(userRepo)
	userChangeEmailUC := userUpdateUC.NewUpdateEmailUsecase(userRepo)
	userProfilePicUC := userProfileUC.NewUserPhotoUsecase(userRepo)
	//admin
	adminSelfUC := adminUpdateUC.NewAdminSelfUpdateUsecase(userRepo)
	adminUpdateUC := adminUpdateUC.NewAdminUpdateUserUsecase(authRepo, fakultasRepo, jurusanRepo, angaktanMapalaRepo)
	adminGetAllUserUC := adminProfileUC.NewAdminGetUserUsecase(userRepo)
	adminSelfProfileUC := adminProfileUC.NewAdminProfileUC(userRepo)
	adminGetDetailUserUC := adminProfileUC.NewAdminUserDetailUsecase(userRepo)
	adminResetPassUC := adminResetPasswordUC.NewAdminResetPasswordUsecase(PasswordResetRepo, userRepo)
	//superadmin
	superAdminUpdateUC := superAdminUpdateUC.NewSuperAdminSelfUpdateUsecase(userRepo)
	superAdminProfileUC := superAdminGPUC.NewSuperAdminProfileUC(userRepo)
	superAdminGetUserUC := superAdminGPUC.NewSuperAdminGetUserUsecase(userRepo)
	superAdminGetDetailUC := superAdminGPUC.NewUserDetailUsecase(userRepo)
	superAdminAccResetUC := superAdminResetPassUC.NewSuperAdminAccResetPasswordUsecase(PasswordResetRepo, userRepo)
	superAdminDeleteUserUC := superAdminDeleteUC.NewSuperAdminDeleteUserUsecase(userRepo)
	//device
	deviceTokenUC := deviceUC.NewSaveDeviceTokenUC(deviceRepo)
	//barang
	barangCrudUC := barangUC.NewBarangUseCase(barangRepo)
	barangPhotoUC := barangUC.NewBarangPhotoUsecase(barangRepo)



	// Init Handlers
	//auth
	registerHandler := authhandler.NewRegisterHandler(registerUC)
	loginHandler := authhandler.NewLoginHandler(loginUC)
	userResetPasswordHandler := authhandler.NewPasswordResetHandler(userResetUC)
	changePasswordHandler := authhandler.NewUpdatePasswordHandler(changePassUC)
	//masterdata
	fakultasHandler := datahandler.NewFakultasHandler(fakultasUC)
	jurusanHandler := datahandler.NewJurusanHandler(jurusanUC)
	angkatanMapalaHandler := angkatanmapalahandler.NewAngkatanMapalaHandler(angkatanMapalaUC)
	//user
	userUpdateHandler := userhandler.NewUserUpdateHandler(userSelfUpdatefUC)
	userProfileHandler := userhandler.NewUserProfileHandler(userGetProfileUc)
	userChangeEmailHandler := userhandler.NewUpdateEmailHandler(userChangeEmailUC)
	userPhotoPicHandler := userhandler.NewUserPhotoHandler(userProfilePicUC)

	//admin
	adminHandler := adminhandler.NewAdminUpdateHandler(adminSelfUC, adminUpdateUC)
	adminResetPassHandler := adminhandler.NewUserForgotPasswordHandler(adminResetPassUC)
	//superadmin
	superAdminUpdateHandler := superadminhandler.NewSuperAdminSelfUpdateHandler(superAdminUpdateUC)
	superAdminProfileHandler := superadminhandler.NewSuperAdminProfileHandler(superAdminProfileUC)
	superAdminGetAllUserHandler := superadminhandler.NewSuperAdminUserHandler(superAdminGetUserUC)
	superAdminGetUserHandler := superadminhandler.NewSuperAdminGetUserHandler(superAdminGetDetailUC)
	superAdminAccResetHandler := superadminhandler.NewAdminResetPasswordHandler(superAdminAccResetUC)
	superAdminDeleteUserHandler := superadminhandler.NewSuperAdminDeleteUserHandler(superAdminDeleteUserUC)
	adminGetAllUserHandler := adminhandler.NewAdminUserHandler(adminGetAllUserUC)
	adminProfileHandler := adminhandler.NewAdminProfileHandler(adminSelfProfileUC)
	adminGetDetailUserHandler := adminhandler.NewAdminGetUserHandler(adminGetDetailUserUC)
	//device
	deviceTokenhandler:= devicehandler.NewDeviceTokenHandler(deviceTokenUC)
	//barang
	barangCrudHandler := baranghandler.NewBarangHandler(barangCrudUC)
	barangPhotoHandler := baranghandler.NewBarangPhotoHandler(barangPhotoUC)
	






	// Init Router + Inject handlers
	r := routes.SetupRouter(
		registerHandler,
		loginHandler,
		userResetPasswordHandler,
		changePasswordHandler,
		fakultasHandler,
		jurusanHandler,
		angkatanMapalaHandler,
		userUpdateHandler,
		userProfileHandler,
		userChangeEmailHandler,
		userPhotoPicHandler,
		adminHandler,
		superAdminUpdateHandler,
		superAdminProfileHandler,
		superAdminGetAllUserHandler,
		superAdminGetUserHandler,
		superAdminAccResetHandler,
		superAdminDeleteUserHandler,
		adminGetAllUserHandler,
		adminProfileHandler,
		adminGetDetailUserHandler,
		adminResetPassHandler,
		deviceTokenhandler,
		barangCrudHandler,
		barangPhotoHandler,
		
		
		
	)

	return &App{
		DB:     db,
		Router: r,
	}, nil
}
