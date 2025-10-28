package app

import (
	"net/http"
	"net/url"
	"time"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/config"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/handler"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/infrastructure/drivers/postgres"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/infrastructure/drivers/rest"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/middleware"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/repository"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/service"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/utils"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/worker"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	server     *http.Server
	engine     *chi.Mux
	stopCh     chan struct{}
	doneCh     chan struct{}
	ordersChan chan *model.Order
	worker     *worker.Worker
	poller     *worker.Poller
	logger     *zap.SugaredLogger
}

func NewApp() *App {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return &App{
		server: &http.Server{},
		engine: chi.NewRouter(),
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}, 5),
		logger: logger.Sugar(),
	}
}

func (a *App) Init() error {
	config, err := config.NewConfig()
	if err != nil {
		return err
	}
	driver, err := postgres.NewSQLDriver(config.DBConfig)
	if err != nil {
		return err
	}

	// this data chan will be used to send orders to worker pool
	const numOfWorkers = 5
	a.ordersChan = make(chan *model.Order, numOfWorkers)
	// HTTP client for accrual system
	accrualHTTPClient := &http.Client{
		Timeout: time.Second * 5,
	}

	accrualBaseURL, err := url.Parse(config.AppConfig.AccrualSystemAddress)
	if err != nil {
		return err
	}
	accrualRestClient := rest.NewRESTClient(
		&rest.RESTClientConfig{
			HTTPClient: accrualHTTPClient,
			BaseURL:    accrualBaseURL,
		},
	)

	// 10 concurrect requests
	s := utils.NewSemaphore(10)

	// repositories
	userRepo := repository.NewUserRepository(driver.DB)
	orderRepo := repository.NewOrderRepository(driver.DB)
	withdrawalRepo := repository.NewWithdrawalRepository(driver.DB)

	// services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(orderRepo, withdrawalRepo)
	orderService := service.NewOrderService(
		&service.OrderServiceConf{
			OrderRepo:      orderRepo,
			AccrualClient:  accrualRestClient,
			OrderQueue:     a.ordersChan,
			RequestLimiter: s,
			Logger:         a.logger,
		},
	)
	withdrawalService := service.NewWithdrawalService(withdrawalRepo)

	// middlewares
	authMiddleware := middleware.NewAuthUserMiddleware(userRepo)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	accountHandler := handler.NewAccountHandler(accountService)
	orderHandler := handler.NewOrderHandler(orderService)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalService)
	a.server = &http.Server{
		Addr:    config.AppConfig.AppAddress,
		Handler: a.engine,
	}

	// attach routes to engine
	userHandler.RegisterRoutes(a.engine)
	accountHandler.RegisterRoutes(a.engine, authMiddleware)
	orderHandler.RegisterRoutes(a.engine, authMiddleware)
	withdrawalHandler.RegisterRoutes(a.engine, authMiddleware)

	// number of workers is hardcoded for simplicity, can be moved to config if needed of course
	wConf := &worker.WorkerConfig{
		BufferSize:   5,
		StopCh:       a.stopCh,
		DoneCh:       a.doneCh,
		Queue:        a.ordersChan,
		OrderService: orderService,
		Logger:       a.logger,
	}
	a.worker = worker.NewWorker(wConf)

	// pollers
	a.poller = worker.NewPoller(
		a.stopCh,
		a.doneCh,
		orderService,
		a.logger,
	)

	return nil
}

func (a *App) Start() error {
	a.worker.Start()
	a.poller.Start()
	a.logger.Infof("starting application on %s...", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown() {
	close(a.stopCh)
	a.logger.Infof("waiting for workers and pollers to stop...")
	for i := 0; i < 5; i++ {
		<-a.doneCh
	}
	a.logger.Infof("all workers and pollers stopped")
}
