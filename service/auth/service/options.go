package service

import (
	"crypto/rsa"
	"time"

	"github.com/ThatsMrTalbot/prototoken"
)

var (
	// DefaultPublicKey is the default public key for token and refresh token
	DefaultPublicKey = prototoken.NewHMACPublicKey([]byte("DefaultSecret"))

	// DefaultPrivateKey is the default private key for token and refresh token
	DefaultPrivateKey = prototoken.NewHMACPrivateKey([]byte("DefaultSecret"))

	// DefaultTokenExpiry is the default expiry for tokens
	DefaultTokenExpiry = time.Minute * 30

	// DefaultRefreshExpiry is the default expiry for refresh tokens
	DefaultRefreshExpiry = time.Hour * 2
)

type options struct {
	TokenPrivateKey prototoken.PrivateKey
	TokenPublicKey  prototoken.PublicKey
	TokenExpiry     time.Duration

	RefreshPrivateKey prototoken.PrivateKey
	RefreshPublicKey  prototoken.PublicKey
	RefreshExpiry     time.Duration
}

func parse(opts ...Option) *options {
	options := &options{
		TokenExpiry:       DefaultTokenExpiry,
		TokenPublicKey:    DefaultPublicKey,
		TokenPrivateKey:   DefaultPrivateKey,
		RefreshPublicKey:  DefaultPublicKey,
		RefreshPrivateKey: DefaultPrivateKey,
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
		o.TokenPublicKey = prototoken.NewHMACPublicKey(secret)
		o.RefreshPublicKey = o.TokenPublicKey
	}
}

// PublicKeyRSA sets the public key for both tokens and refresh tokens
func PublicKeyRSA(key *rsa.PublicKey) Option {
	return func(o *options) {
		o.TokenPublicKey = prototoken.NewRSAPublicKey(key)
		o.RefreshPublicKey = o.TokenPublicKey
	}
}

// PrivateKeyRSA sets the private key for both tokens and refresh tokens
func PrivateKeyRSA(key *rsa.PrivateKey) Option {
	return func(o *options) {
		o.TokenPrivateKey = prototoken.NewRSAPrivateKey(key)
		o.RefreshPrivateKey = o.TokenPrivateKey
	}
}

// PrivateKeyHMAC sets the private key for both tokens and refresh tokens
func PrivateKeyHMAC(secret []byte) Option {
	return func(o *options) {
		o.TokenPrivateKey = prototoken.NewHMACPrivateKey(secret)
		o.RefreshPrivateKey = o.TokenPrivateKey
	}
}

// TokenPublicKeyHMAC sets the public key to a hmac secret
func TokenPublicKeyHMAC(secret []byte) Option {
	return func(o *options) {
		o.TokenPublicKey = prototoken.NewHMACPublicKey(secret)
	}
}

// TokenPublicKeyRSA sets the public key to a rsa public key
func TokenPublicKeyRSA(key *rsa.PublicKey) Option {
	return func(o *options) {
		o.TokenPublicKey = prototoken.NewRSAPublicKey(key)
	}
}

// TokenPrivateKeyHMAC sets the private key to a hmac secret
func TokenPrivateKeyHMAC(secret []byte) Option {
	return func(o *options) {
		o.TokenPrivateKey = prototoken.NewHMACPrivateKey(secret)
	}
}

// TokenPrivateKeyRSA sets the private key to a rsa private key
func TokenPrivateKeyRSA(key *rsa.PrivateKey) Option {
	return func(o *options) {
		o.TokenPrivateKey = prototoken.NewRSAPrivateKey(key)
	}
}

// TokenExpiry sets the token expiry
func TokenExpiry(d time.Duration) Option {
	return func(o *options) {
		o.TokenExpiry = d
	}
}

// RefreshTokenPublicKeyHMAC sets the public key to a hmac secret
func RefreshTokenPublicKeyHMAC(secret []byte) Option {
	return func(o *options) {
		o.TokenPublicKey = prototoken.NewHMACPublicKey(secret)
	}
}

// RefreshTokenPublicKeyRSA sets the public key to a rsa public key
func RefreshTokenPublicKeyRSA(key *rsa.PublicKey) Option {
	return func(o *options) {
		o.RefreshPublicKey = prototoken.NewRSAPublicKey(key)
	}
}

// RefreshTokenPrivateKeyHMAC sets the private key to a hmac secret
func RefreshTokenPrivateKeyHMAC(secret []byte) Option {
	return func(o *options) {
		o.TokenPrivateKey = prototoken.NewHMACPrivateKey(secret)
	}
}

// RefreshTokenPrivateKeyRSA sets the private key to a rsa private key
func RefreshTokenPrivateKeyRSA(key *rsa.PrivateKey) Option {
	return func(o *options) {
		o.RefreshPrivateKey = prototoken.NewRSAPrivateKey(key)
	}
}

func RefreshTokenExpiry(d time.Duration) Option {
	return func(o *options) {
		o.RefreshExpiry = d
	}
}
