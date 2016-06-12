package state

import (
	"errors"

	"github.com/micro/go-micro/registry"
)

type Watch struct {
	close  chan struct{}
	result chan *registry.Result
}

func (w *Watch) Next() (*registry.Result, error) {
	select {
	case <-w.close:
		return nil, errors.New("Watcher has been stopped")
	case r := <-w.result:
		return r, nil
	}
}

func (w *Watch) Stop() {
	select {
	case <-w.close:
		return
	default:
		close(w.close)
	}
}
