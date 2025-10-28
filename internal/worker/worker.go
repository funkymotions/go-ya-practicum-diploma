package worker

import (
	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type WorkerConfig struct {
	BufferSize   int
	StopCh       chan struct{}
	DoneCh       chan struct{}
	Queue        chan *model.Order
	OrderService interfaces.OrderService
	Logger       interfaces.Logger
}

type Worker struct {
	queue        chan *model.Order
	stopCh       chan struct{}
	doneCh       chan struct{}
	orderService interfaces.OrderService
	logger       interfaces.Logger
}

func NewWorker(
	cfg *WorkerConfig,
) *Worker {
	return &Worker{
		queue:        cfg.Queue,
		stopCh:       cfg.StopCh,
		doneCh:       cfg.DoneCh,
		orderService: cfg.OrderService,
		logger:       cfg.Logger,
	}
}

func (w *Worker) Start() {
	go func() {
		// hardcoded number of workers for simplicity
		for i := 0; i < 5; i++ {
			go w.ProcessOrderWorker()
		}
	}()
}

func (w *Worker) ProcessOrderWorker() {
	for {
		select {
		case order := <-w.queue:
			w.logger.Infof("Processing order...%+v\n", order)
			if err := w.orderService.ProcessOrderFromAccrual(order); err != nil {
				w.logger.Errorf("Error processing order with ID = %s, err: %v\n", order.ID, err)
			}
			// Process the order
		case <-w.stopCh:
			w.logger.Infof("Stopping worker...\n")
			// TODO: Clean up resources if needed
			w.doneCh <- struct{}{}
			return
		}
	}
}
