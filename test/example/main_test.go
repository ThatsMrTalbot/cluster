package main

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ThatsMrTalbot/cluster/test"
	"github.com/facebookgo/freeport"
	"github.com/hashicorp/memberlist"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMemberList(t *testing.T) {
	Convey("Given a memberlist", t, WithMemberList(func(list *memberlist.Memberlist) {

		Convey("When members are counted", func() {
			alive := list.NumMembers()

			Convey("Then there should be at least two", func() {
				So(alive, ShouldBeGreaterThanOrEqualTo, 2)
				Printf("\nThere are %d members!\n", alive)
			})
		})
	}))
}

func WithMemberList(runner func(*memberlist.Memberlist)) func() {
	return func() {
		config := memberlist.DefaultLocalConfig()
		port, _ := freeport.Get()

		config.Name = "TestRunner"
		config.BindPort = port
		config.AdvertiseAddr = test.Address()
		config.AdvertisePort = port
		config.LogOutput = ioutil.Discard

		list, err := memberlist.Create(config)
		So(err, ShouldBeNil)

		list.Join([]string{"127.0.0.1:7946"})
		time.Sleep(time.Second * 10)

		runner(list)

		list.Leave(time.Second * 10)
		list.Shutdown()
	}
}

func TestMain(m *testing.M) {
	docker := flag.Bool("docker", false, "Flag to run tests against a docker cluster")
	write := flag.Bool("write", false, "Flag write cluster logs to a file")
	flag.Parse()

	if *docker {
		runner := test.NewRunner(
			test.Size(10),
			test.WriteToFile(*write),
		)
		runner.StartAndTest()
	}

	quiet = true

	go main()

	time.Sleep(time.Second * 10)

	os.Exit(m.Run())
}
