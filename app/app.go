package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/config"
	handler "backend/handler/auth"
	"backend/model"
	"backend/repo"
	"backend/routes"

	authUC "backend/usecase/auth"
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
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	// Init Repos
	authRepo := repo.NewAuthRepository(db)

	// Init Usecases
	registerUC := authUC.NewRegisterUsecase(authRepo)
	loginUC := authUC.NewLoginUsecase(authRepo)

	// Init Handlers
	registerHandler := handler.NewRegisterHandler(registerUC)
	loginHandler := handler.NewLoginHandler(loginUC)

	// Init Router + Inject handlers
	r := routes.SetupRouter(
		registerHandler,
		loginHandler,
	)

	return &App{
		DB:     db,
		Router: r,
	}, nil
}
