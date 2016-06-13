package gossip

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/facebookgo/freeport"
	"github.com/micro/go-micro/registry"
	"github.com/pborman/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegistry(t *testing.T) {
	Convey("Given a gossip registy", t, WithRegistry(nil, func(r1 registry.Registry, addr string, port int) {
		r1Address := fmt.Sprintf("%s:%d", addr, port)

		Convey("When no services are registered", func() {
			Convey("Then there should not be any services available", func() {
				list, err := r1.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 0)

				service, err := r1.GetService("test")
				So(err, ShouldNotBeNil)
				So(service, ShouldHaveLength, 0)
			})
		})

		Convey("When a service is registered", WithService(r1, "test", addr, port, func(service *registry.Service) {
			Convey("Then the service should be registered", func() {
				list, err := r1.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 1)

				service, err := r1.GetService("test")
				So(err, ShouldBeNil)
				So(service, ShouldHaveLength, 1)
				So(service[0].Nodes, ShouldHaveLength, 1)
			})
		}))

		Convey("When multiple services are registered", WithService(r1, "test", addr, port, func(service *registry.Service) {
			WithService(r1, "test", addr, port+1, nil)()

			Convey("Then the services should be registered", func() {
				list, err := r1.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 1)

				service, err := r1.GetService("test")
				So(err, ShouldBeNil)
				So(service, ShouldHaveLength, 1)
				So(service[0].Nodes, ShouldHaveLength, 2)
			})
		}))

		Convey("When another registry joins", WithRegistry([]string{r1Address}, func(r2 registry.Registry, _ string, _ int) {
			Convey("Then services should propigate", WithService(r1, "test", addr, port, func(*registry.Service) {
				time.Sleep(time.Second * 1)

				list, err := r2.ListServices()
				So(err, ShouldBeNil)
				So(list, ShouldHaveLength, 1)

				service, err := r2.GetService("test")
				So(err, ShouldBeNil)
				So(service, ShouldHaveLength, 1)
				So(service[0].Nodes, ShouldHaveLength, 1)
			}))
		}))
	}))
}

func WithRegistry(addrs []string, f func(registry.Registry, string, int)) func() {
	return func() {
		portInt, err := freeport.Get()
		So(err, ShouldBeNil)

		port := strconv.Itoa(portInt)

		reg := NewRegistry(
			Address("127.0.0.1:"+port),
			Advertise("127.0.0.1:"+port),
			Logger(log.New(ioutil.Discard, "", log.LstdFlags)),
			NetworkMode(Local),
			SecretKey([]byte("SixteenBytTstKey")),
			registry.Addrs(addrs...),
			registry.Secure(true),
		)

		Reset(func() {
			m := reg.(*gossip).m
			m.Leave(time.Second * 10)
			m.Shutdown()
		})

		if f != nil {
			f(reg, "127.0.0.1", portInt)
		}
	}
}

func WithService(reg registry.Registry, name string, addr string, port int, f func(*registry.Service)) func() {
	return func() {
		service := &registry.Service{
			Name:     name,
			Metadata: make(map[string]string),
			Version:  "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:       uuid.NewUUID().String(),
					Address:  addr,
					Port:     port,
					Metadata: make(map[string]string),
				},
			},
		}

		err := reg.Register(service)
		So(err, ShouldBeNil)

		Reset(func() {
			reg.Deregister(service)
		})

		if f != nil {
			f(service)
		}
	}
}
