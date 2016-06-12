package main

import (
	"io/ioutil"
	"time"

	"github.com/ThatsMrTalbot/cluster/test"
	"github.com/hashicorp/memberlist"
)

var quiet bool

func main() {
	config := memberlist.DefaultLANConfig()
	config.AdvertiseAddr = test.Address()

	if quiet {
		config.LogOutput = ioutil.Discard
	}

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	list.Join([]string{"cluster_1"})

	time.Sleep(time.Hour)
}
