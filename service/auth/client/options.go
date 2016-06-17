package client

import (
	"crypto/rsa"

	"github.com/ThatsMrTalbot/prototoken"
)

var (
	// DefaultPublicKey is the default public key for token and refresh token
	DefaultPublicKey = prototoken.NewHMACPublicKey([]byte("DefaultSecret"))
)

type options struct {
	PublicKey prototoken.PublicKey
}

func parse(opts ...Option) *options {
	options := &options{
		PublicKey: DefaultPublicKey,
	}

	for _, o := range opts {
		o(options)
	}

	return options
}

// Option is a options for client or service
type Option func(*options)

// PublicKeyHMAC sets the public key for both tokens and refresh tokens
func PublicKeyHMAC(secret []byte) Option {
	return func(o *options) {
		o.PublicKey = prototoken.NewHMACPublicKey(secret)
	}
}

// PublicKeyRSA sets the public key for both tokens and refresh tokens
func PublicKeyRSA(key *rsa.PublicKey) Option {
	return func(o *options) {
		o.PublicKey = prototoken.NewRSAPublicKey(key)
	}
}
