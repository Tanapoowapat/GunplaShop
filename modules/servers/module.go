package servers

import (
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
	router.Post("/signup", handlers.SignUpCustomer)
	router.Post("/signin", handlers.SignIn)
	router.Post("/signout", handlers.SignOut)
	router.Post("/refresh", handlers.RefreshPassport)
	router.Post("/signup-admin", handlers.SignUpAdmin)

	//Get
	router.Get("/:userId", m.mid.JwtAuth(), m.mid.ParamsCheck(), handlers.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorization(2), handlers.GenaerateAdminToken)
}

func (mf *moduleFactory) MonitorModule() {
	handlers := monitorhandlers.NewMonitorHandlers(mf.server.cfg)
	mf.router.Get("/", handlers.HealthCheck)
}
