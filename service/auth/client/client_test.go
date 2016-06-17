package client

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/ThatsMrTalbot/cluster/service/auth/proto"
	"github.com/ThatsMrTalbot/prototoken"
	"github.com/micro/go-micro/client/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("Given a valid auth client", t, WithValidClient(func(client *Client) {
		Convey("When auth is called", func() {
			ctx := context.TODO()
			rsp, err := client.Auth(ctx, "username", "password")

			Convey("Then a valid response should be returned", func() {
				So(err, ShouldBeNil)
				So(rsp.Token, ShouldResemble, "token")
				So(rsp.Refresh, ShouldResemble, "refresh")
			})
		})

		Convey("When refresh is called", func() {
			ctx := context.TODO()
			rsp, err := client.Refresh(ctx, "token")

			Convey("Then a valid response should be returned", func() {
				So(err, ShouldBeNil)
				So(rsp.Token, ShouldResemble, "token")
				So(rsp.Refresh, ShouldResemble, "refresh")
			})
		})

		Convey("When a valid token is validated", func() {
			token := &proto.Token{
				Type: proto.TokenType_Auth,
				User: &proto.User{UID: "test", Permissions: []string{"a"}},
			}

			refresh, err := prototoken.GenerateString(token, prototoken.NewHMACPrivateKey([]byte("DefaultSecret")))
			So(err, ShouldBeNil)

			result, err := client.Validate(refresh)

			Convey("Then the token should be validated", func() {
				So(err, ShouldBeNil)
				So(UID(result), ShouldEqual, "test")
				So(HasPermission(result, "a"), ShouldBeTrue)
				So(HasPermission(result, "b"), ShouldBeFalse)
			})
		})

		Convey("When an expired token is validated", func() {
			token := &proto.Token{
				Type:   proto.TokenType_Auth,
				User:   &proto.User{UID: "test"},
				Expiry: 1,
			}

			refresh, err := prototoken.GenerateString(token, prototoken.NewHMACPrivateKey([]byte("DefaultSecret")))
			So(err, ShouldBeNil)

			_, err = client.Validate(refresh)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When a refresh token is validated", func() {
			token := &proto.Token{
				Type: proto.TokenType_Refresh,
				User: &proto.User{UID: "test"},
			}

			refresh, err := prototoken.GenerateString(token, prototoken.NewHMACPrivateKey([]byte("DefaultSecret")))
			So(err, ShouldBeNil)

			_, err = client.Validate(refresh)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When an invalid token is validated", func() {
			_, err := client.Validate("invalid")

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	}))

	Convey("Given an invalid auth client", t, WithInvalidClient(func(client *Client) {
		Convey("When auth is called", func() {
			ctx := context.TODO()
			_, err := client.Auth(ctx, "username", "password")

			Convey("Then a valid response should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When refresh is called", func() {
			ctx := context.TODO()
			_, err := client.Refresh(ctx, "token")

			Convey("Then a valid response should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	}))
}

func WithValidClient(f func(*Client)) func() {
	return func() {
		client := mock.NewClient(
			mock.Response("service", []mock.MockResponse{
				{
					Method: "Auth.Auth",
					Response: proto.Response{
						Token:   "token",
						Refresh: "refresh",
					},
				},
				{
					Method: "Auth.Refresh",
					Response: proto.Response{
						Token:   "token",
						Refresh: "refresh",
					},
				},
			}),
		)
		c := NewClient(client, "service")
		f(c)
	}
}

func WithInvalidClient(f func(*Client)) func() {
	return func() {
		client := mock.NewClient(
			mock.Response("service", []mock.MockResponse{
				{
					Method: "Auth.Auth",
					Error:  errors.New("Some error"),
				},
				{
					Method: "Auth.Refresh",
					Error:  errors.New("Some error"),
				},
			}),
		)
		c := NewClient(client, "service")
		f(c)
	}
}
