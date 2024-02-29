package servers

import (
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresHandlers"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresRepositories"
	"github.com/Tanapoowapat/GunplaShop/modules/middlewares/middlewaresUsecase"
	monitorhandlers "github.com/Tanapoowapat/GunplaShop/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
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

func (mf *moduleFactory) MonitorModule() {
	handlers := monitorhandlers.NewMonitorHandlers(mf.server.cfg)
	mf.router.Get("/", handlers.HealthCheck)
}
