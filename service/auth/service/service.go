package service

import (
	"time"

	"github.com/ThatsMrTalbot/cluster/service/auth/proto"
	"github.com/ThatsMrTalbot/prototoken"
	"github.com/micro/go-micro/server"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Interface defines how users are authenticated
type Interface interface {
	Auth(username string, password string) (*proto.User, error)
}

// Auth defines a go micro service
type Auth struct {
	opts  *options
	iface Interface
}

// RegisterAuthHandler registers an auth handler on the server
func RegisterAuthHandler(s server.Server, iface Interface, opts ...Option) error {
	auth := &Auth{
		opts:  parse(opts...),
		iface: iface,
	}

	handler := s.NewHandler(auth)
	err := s.Handle(handler)
	return errors.Wrap(err, "Could not attach handler to server")
}

// Auth authenticates and returns a token
func (a *Auth) Auth(ctx context.Context, req *proto.AuthRequest, rsp *proto.Response) error {
	user, err := a.iface.Auth(req.Username, req.Password)
	if err != nil {
		return errors.Wrap(err, "Authenticator returned error")
	}

	rsp.Token, rsp.Refresh, err = a.generate(user)
	if err != nil {
		return err
	}

	return nil
}

// Refresh generates a new token when provided with a valid refresh token
func (a *Auth) Refresh(ctx context.Context, req *proto.RefreshRequest, rsp *proto.Response) error {
	var tok proto.Token
	_, err := prototoken.ValidateString(req.Token, a.opts.RefreshPublicKey, &tok)
	if err != nil {
		return errors.Wrap(err, "Unable to validate refresh token")
	}

	if tok.Type != proto.TokenType_Refresh {
		return errors.New("Provided token is not a refresh token")
	}

	if tok.Expiry < time.Now().UTC().Unix() {
		return errors.New("Provided token has expired")
	}

	rsp.Token, rsp.Refresh, err = a.generate(tok.User)
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) generate(user *proto.User) (string, string, error) {
	tokenExp := int64(0)
	if a.opts.TokenExpiry > 0 {
		tokenExp = time.Now().UTC().Add(a.opts.TokenExpiry).Unix()
	}

	token := &proto.Token{
		Type:   proto.TokenType_Auth,
		User:   user,
		Expiry: tokenExp,
	}

	refreshExp := int64(0)
	if a.opts.TokenExpiry > 0 {
		refreshExp = time.Now().UTC().Add(a.opts.RefreshExpiry).Unix()
	}

	refresh := &proto.Token{
		Type:   proto.TokenType_Refresh,
		User:   user,
		Expiry: refreshExp,
	}

	tokenString, err := prototoken.GenerateString(token, a.opts.TokenPrivateKey)
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to generate token")
	}

	refreshString, err := prototoken.GenerateString(refresh, a.opts.RefreshPrivateKey)
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to generate refresh token")
	}

	return tokenString, refreshString, nil
}
