package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro/go-micro/registry"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWatch(t *testing.T) {
	Convey("Given a gossip registry", t, WithState(func(s *State) {
		Convey("When a local watcher is initiated", WithWatcher(s, func(w registry.Watcher) {
			Convey("Then changes should be published to the watcher", func() {
				WithService(s, func(service *registry.Service) {
					expected := &registry.Result{
						Action:  "create",
						Service: service,
					}

					So(w, ShouldHaveNext, expected)
				})()

			})
		}))
	}))
}

func WithWatcher(s *State, f func(registry.Watcher)) func() {
	return func() {
		watcher, err := s.Watch()
		So(err, ShouldBeNil)

		Reset(func() {
			watcher.Stop()
		})

		f(watcher)
	}
}

func ShouldHaveNext(actual interface{}, expected ...interface{}) string {
	watcher, ok := actual.(registry.Watcher)
	if !ok {
		return fmt.Sprintf("Expected registry.Watcher got %t", actual)
	}

	result := make(chan *registry.Result)
	err := error(nil)

	go func() {
		r, e := watcher.Next()
		err = e
		result <- r
	}()

	timeout := time.After(time.Second * 5)
	select {
	case r := <-result:
		if err != nil {
			return fmt.Sprintf("Next returned error: %s", err)
		}
		if len(expected) > 0 {
			return ShouldResemble(r, expected...)
		}
		return ""
	case <-timeout:
		return "No next event triggered"
	}
}
