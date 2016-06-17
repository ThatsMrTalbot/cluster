package client

import (
	"github.com/ThatsMrTalbot/cluster/service/auth/proto"
	"github.com/ThatsMrTalbot/prototoken"
	"github.com/micro/go-micro/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Client performs auth requests against the service
type Client struct {
	opts    *options
	service string
	client  client.Client
}

// NewClient creates an auth client
func NewClient(client client.Client, service string, opts ...Option) *Client {
	return &Client{
		client:  client,
		service: service,
		opts:    parse(opts...),
	}
}

// Auth performs an auth request against the service
func (c *Client) Auth(ctx context.Context, username string, password string) (*proto.Response, error) {
	rsp := new(proto.Response)
	req := client.NewRequest(c.service, "Auth.Auth", &proto.AuthRequest{
		Username: username,
		Password: password,
	})

	err := c.client.Call(ctx, req, rsp)
	if err != nil {
		return nil, errors.Wrap(err, "Could not authenticate against service")
	}

	return rsp, nil
}

// Refresh performs a refresh request against the service
func (c *Client) Refresh(ctx context.Context, token string) (*proto.Response, error) {
	rsp := new(proto.Response)
	req := client.NewRequest(c.service, "Auth.Refresh", &proto.RefreshRequest{
		Token: token,
	})

	err := c.client.Call(ctx, req, rsp)
	if err != nil {
		return nil, errors.Wrap(err, "Could not authenticate against service")
	}

	return rsp, nil
}

// Validate validates a token locally
// this will fail on refresh tokens
func (c *Client) Validate(token string) (*proto.Token, error) {
	data := new(proto.Token)
	_, err := prototoken.ValidateString(token, c.opts.PublicKey, data)
	if err != nil {
		return nil, errors.Wrap(err, "Could not validate token")
	}

	if data.Type == proto.TokenType_Refresh {
		return nil, errors.New("Refresh tokens cannot be validated")
	}

	if Expired(data) {
		return nil, errors.New("Token has expired")
	}

	return data, nil
}
