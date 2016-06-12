package gossip

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro/go-micro/registry"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWatch(t *testing.T) {
	Convey("Given a gossip registry", t, WithRegistry(nil, func(r1 registry.Registry, addr string, port int) {
		r1Address := fmt.Sprintf("%s:%d", addr, port)

		Convey("When a local watcher is initiated", WithWatcher(r1, func(w registry.Watcher) {
			Convey("Then changes should be published to the watcher", func() {
				WithService(r1, "test", addr, port, func(s *registry.Service) {
					expected := &registry.Result{
						Action:  "create",
						Service: s,
					}

					So(w, ShouldHaveNext, expected)
				})()

			})
		}))

		Convey("When a watcher is initiated from a joined node", WithRegistry([]string{r1Address}, func(r2 registry.Registry, _ string, _ int) {
			WithWatcher(r2, func(w registry.Watcher) {
				Convey("Then changes should be published to the watcher", WithService(r1, "test", addr, port, func(s *registry.Service) {
					expected := &registry.Result{
						Action:  "create",
						Service: s,
					}

					So(w, ShouldHaveNext, expected)
				}))
			})()
		}))
	}))
}

func WithWatcher(r registry.Registry, f func(registry.Watcher)) func() {
	return func() {
		watcher, err := r.Watch()
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
