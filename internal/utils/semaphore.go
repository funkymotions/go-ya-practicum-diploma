package utils

type Semaphore struct {
	syncChan chan struct{}
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		syncChan: make(chan struct{}, limit),
	}
}

func (s *Semaphore) Acquire() {
	s.syncChan <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.syncChan
}
