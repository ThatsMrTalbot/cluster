package gossip

import (
	"log"
	"os"

	"github.com/ThatsMrTalbot/cluster/registry/gossip/state"
	"github.com/hashicorp/memberlist"
	"github.com/micro/go-micro/registry"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type gossip struct {
	*state.State
	*memberlist.TransmitLimitedQueue

	m *memberlist.Memberlist
	l *log.Logger
}

func (g *gossip) NodeMeta(int) []byte {
	return nil
}

func (g *gossip) NotifyMsg(buf []byte) {
	err := g.State.MergeRemote(buf)
	if err != nil {
		g.l.Printf("[ERROR] Error merging broadcast: %s", err)
	}
}

func (g *gossip) LocalState(join bool) []byte {
	byt, err := g.State.LocalState()
	if err != nil {
		g.l.Printf("[ERROR] Error getting local state: %s", err)
	}
	return byt
}

func (g *gossip) MergeRemoteState(buf []byte, join bool) {
	err := g.State.MergeRemote(buf)
	if err != nil {
		g.l.Printf("[ERROR] Error merging remote state: %s", err)
	}
}

func (g *gossip) Deregister(s *registry.Service) error {
	change, err := g.DeregisterAndReturnChange(s)
	if err != nil {
		return errors.Wrap(err, "Error deregistering service")
	}

	// Broadcast change
	g.QueueBroadcast(broadcast(change))
	return nil
}

func (g *gossip) Register(s *registry.Service, ops ...registry.RegisterOption) error {
	change, err := g.RegisterAndReturnChange(s, ops...)
	if err != nil {
		return errors.Wrap(err, "Error registering service")
	}

	// Broadcast change
	g.QueueBroadcast(broadcast(change))
	return nil
}

// NewRegistry creates a new registry
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := &registry.Options{
		Context: context.TODO(),
	}

	for _, o := range opts {
		o(options)
	}

	hostname, _ := os.Hostname()
	config := getMemberlistConfig(options)
	config.Name = hostname + "-" + uuid.NewUUID().String()

	log := applyLogger(options, config)

	if err := applySecretKey(options, config); err != nil {
		log.Fatalf("Error creating memberlist: %s", err)
	}

	if err := applyAddress(options, config); err != nil {
		log.Fatalf("Error creating memberlist: %s", err)
	}

	if err := applyAdvertise(options, config); err != nil {
		log.Fatalf("Error creating memberlist: %s", err)
	}

	g := new(gossip)

	config.Delegate = g

	m, err := memberlist.Create(config)
	if err != nil {
		log.Fatalf("Error creating memberlist: %s", err)
	}

	g.m = m
	g.l = log
	g.State = state.NewState(ExpiryTick)
	g.TransmitLimitedQueue = &memberlist.TransmitLimitedQueue{
		NumNodes:       m.NumMembers,
		RetransmitMult: 3,
	}

	if len(options.Addrs) != 0 {
		_, err := m.Join(options.Addrs)
		if err != nil {
			log.Fatalf("Error creating memberlist: %s", err)
		}
	}

	return g
}
