package gossip

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/facebookgo/freeport"
	"github.com/hashicorp/memberlist"
	"github.com/micro/go-micro/registry"

	"golang.org/x/net/context"
)

// Mode is a network mdoe
type mode int

// Nework modes
const (
	Local mode = iota
	LAN
	WAN
)

var (
	// DefaultKey is the default key, you should change this if using secure
	DefaultKey = []byte("DefaultGossipKey")
)

type contextModeKey struct{}

// NetworkMode sets the network mode
func NetworkMode(m mode) registry.Option {
	return func(o *registry.Options) {
		o.Context = context.WithValue(o.Context, contextModeKey{}, m)
	}
}

func getMemberlistConfig(options *registry.Options) *memberlist.Config {
	if m, ok := options.Context.Value(contextModeKey{}).(mode); ok {
		switch m {
		case Local:
			return memberlist.DefaultLocalConfig()
		case LAN:
			return memberlist.DefaultLANConfig()
		case WAN:
			return memberlist.DefaultWANConfig()
		}
	}
	return memberlist.DefaultWANConfig()
}

type contextSecretKey struct{}

// SecretKey sets the secret key for gossip
func SecretKey(k []byte) registry.Option {
	return func(o *registry.Options) {
		o.Context = context.WithValue(o.Context, contextSecretKey{}, k)
	}
}

func applySecretKey(options *registry.Options, config *memberlist.Config) error {
	if options.Secure {
		key := DefaultKey
		if k, ok := options.Context.Value(contextSecretKey{}).([]byte); ok {
			key = k
		}
		config.SecretKey = key
	}
	return nil
}

type contextAddressKey struct{}

// Address sets the bind address:port
func Address(address string) registry.Option {
	return func(o *registry.Options) {
		o.Context = context.WithValue(o.Context, contextAddressKey{}, address)
	}
}

func applyAddress(options *registry.Options, config *memberlist.Config) error {
	if addr, ok := options.Context.Value(contextAddressKey{}).(string); ok {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}

		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		if p == 0 {
			p, err = freeport.Get()
			if err != nil {
				return err
			}
		}

		config.BindAddr = host
		config.BindPort = p
	}
	return nil
}

type contextAdvertiseKey struct{}

// Advertise sets the address:port avertised
func Advertise(address string) registry.Option {
	return func(o *registry.Options) {
		o.Context = context.WithValue(o.Context, contextAdvertiseKey{}, address)
	}
}

func applyAdvertise(options *registry.Options, config *memberlist.Config) error {
	if addr, ok := options.Context.Value(contextAdvertiseKey{}).(string); ok {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}

		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		if p == 0 {
			p = config.BindPort
		}

		config.AdvertiseAddr = host
		config.AdvertisePort = p
	}
	return nil
}

type contextLoggerKey struct{}

// Logger sets the logger used for the memberlist instance
func Logger(l *log.Logger) registry.Option {
	return func(o *registry.Options) {
		o.Context = context.WithValue(o.Context, contextLoggerKey{}, l)
	}
}

func applyLogger(options *registry.Options, config *memberlist.Config) *log.Logger {
	if logger, ok := options.Context.Value(contextLoggerKey{}).(*log.Logger); ok {
		config.Logger = logger
		return logger
	}
	return log.New(os.Stderr, "", log.LstdFlags)
}
