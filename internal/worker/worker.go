package worker

import (
	"fmt"

	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type Worker struct {
	queue   chan *model.Order
	results chan error
	stopCh  chan struct{}
	doneCh  chan struct{}
}

func NewWorker(bufferSize int, stopCh chan struct{}, doneCh chan struct{}, queue chan *model.Order) *Worker {
	return &Worker{
		queue:   queue,
		results: make(chan error, bufferSize),
		stopCh:  stopCh,
		doneCh:  doneCh,
	}
}

func (w *Worker) Start() {
	go func() {
		for i := 0; i < 5; i++ {
			go w.ProcessOrderWorker()
		}
	}()
}

func (w *Worker) ProcessOrderWorker() {
	for {
		select {
		case order := <-w.queue:
			fmt.Printf("Processing order... %+v\n", order)
			// Process the order
		case <-w.results:
			fmt.Printf("Handling order error...\n")
			// Handle the result
		case <-w.stopCh:
			fmt.Printf("Stopping worker...\n")
			// TODO: Clean up resources if needed
			w.doneCh <- struct{}{}
			return
		}
	}
}
