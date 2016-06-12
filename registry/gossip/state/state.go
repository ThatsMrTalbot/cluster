//go:generate protoc --go_out=. *.proto

package state

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"

	"github.com/micro/go-micro/registry"
	"github.com/micro/protobuf/proto"
	"github.com/pborman/uuid"
)

//TODO
// State tests
// State hooks (to push changes)

// State is the state of the registry
type State struct {
	mu       sync.RWMutex
	services map[string][]*registry.Service
	index    *Index
	subs     map[string]chan *registry.Result
}

// NewState creates a new state
func NewState() *State {
	return &State{
		services: make(map[string][]*registry.Service),
		index:    &Index{},
		subs:     make(map[string]chan *registry.Result),
	}
}

// Watch creates a watcher
func (state *State) Watch() (registry.Watcher, error) {
	kill := make(chan struct{})
	result := make(chan *registry.Result, 10)
	id := uuid.NewUUID().String()

	watch := &Watch{
		result: result,
		close:  kill,
	}

	state.mu.Lock()
	if state.subs == nil {
		state.subs = make(map[string]chan *registry.Result)
	}
	state.subs[id] = result
	state.mu.Unlock()

	go func() {
		<-kill

		// unsub
		state.mu.Lock()
		delete(state.subs, id)
		state.mu.Unlock()

		// close
		close(result)
	}()

	return watch, nil
}

func (state *State) pub(res []*registry.Result) {
	for _, r := range res {
		for _, c := range state.subs {
			c <- r
		}
	}
}

// String implements strings
func (state *State) String() string {
	return "state"
}

// GetService gets a service by name
func (state *State) GetService(name string) ([]*registry.Service, error) {
	state.mu.RLock()
	defer state.mu.RUnlock()

	if state.services != nil && len(state.services[name]) != 0 {
		return state.services[name], nil
	}
	return nil, fmt.Errorf("Service %s not found", name)
}

// ListServices lists all services
func (state *State) ListServices() ([]*registry.Service, error) {
	state.mu.RLock()
	defer state.mu.RUnlock()

	s := make([]*registry.Service, 0, len(state.services))
	for _, services := range state.services {
		s = append(s, services...)
	}
	return s, nil
}

// Register a service
func (state *State) Register(s *registry.Service, ops ...registry.RegisterOption) error {
	_, err := state.RegisterAndReturnChange(s, ops...)
	return err
}

// RegisterAndReturnChange registers and returns a mergable change
func (state *State) RegisterAndReturnChange(s *registry.Service, ops ...registry.RegisterOption) ([]byte, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	options := registry.RegisterOptions{
		Context: context.TODO(),
	}

	for _, o := range ops {
		o(&options)
	}

	diff, change, err := state.index.Add(nil, s, options.TTL)
	if err != nil {
		return nil, err
	}

	state.pub(diff)

	services, err := state.index.ToMap()
	if err != nil {
		return nil, err
	}

	c, err := proto.Marshal(change)
	if err != nil {
		return nil, err
	}

	state.services = services

	return c, nil
}

// Deregister a service
func (state *State) Deregister(s *registry.Service) error {
	_, err := state.DeregisterAndReturnChange(s)

	return err
}

// DeregisterAndReturnChange deregisters and returns a mergable change
func (state *State) DeregisterAndReturnChange(s *registry.Service) ([]byte, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	diff, change, err := state.index.Remove(nil, s)
	if err != nil {
		return nil, err
	}

	state.pub(diff)

	services, err := state.index.ToMap()
	if err != nil {
		return nil, err
	}

	c, err := proto.Marshal(change)
	if err != nil {
		return nil, err
	}

	state.services = services

	return c, nil
}

// MergeRemote merges remote state
func (state *State) MergeRemote(byt []byte) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	var index Index
	if err := proto.Unmarshal(byt, &index); err != nil {
		return err
	}

	diff, err := state.index.Merge(nil, &index)
	if err != nil {
		return err
	}

	services, err := state.index.ToMap()
	if err != nil {
		return err
	}

	state.services = services

	state.pub(diff)

	return nil
}

// LocalState returns the local state
func (state *State) LocalState() ([]byte, error) {
	state.mu.RLock()
	defer state.mu.RUnlock()

	return proto.Marshal(state.index)
}
