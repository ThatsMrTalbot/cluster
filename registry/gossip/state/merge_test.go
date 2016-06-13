package state

import (
	"testing"
	"time"

	"github.com/micro/go-micro/registry"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMerge(t *testing.T) {
	Convey("Given an index", t, func() {
		i := NewIndex()

		Convey("When a service is added", func() {
			service := &registry.Service{
				Name:     "test",
				Metadata: map[string]string{"a": "b"},
				Nodes: []*registry.Node{
					{
						Id:       "node1",
						Address:  "127.0.0.1",
						Port:     123,
						Metadata: make(map[string]string),
					},
				},
			}

			diff, _, err := i.Add(nil, service, 0)
			So(err, ShouldBeNil)

			Convey("Then the diff should contain a created event", func() {
				So(diff, ShouldHaveLength, 1)
				So(diff[0], ShouldResemble, &registry.Result{
					Action:  "create",
					Service: service,
				})
			})

			Convey("Then the service should exist", func() {
				m, err := i.ToMap()
				So(err, ShouldBeNil)
				So(m, ShouldContainKey, "test")
				So(m["test"], ShouldHaveLength, 1)
				So(m["test"][0], ShouldResemble, service)
			})
		})

		Convey("When a service is removed", func() {
			service := &registry.Service{
				Name:     "test",
				Metadata: map[string]string{"a": "b"},
				Nodes: []*registry.Node{
					{
						Id:      "node1",
						Address: "127.0.0.1",
						Port:    123,
					},
				},
			}

			_, _, err := i.Add(nil, service, 0)
			So(err, ShouldBeNil)

			diff, _, err := i.Remove(nil, service)
			So(err, ShouldBeNil)

			Convey("Then the diff should contain a removed event", func() {
				So(diff, ShouldHaveLength, 1)
				So(diff[0], ShouldResemble, &registry.Result{
					Action:  "delete",
					Service: service,
				})
			})

			Convey("Then the service should not exist", func() {
				m, err := i.ToMap()
				So(err, ShouldBeNil)
				So(m, ShouldNotContainKey, "test")
			})
		})

		Convey("When a service is changed", func() {
			service := &registry.Service{
				Name:     "test",
				Metadata: map[string]string{"a": "b"},
				Nodes: []*registry.Node{
					{
						Id:       "node1",
						Address:  "127.0.0.1",
						Port:     123,
						Metadata: make(map[string]string),
					},
				},
			}

			_, _, err := i.Add(nil, service, 0)
			So(err, ShouldBeNil)

			service.Metadata = make(map[string]string)

			diff, _, err := i.Add(nil, service, 0)
			So(err, ShouldBeNil)

			Convey("Then the diff should contain an update event", func() {
				So(diff, ShouldHaveLength, 1)
				So(diff[0], ShouldResemble, &registry.Result{
					Action:  "update",
					Service: service,
				})
			})

			Convey("Then the service should exist", func() {
				m, err := i.ToMap()
				So(err, ShouldBeNil)
				So(m, ShouldContainKey, "test")
				So(m["test"], ShouldHaveLength, 1)
				So(m["test"][0], ShouldResemble, service)
			})
		})

		Convey("When index is cleaned", func() {
			service := &registry.Service{
				Name:     "test",
				Metadata: map[string]string{"a": "b"},
				Nodes: []*registry.Node{
					{
						Id:      "node1",
						Address: "127.0.0.1",
						Port:    123,
					},
				},
			}

			_, _, err := i.Add(nil, service, -time.Second)
			So(err, ShouldBeNil)

			diff, err := i.Clean()

			Convey("Then the diff should contain an update event", func() {
				service.Nodes = []*registry.Node{}

				So(diff, ShouldHaveLength, 1)
				So(diff[0], ShouldResemble, &registry.Result{
					Action:  "delete",
					Service: service,
				})
			})

			Convey("Then the service should exist", func() {
				m, err := i.ToMap()
				So(err, ShouldBeNil)
				So(m, ShouldNotContainKey, "test")
			})
		})
	})
}
