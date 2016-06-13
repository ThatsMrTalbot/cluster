package state

import (
	"math/rand"
	"testing"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/pborman/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestState(t *testing.T) {
	Convey("Given a state instance", t, WithState(func(s *State) {

		Convey("When no services are registered", func() {
			Convey("Then there should not be any services available", func() {
				list, err := s.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 0)

				service, err := s.GetService("test")
				So(err, ShouldNotBeNil)
				So(service, ShouldHaveLength, 0)
			})
		})

		Convey("When a service is registered", WithService(s, func(service *registry.Service) {
			Convey("Then the service should be registered", func() {
				list, err := s.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 1)

				service, err := s.GetService("test")
				So(err, ShouldBeNil)
				So(service, ShouldHaveLength, 1)
				So(service[0].Nodes, ShouldHaveLength, 1)
			})
		}))

		Convey("When multiple services are registered", WithService(s, func(service *registry.Service) {
			WithService(s, nil)()

			Convey("Then the services should be registered", func() {
				list, err := s.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 1)

				service, err := s.GetService("test")
				So(err, ShouldBeNil)
				So(service, ShouldHaveLength, 1)
				So(service[0].Nodes, ShouldHaveLength, 2)
			})
		}))
	}))
}

func WithState(f func(*State)) func() {
	return func() {
		s := NewState(time.Second)
		f(s)
	}
}

func WithService(s *State, f func(*registry.Service)) func() {
	return func() {
		service := &registry.Service{
			Name:    "test",
			Version: "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:      uuid.NewUUID().String(),
					Address: "127.0.0.1",
					Port:    rand.Int(),
				},
			},
		}

		err := s.Register(service)
		So(err, ShouldBeNil)

		Reset(func() {
			s.Deregister(service)
		})

		if f != nil {
			f(service)
		}
	}
}
