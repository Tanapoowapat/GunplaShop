package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	//Middlewares
	middlewares := NewMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())
	//Modules
	v1 := s.app.Group("/v1")
	modules := NewModuleFactory(v1, s, middlewares)
	modules.MonitorModule()
	modules.UserMoudle()
	modules.AppinfoModule()
	modules.FileModule()
	modules.ProductsModule()
	modules.OrdersModule()
	s.app.Use(middlewares.RouterCheck())

	//Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Server is shutting down...")
		_ = s.app.Shutdown()
	}()

	//Listen to host:port
	log.Printf("server is running on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())

}
