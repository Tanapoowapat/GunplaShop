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
	ordershandlers "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersHandlers"
	ordersrepositories "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersRepositories"
	ordersusecase "github.com/Tanapoowapat/GunplaShop/modules/orders/ordersUsecase"
	productshandlers "github.com/Tanapoowapat/GunplaShop/modules/products/productsHandlers"
	productsrepositories "github.com/Tanapoowapat/GunplaShop/modules/products/productsRepositories"
	productsusecase "github.com/Tanapoowapat/GunplaShop/modules/products/productsUsercase"
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
	ProductsModule()
	OrdersModule()
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

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorization(2), handler.UploadImageLocal)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorization(2), handler.DeleteImagesLocal)
}

func (m *moduleFactory) ProductsModule() {
	fileUsecase := filesusecase.NewFileUsecase(m.server.cfg)
	repo := productsrepositories.NewProductRepositories(m.server.db, m.server.cfg, fileUsecase)
	usecase := productsusecase.NewProductsUsecase(repo)
	handler := productshandlers.NewProductsHandler(m.server.cfg, usecase, fileUsecase)

	router := m.router.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorization(2), handler.AddProducts)
	router.Patch("/", m.mid.JwtAuth(), m.mid.Authorization(2), handler.UpdateProducts)

	router.Get("/", m.mid.CheckApiKey(), handler.FindProducts)
	router.Get("/:product_id", m.mid.CheckApiKey(), handler.FindOneProduct)

	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorization(2), handler.DeleteProducts)

}

func (m *moduleFactory) OrdersModule() {
	fileUsecase := filesusecase.NewFileUsecase(m.server.cfg)
	productsRepo := productsrepositories.NewProductRepositories(m.server.db, m.server.cfg, fileUsecase)

	repo := ordersrepositories.NewOrdersRepositories(m.server.db)
	usecase := ordersusecase.NewOrdersUsecase(repo, productsRepo)
	handler := ordershandlers.NewOrdersHandlers(usecase, m.server.cfg)

	router := m.router.Group("/orders")

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorization(2), handler.FindOrder)
	router.Get("/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.FindOnceOrders)
}
