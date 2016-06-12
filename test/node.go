package test

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
	"golang.org/x/net/context"
)

// StartHook is run after create and before start
type StartHook func(client *client.Client, id string) error

// Node is a docker node
type Node struct {
	id   string
	name string
	test bool
	ops  *Options
}

// NewClusterNode creates a new node
func NewClusterNode(options ...Option) *Node {
	ops := DefaultOptions()
	for _, o := range options {
		o(ops)
	}

	return NewClusterNodeFromOptions(ops)
}

// NewTestNode creates a test runner node
func NewTestNode(options ...Option) *Node {
	ops := DefaultOptions()
	for _, o := range options {
		o(ops)
	}

	return NewTestNodeFromOptions(ops)
}

// NewClusterNodeFromOptions creates a new node
func NewClusterNodeFromOptions(options *Options) *Node {
	return &Node{
		ops: options,
	}
}

// NewTestNodeFromOptions creates a test runner node
func NewTestNodeFromOptions(options *Options) *Node {
	return &Node{
		test: true,
		ops:  options,
	}
}

func (n *Node) command() strslice.StrSlice {
	cmd := strslice.StrSlice{}
	if n.test {
		cmd = append(cmd, n.ops.TestCommand...)
	} else {
		cmd = append(cmd, n.ops.NodeCommand...)
	}
	return cmd
}

func (n *Node) init(client *client.Client) error {
	ctx := context.TODO()

	config := &container.Config{
		Image:      n.ops.Image,
		WorkingDir: n.ops.WorkingDirectory,
		Entrypoint: n.command(),
		Volumes:    make(map[string]struct{}),
	}

	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}

	networkConfig := &network.NetworkingConfig{}

	for host, container := range n.ops.Mount {
		config.Volumes[container] = struct{}{}
		hostConfig.Binds = append(hostConfig.Binds, mountString(host, container))
	}

	n.name = createID()

	c, err := client.ContainerCreate(ctx, config, hostConfig, networkConfig, n.name)
	n.id = c.ID

	return err
}

// Start node
func (n *Node) Start(hooks ...StartHook) error {
	ctx := context.TODO()

	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	err = n.init(client)
	if err != nil {
		return err
	}

	for _, hook := range hooks {
		hook(client, n.id)
	}

	if n.ops.WriteToFile {
		defer func() {
			go n.WriteToFile()
		}()
	}

	return client.ContainerStart(ctx, n.id, types.ContainerStartOptions{})
}

// ID gets the container id
func (n *Node) ID() string {
	return n.id
}

// Kill node
func (n *Node) Kill() error {
	ctx := context.TODO()

	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	return client.ContainerRemove(ctx, n.id, types.ContainerRemoveOptions{Force: true})
}

// Wait for node to exit
func (n *Node) Wait() (int, error) {
	ctx := context.TODO()

	client, err := client.NewEnvClient()
	if err != nil {
		return 0, err
	}

	return client.ContainerWait(ctx, n.id)
}

// Logs gets a log reader
func (n *Node) Logs() (io.ReadCloser, error) {
	ctx := context.TODO()

	client, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return client.ContainerLogs(ctx, n.id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "all",
		Follow:     true,
	})
}

// Output logs to stdout
func (n *Node) Output() error {
	reader, err := n.Logs()
	if err != nil {
		return err
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	return nil
}

// WriteToFile writes logs to file
func (n *Node) WriteToFile() error {
	reader, err := n.Logs()
	if err != nil {
		return err
	}
	defer reader.Close()

	fname := fmt.Sprintf("test_%s.log", n.name)
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, reader)

	return nil
}
