package monitorhandlers

import (
	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/entities"
	"github.com/Tanapoowapat/GunplaShop/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandlers interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandlers struct {
	cfg config.IConfig
}

func NewMonitorHandlers(cfg config.IConfig) IMonitorHandlers {
	return &monitorHandlers{
		cfg: cfg,
	}
}

func (mh *monitorHandlers) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    mh.cfg.App().Name(),
		Version: mh.cfg.App().Version(),
	}
	return entities.NewResponse(c).Sucess(fiber.StatusOK, res).Res()
}
