package client

import (
	"time"

	"github.com/ThatsMrTalbot/cluster/service/auth/proto"
)

// Expired checks if token has expired
func Expired(token *proto.Token) bool {
	return token.Expiry != 0 && token.Expiry < time.Now().Unix()
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
