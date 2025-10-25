package app

import (
	"log"
	"net/http"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/config"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/handler"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/infrastructure/drivers/postgres"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/repository"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/worker"
	"github.com/go-chi/chi/v5"
)

type App struct {
	server     *http.Server
	engine     *chi.Mux
	stopCh     chan struct{}
	doneCh     chan struct{}
	ordersChan chan *model.Order
	worker     *worker.Worker
}

func NewApp() *App {
	return &App{
		server: &http.Server{},
		engine: chi.NewRouter(),
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}, 5),
	}
}

func (a *App) Init() error {
	appConf, err := config.NewAppConfig()
	if err != nil {
		return err
	}
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		return err
	}
	driver, err := postgres.NewSQLDriver(dbConfig)
	if err != nil {
		return err
	}

	// repositories
	userRepo := repository.NewUserRepository(driver.DB)
	orderRepo := repository.NewOrderRepository(driver.DB)
	withdrawalRepo := repository.NewWithdrawalRepository(driver.DB)

	// services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(orderRepo, withdrawalRepo)
	orderService := service.NewOrderService(orderRepo)
	withdrawalService := service.NewWithdrawalService(withdrawalRepo)

	// middlewares
	authMiddleware := middleware.NewAuthUserMiddleware(userRepo)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	accountHandler := handler.NewAccountHandler(accountService)
	orderHandler := handler.NewOrderHandler(orderService)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalService)
	a.server = &http.Server{
		Addr:    appConf.AppAddress,
		Handler: a.engine,
	}

	// attach routes to engine
	userHandler.RegisterRoutes(a.engine)
	accountHandler.RegisterRoutes(a.engine, authMiddleware)
	orderHandler.RegisterRoutes(a.engine, authMiddleware)
	withdrawalHandler.RegisterRoutes(a.engine, authMiddleware)

	// this data chan will be used to send orders to worker pool
	const numOfWorkers = 10
	a.ordersChan = make(chan *model.Order, numOfWorkers)

	// worker
	// number of workers is hardcoded for simplicity, can be moved to config if needed of course
	a.worker = worker.NewWorker(numOfWorkers, a.stopCh, a.doneCh, a.ordersChan)

	return nil
}

func (a *App) Start() error {
	a.worker.Start()
	log.Printf("starting application on %s...", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown() {
	close(a.stopCh)
	log.Printf("shutting down application...")
	for i := 0; i < 5; i++ {
		<-a.doneCh
	}
	log.Printf("application stopped")
}
