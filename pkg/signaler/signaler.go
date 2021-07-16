package signaler

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SignalerService interface {
	WaitForSignal(func(error))
}

type SignalerOption func(*Signaler)

type SignalHandler func(func()) error

type Signaler struct {
	wg            sync.WaitGroup
	onInterrupt   []SignalHandler
	onHangup      []SignalHandler
	onTermination []SignalHandler
}

func WithOnInterrupt(handlers ...SignalHandler) SignalerOption {
	return func(s *Signaler) {
		s.onInterrupt = append(s.onInterrupt, handlers...)
	}
}

func WithOnHangup(handlers ...SignalHandler) SignalerOption {
	return func(s *Signaler) {
		s.onHangup = append(s.onInterrupt, handlers...)
	}
}

func WithOnTermination(handlers ...SignalHandler) SignalerOption {
	return func(s *Signaler) {
		s.onTermination = append(s.onInterrupt, handlers...)
	}
}

func NewSignaler(options ...SignalerOption) *Signaler {

	s := &Signaler{}

	for _, option := range options {
		option(s)
	}

	return s

}

func (s *Signaler) WaitForSignal(errorHandler func(error)) {

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	s.wg.Add(1)

	go s.waitForSignal(sigChan, errorHandler)

	s.wg.Wait()

}

func (s *Signaler) waitForSignal(sigChan chan os.Signal, errorHandler func(error)) {

	var handlers []SignalHandler
	var once sync.Once

	release := func() {
		once.Do(func() {
			close(sigChan)
		})
	}

	for sig := range sigChan {

		handlers = nil

		switch sig {
		case syscall.SIGINT:
			handlers = s.onInterrupt
		case syscall.SIGTERM:
			handlers = s.onTermination
		case syscall.SIGHUP:
			handlers = s.onHangup
		default:
			continue
		}

		for _, handler := range handlers {
			err := handler(release)
			if err != nil {
				errorHandler(err)
			}
		}

	}

	s.wg.Done()

}
