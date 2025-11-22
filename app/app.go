package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/config"
	authhandler "backend/handler/auth"
	datahandler "backend/handler/masterdata"
	"backend/model"
	"backend/repo"
	"backend/routes"

	authUC "backend/usecase/auth"
	dataUC "backend/usecase/masterdata"
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
	if err := db.AutoMigrate(&model.User{}, &model.Fakultas{}, &model.Jurusan{}, &model.AngkatanMapala{}); err != nil {
		return nil, err
	}

	// Init Repos
	authRepo := repo.NewAuthRepository(db)
	jurusanRepo := repo.NewJurusanRepo(db)
	fakultasRepo := repo.NewFakultasRepo(db)
	angaktanMapalaRepo := repo.NewAngkatanMapalaRepo(db)

	// Init Usecases
	registerUC := authUC.NewRegisterUsecase(authRepo, fakultasRepo, jurusanRepo, angaktanMapalaRepo)
	loginUC := authUC.NewLoginUsecase(authRepo)
	fakultasUC := dataUC.NewFakultasUsecase(fakultasRepo)
	jurusanUC := dataUC.NewJurusanUsecase(jurusanRepo, fakultasRepo)
	angkatanMapalaUC := dataUC.NewAngkatanMapalaUsecase(angaktanMapalaRepo)

	// Init Handlers
	registerHandler := authhandler.NewRegisterHandler(registerUC)
	loginHandler := authhandler.NewLoginHandler(loginUC)
	fakultasHandler := datahandler.NewFakultasHandler(fakultasUC)
	jurusanHandler := datahandler.NewJurusanHandler(jurusanUC)
	angkatanMapalaHandler := datahandler.NewAngkatanMapalaHandler(angkatanMapalaUC)




	// Init Router + Inject handlers
	r := routes.SetupRouter(
		registerHandler,
		loginHandler,
		fakultasHandler,
		jurusanHandler,
		angkatanMapalaHandler,
	)

	return &App{
		DB:     db,
		Router: r,
	}, nil
}
