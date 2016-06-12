package state

import (
	"time"

	"golang.org/x/net/context"

	"github.com/micro/go-micro/registry"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// NewIndex creates a new index
func NewIndex() *Index {
	return &Index{
		Services: make(map[string]*Services),
	}
}

// Add merges a new service
func (i *Index) Add(ctx context.Context, s *registry.Service, timeout time.Duration) ([]*registry.Result, *Index, error) {
	mod := time.Now().UnixNano()

	ttl := int64(0)
	if timeout != 0 {
		ttl = time.Now().Add(timeout).UnixNano()
	}

	raw, err := msgpack.Marshal(s)
	if err != nil {
		return nil, nil, err
	}

	nodes := make(map[string]*Node)
	for _, node := range s.Nodes {
		nodes[node.Id] = &Node{
			Enabled: true,
			Mod:     mod,
			Expiry:  ttl,
		}
	}

	merge := &Index{
		map[string]*Services{
			s.Name: {
				Services: map[string]*Service{
					s.Version: {
						Nodes: nodes,
						Mod:   mod,
						Raw:   raw,
					},
				},
			},
		},
	}

	ctx = withService(ctx, s)
	diff, err := i.Merge(ctx, merge)
	return diff, merge, err
}

// Remove merges a removed service
func (i *Index) Remove(ctx context.Context, s *registry.Service) ([]*registry.Result, *Index, error) {
	mod := time.Now().UnixNano()

	nodes := make(map[string]*Node)
	for _, node := range s.Nodes {
		nodes[node.Id] = &Node{
			Enabled: false,
			Mod:     mod,
		}
	}

	s.Nodes = []*registry.Node{}

	raw, err := msgpack.Marshal(s)
	if err != nil {
		return nil, nil, err
	}

	merge := &Index{
		map[string]*Services{
			s.Name: {
				Services: map[string]*Service{
					s.Version: {
						Nodes: nodes,
						Mod:   mod,
						Raw:   raw,
					},
				},
			},
		},
	}

	ctx = withService(ctx, s)
	diff, err := i.Merge(ctx, merge)
	return diff, merge, err
}

// Merge merges one index into another
func (i *Index) Merge(ctx context.Context, merge *Index) ([]*registry.Result, error) {
	// Diff
	diff := []*registry.Result{}

	if i.Services == nil {
		i.Services = make(map[string]*Services)
	}

	for name, services := range merge.Services {
		for version, s1 := range services.Services {
			s2 := i.GetService(name, version)

			// If we dont have the sevice, so create it
			if s2 == nil {
				v, ok := i.Services[name]
				if !ok {
					v = &Services{Services: make(map[string]*Service)}
					i.Services[name] = v
				}
				s2 = &Service{Nodes: make(map[string]*Node)}
				v.Services[version] = s2
			}

			// Merge
			if err := s2.Merge(ctx, name, version, s1, &diff); err != nil {
				return nil, err
			}
		}
	}
	return diff, nil
}

// ToMap generates a lookup map for registry access
func (i *Index) ToMap() (map[string][]*registry.Service, error) {
	m := make(map[string][]*registry.Service)

	for name, services := range i.Services {
		slice := make([]*registry.Service, 0, len(services.Services))

		for _, service := range services.Services {
			var s *registry.Service

			if err := msgpack.Unmarshal(service.Raw, &s); err != nil {
				return nil, err
			}

			if len(s.Nodes) == 0 {
				continue
			}

			slice = append(slice, s)
		}

		if len(slice) != 0 {
			m[name] = slice
		}
	}

	return m, nil
}

// GetService gets a service by name and version
func (i *Index) GetService(name string, version string) *Service {
	s, ok := i.Services[name]
	if !ok {
		return nil
	}

	v, _ := s.Services[version]
	return v
}

// Clean index
func (i *Index) Clean() ([]*registry.Result, error) {
	ctx := context.TODO()
	now := time.Now().UnixNano()

	if i.Services == nil {
		return nil, nil
	}

	for _, services := range i.Services {
		if services.Services == nil {
			continue
		}

		for _, service := range services.Services {
			if service.Nodes == nil {
				continue
			}

			for _, node := range service.Nodes {
				if node.Expiry != 0 && node.Expiry < now {
					node.Enabled = false
				}
			}
		}
	}

	// Merge into self to trigger update
	return i.Merge(ctx, i)
}

// Merge one service into another
func (s *Service) Merge(ctx context.Context, name string, version string, merge *Service, diff *[]*registry.Result) error {
	// Unmarshal service
	s1, err := getService(nil, name, version, s.Raw)
	if err != nil {
		return err
	}

	s2, err := getService(ctx, name, version, merge.Raw)
	if err != nil {
		return err
	}

	if s.Raw == nil {
		s1 = &registry.Service{}
	} else {
		if err := msgpack.Unmarshal(s.Raw, &s1); err != nil {
			return err
		}
	}

	if merge.Raw == nil {
		s2 = &registry.Service{}
	} else {
		if err := msgpack.Unmarshal(merge.Raw, &s2); err != nil {
			return err
		}
	}

	// Has the service changed
	changed := false
	count := len(s1.Nodes)

	// If s2 is more recent than s1 then replace meta data
	if s.Mod <= merge.Mod {
		s1.Name = s2.Name
		s1.Version = s2.Version
		s1.Endpoints = s2.Endpoints
		s1.Metadata = s2.Metadata

		changed = true
	}

	// For each node in service
	for id, n2 := range merge.Nodes {
		n1, ok := s.Nodes[id]
		if !ok {
			s.Nodes[id] = n2
			n1 = n2
		}

		// If n2 is newer than n1 then replace meta data
		if n1.Mod <= n2.Mod {
			n1.Enabled = n2.Enabled
			n1.Mod = n2.Mod
			n1.Expiry = n2.Expiry
			changed = true

			i, node1 := NodeByID(s1.Nodes, id)
			_, node2 := NodeByID(s2.Nodes, id)

			// If non existent then insert
			if i == -1 && n2.Enabled {
				s1.Nodes = append(s1.Nodes, node2)
				continue
			} else if i != -1 {
				// If enabled replace node meta data, else delete
				if n2.Enabled {
					node1.Address = node2.Address
					node1.Metadata = node2.Metadata
					node1.Port = node2.Port
				} else {
					s1.Nodes = append(s1.Nodes[:i], s1.Nodes[i+1:]...)
				}
			}
		}
	}

	if changed {
		if len(s1.Nodes) == 0 && count != 0 {
			*diff = append(*diff, &registry.Result{
				Action:  "delete",
				Service: s1,
			})
		} else if count == 0 {
			*diff = append(*diff, &registry.Result{
				Action:  "create",
				Service: s1,
			})
		} else {
			*diff = append(*diff, &registry.Result{
				Action:  "update",
				Service: s1,
			})
		}
	}

	// Marshal service and store
	raw, err := msgpack.Marshal(s1)
	if err != nil {
		return err
	}

	// Update state
	s.Raw = raw

	// Return
	return nil
}

// NodeByID is a helper to loop through nodes to find a specific one
func NodeByID(nodes []*registry.Node, id string) (int, *registry.Node) {
	for i, node := range nodes {
		if node.Id == id {
			return i, node
		}
	}
	return -1, nil
}

type contextServiceID struct {
	Name    string
	Version string
}

func withService(ctx context.Context, service *registry.Service) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}

	id := contextServiceID{
		Name:    service.Name,
		Version: service.Version,
	}

	return context.WithValue(ctx, id, service)
}

func getService(ctx context.Context, name string, version string, raw []byte) (*registry.Service, error) {
	service := &registry.Service{}
	id := contextServiceID{
		Name:    name,
		Version: version,
	}

	if ctx != nil {
		if service, ok := ctx.Value(id).(*registry.Service); ok {
			return service, nil
		}
	}

	if len(raw) == 0 {
		return &registry.Service{}, nil
	}

	if err := msgpack.Unmarshal(raw, service); err != nil {
		return nil, err
	}

	return service, nil
}
