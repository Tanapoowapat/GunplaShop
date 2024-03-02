package servers

import (
	appinfohandlers "github.com/Tanapoowapat/GunplaShop/modules/appinfo/appinfoHandlers"
	appinforepositories "github.com/Tanapoowapat/GunplaShop/modules/appinfo/appinfoRepositories"
	appinfousecase "github.com/Tanapoowapat/GunplaShop/modules/appinfo/appinfoUsecase"
	filehandler "github.com/Tanapoowapat/GunplaShop/modules/file/fileHandler"
	filesusecase "github.com/Tanapoowapat/GunplaShop/modules/file/filesUsecase"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresHandlers"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresRepositories"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresUsecase"
	monitorhandlers "github.com/Tanapoowapat/GunplaShop/modules/monitor/monitorHandlers"
	"github.com/Tanapoowapat/GunplaShop/modules/users/usersHandlers"
	"github.com/Tanapoowapat/GunplaShop/modules/users/usersRepositories"
	"github.com/Tanapoowapat/GunplaShop/modules/users/usersUsecase"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UserMoudle()
	AppinfoModule()
	FileModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.IMiddlewaresHandlers
}

func NewModuleFactory(router fiber.Router, server *server, mid middlewaresHandlers.IMiddlewaresHandlers) IModuleFactory {
	return &moduleFactory{
		router: router,
		server: server,
		mid:    mid,
	}
}

func NewMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandlers {
	repo := middlewaresRepositories.NewMiddlewaresRepositories(s.db)
	usecase := middlewaresUsecase.NewMiddlewaresUsecase(repo)
	return middlewaresHandlers.NewMiddlewaresHandlers(s.cfg, usecase)

}

func (m *moduleFactory) UserMoudle() {
	repo := usersRepositories.UsersRepositories(m.server.db)
	usecase := usersUsecase.UsersUsecase(m.server.cfg, repo)
	handlers := usersHandlers.NewUsersHandlers(m.server.cfg, usecase)

	router := m.router.Group("/users")

	//Post
	router.Post("/signup", m.mid.CheckApiKey(), handlers.SignUpCustomer)
	router.Post("/signin", m.mid.CheckApiKey(), handlers.SignIn)
	router.Post("/signout", m.mid.CheckApiKey(), handlers.SignOut)
	router.Post("/refresh", m.mid.CheckApiKey(), handlers.RefreshPassport)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorization(2), handlers.SignUpAdmin)

	//Get
	router.Get("/:userId", m.mid.JwtAuth(), m.mid.ParamsCheck(), handlers.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorization(2), handlers.GenaerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	appinfo_repo := appinforepositories.AppinfoRepositories(m.server.db)
	appinfo_usecase := appinfousecase.AppinfoRepositories(appinfo_repo)
	appinfo_handlers := appinfohandlers.NewAppinfoHandlers(m.server.cfg, appinfo_usecase)

	router := m.router.Group("/appinfo")
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorization(2), appinfo_handlers.GenerateApiKey)

	router.Get("/categories", m.mid.CheckApiKey(), appinfo_handlers.FindCategory)
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorization(2), appinfo_handlers.InsertCategory)

	router.Delete("/categories/:categoryId", m.mid.JwtAuth(), m.mid.Authorization(2), appinfo_handlers.RemoveCategory)
}

func (m *moduleFactory) MonitorModule() {
	handlers := monitorhandlers.NewMonitorHandlers(m.server.cfg)
	m.router.Get("/", handlers.HealthCheck)
}

func (m *moduleFactory) FileModule() {
	usecase := filesusecase.NewFileUsecase(m.server.cfg)
	handler := filehandler.NewfileHandler(m.server.cfg, usecase)

	router := m.router.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorization(2), handler.UploadFile)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorization(2), handler.DeleteFile)
}
