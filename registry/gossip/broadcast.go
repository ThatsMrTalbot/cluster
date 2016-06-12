package gossip

import "github.com/hashicorp/memberlist"

type broadcast []byte

func (b broadcast) Invalidates(memberlist.Broadcast) bool {
	return false
}

func (b broadcast) Message() []byte {
	return []byte(b)
}

func (b broadcast) Finished() {
	// Finished
}
