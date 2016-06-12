package test

import (
	"os"
	"path"
	"path/filepath"
)

// Options contains the runnner options
type Options struct {
	Image            string
	Size             int
	Mount            map[string]string
	Env              map[string]string
	WorkingDirectory string
	NodeCommand      []string
	TestCommand      []string
	WriteToFile      bool
}

// DefaultOptions returns the default options
func DefaultOptions() *Options {
	goPath := os.Getenv("GOPATH")

	wd, _ := os.Getwd()
	rel, _ := filepath.Rel(goPath, wd)
	toSlash := filepath.ToSlash(rel)
	cwd := path.Join("/go", toSlash)

	return &Options{
		Image:            "golang",
		Size:             2,
		Mount:            map[string]string{goPath: "/go"},
		Env:              map[string]string{"GOPATH": "/go"},
		WorkingDirectory: cwd,
		NodeCommand:      []string{"sh", "-c", "go build -o=/node && /node"},
		TestCommand:      []string{"go", "test"},
	}
}

// Option is a test runner option
type Option func(*Options)

// Mount adds a mount
func Mount(host string, container string) Option {
	return func(o *Options) {
		if o.Mount == nil {
			o.Mount = make(map[string]string)
		}

		o.Mount[host] = container
	}
}

// Env sets an enviroment var
func Env(key string, value string) Option {
	return func(o *Options) {
		if o.Env == nil {
			o.Env = make(map[string]string)
		}

		o.Env[key] = value
	}
}

// Size sets the cluster size
func Size(size int) Option {
	return func(o *Options) {
		o.Size = size
	}
}

// Image sets the image to use
func Image(image string) Option {
	return func(o *Options) {
		o.Image = image
	}
}

// WorkingDirectory sets the working dir in the container
func WorkingDirectory(dir string) Option {
	return func(o *Options) {
		o.WorkingDirectory = dir
	}
}

// StartNodeCommand sets the command to run per cluster node
func StartNodeCommand(command ...string) Option {
	return func(o *Options) {
		o.NodeCommand = command
	}
}

// RunTestsCommand sets the command to run on the test comtainer
func RunTestsCommand(command ...string) Option {
	return func(o *Options) {
		o.TestCommand = command
	}
}

// WriteToFile ensures container logs are written to a file
func WriteToFile(enabled bool) Option {
	return func(o *Options) {
		o.WriteToFile = enabled
	}
}
