package service

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/ThatsMrTalbot/cluster/service/auth/proto"
	"github.com/micro/go-micro/server/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("Given a server", t, func() {
		server := mock.NewServer()

		Convey("When auth is registered", func() {
			err := RegisterAuthHandler(server, DummyInterface{})
			So(err, ShouldBeNil)

			Convey("Then the handler should be registered", func() {
				So(server.Handlers, ShouldHaveLength, 1)
			})
		})
	})
}

func TestAuth(t *testing.T) {
	Convey("Given a service", t, func() {
		auth := &Auth{
			opts:  parse(),
			iface: DummyInterface{},
		}

		Convey("When auth is called with a valid username and password", func() {
			ctx := context.TODO()
			req := &proto.AuthRequest{
				Username: "username",
				Password: "password",
			}
			rsp := &proto.Response{}
			err := auth.Auth(ctx, req, rsp)

			Convey("Then the result should contain a token", func() {
				So(err, ShouldBeNil)
				So(rsp.Token, ShouldNotBeEmpty)
				So(rsp.Refresh, ShouldNotBeEmpty)
			})
		})

		Convey("When auth is called with an invalid username and password", func() {
			ctx := context.TODO()
			req := &proto.AuthRequest{
				Username: "invalid",
				Password: "invalid",
			}
			rsp := &proto.Response{}
			err := auth.Auth(ctx, req, rsp)

			Convey("Then the result should contain a token", func() {
				So(err, ShouldNotBeNil)
				So(rsp.Token, ShouldBeEmpty)
				So(rsp.Refresh, ShouldBeEmpty)
			})
		})

		Convey("When refresh is called with a valid refresh token", func() {
			_, token, err := auth.generate(&proto.User{})
			So(err, ShouldBeNil)

			ctx := context.TODO()
			req := &proto.RefreshRequest{
				Token: token,
			}
			rsp := &proto.Response{}
			err = auth.Refresh(ctx, req, rsp)

			Convey("Then the result should contain a token", func() {
				So(err, ShouldBeNil)
				So(rsp.Token, ShouldNotBeEmpty)
				So(rsp.Refresh, ShouldNotBeEmpty)
			})
		})

		Convey("When refresh is called with an invalid refresh token", func() {
			ctx := context.TODO()
			req := &proto.RefreshRequest{
				Token: "invalid",
			}
			rsp := &proto.Response{}
			err := auth.Refresh(ctx, req, rsp)

			Convey("Then the result should contain a token", func() {
				So(err, ShouldNotBeNil)
				So(rsp.Token, ShouldBeEmpty)
				So(rsp.Refresh, ShouldBeEmpty)
			})
		})
	})
}

type DummyInterface struct{}

func (DummyInterface) Auth(u string, p string) (*proto.User, error) {
	if u == "username" && p == "password" {
		return &proto.User{
			UID:         "123",
			Permissions: []string{"a", "b", "c"},
		}, nil
	}
	return nil, errors.New("Incorrect username or password")
}
