package client

import (
	"time"

	"github.com/ThatsMrTalbot/cluster/service/auth/proto"
	"github.com/ThatsMrTalbot/prototoken"
	"github.com/pkg/errors"
)

// Expired checks if token has expired
func Expired(token string) (bool, error) {
	tok, err := prototoken.UnpackString(token)
	if err != nil {
		return false, errors.Wrap(err, "Could not unpack token")
	}

	var value proto.Token
	err = prototoken.ExtractMessage(tok, &value)
	if err != nil {
		return false, errors.Wrap(err, "Could not extract value")
	}

	return expired(&value), nil
}

func expired(token *proto.Token) bool {
	return token.Expiry != 0 && token.Expiry < time.Now().UTC().Unix()
}

// HasPermission returns true if the user has a permission
func HasPermission(token *proto.Token, permission string) bool {
	for _, perm := range token.User.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// UID returns the user id contained in the token
func UID(token *proto.Token) string {
	return token.User.UID
}
