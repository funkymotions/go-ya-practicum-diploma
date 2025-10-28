package worker

import (
	"context"
	"sync"
	"time"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/interfaces"
)

type Poller struct {
	stopCh       chan struct{}
	doneCh       chan struct{}
	orderService interfaces.OrderService
	logger       interfaces.Logger
}

func NewPoller(
	stopCh chan struct{},
	doneCh chan struct{},
	orderService interfaces.OrderService,
	logger interfaces.Logger,
) *Poller {
	return &Poller{
		stopCh:       stopCh,
		doneCh:       doneCh,
		orderService: orderService,
		logger:       logger,
	}
}

func (p *Poller) Start() {
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		wg := &sync.WaitGroup{}
		ctx, cancelFunc := context.WithCancel(context.Background())
		for {
			select {
			case <-p.stopCh:
				p.logger.Infof("Stopping poller...\n")
				cancelFunc()
				p.logger.Infof("Cancelling async tasks...\n")
				wg.Wait()
				p.logger.Infof("Finished cancelling async tasks.\n")
				p.doneCh <- struct{}{}
				return
			case <-ticker.C:
				p.logger.Infof("Polling orders from accrual system...\n")
				// run in a separate goroutine to avoid blocking
				go p.orderService.EnqueueOrdersFromAccrualSystem(ctx, wg)
			}
		}
	}()
}
