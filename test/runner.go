package test

import (
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/network"
)

// Runner is a cluster test runner
type Runner struct {
	ops    *Options
	nodes  []*Node
	runner *Node
	nid    string
}

// NewRunner creates a new runner
func NewRunner(options ...Option) *Runner {
	ops := DefaultOptions()
	for _, o := range options {
		o(ops)
	}

	return NewRunnerFromOptions(ops)
}

// NewRunnerFromOptions creates a new runner from options struct
func NewRunnerFromOptions(options *Options) *Runner {
	nodes := []*Node{}
	for i := 0; i < options.Size; i++ {
		n := NewClusterNodeFromOptions(options)
		nodes = append(nodes, n)
	}

	n := NewTestNodeFromOptions(options)

	return &Runner{
		ops:    options,
		nodes:  nodes,
		runner: n,
	}
}

func (r *Runner) startNetwork() error {
	if r.nid != "" {
		return nil
	}

	ctx := context.TODO()
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	n, err := client.NetworkCreate(ctx, "test", types.NetworkCreate{
		CheckDuplicate: true,
	})
	if err != nil {
		return err
	}

	r.nid = n.ID

	return nil
}

func (r *Runner) killNetwork() error {
	ctx := context.TODO()
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	err = client.NetworkRemove(ctx, r.nid)
	r.nid = ""

	return err
}

func (r *Runner) joinNetwork(client *client.Client, cid string) error {
	if r.nid == "" {
		return nil
	}

	ctx := context.TODO()
	return client.NetworkConnect(ctx, r.nid, cid, &network.EndpointSettings{})
}

// Start cluster
func (r *Runner) Start() error {
	err := r.startNetwork()
	if err != nil {
		return err
	}

	for i, node := range r.nodes {
		err := node.Start(r.joinNetwork)
		if err != nil {
			return err
		}

		if i == 0 {
			// Give the first node a head start
			time.Sleep(time.Second * 5)
		}
	}
	return nil
}

// Kill cluster
func (r *Runner) Kill() error {
	for _, node := range r.nodes {
		err := node.Kill()
		if err != nil {
			return err
		}
	}

	err := r.killNetwork()
	if err != nil {
		return err
	}

	return nil
}

// Test runs tests and returns exit code
func (r *Runner) Test() (int, error) {
	err := r.runner.Start(r.joinNetwork)
	if err != nil {
		return 0, err
	}
	defer r.runner.Kill()

	go r.runner.Output()

	return r.runner.Wait()
}

// StartAndTest starts and runs test, exiting on completion
func (r *Runner) StartAndTest() {
	logrus.Info("Starting cluster")

	err := r.Start()
	if err != nil {
		logrus.WithError(err).Fatal("Could not start cluster")
	}

	logrus.Info("Waiting for cluster")

	time.Sleep(20 * time.Second)

	logrus.Info("Running tests")

	code, err := r.Test()
	if err != nil {
		logrus.WithError(err).Error("Could not run tests")
	}

	logrus.Info("Killing cluster")

	r.Kill()

	logrus.Info("Finished")

	os.Exit(code)
}
