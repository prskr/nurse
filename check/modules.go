package check

import (
	"fmt"
	"strings"
	"sync"

	"github.com/baez90/nurse/config"
	"github.com/baez90/nurse/grammar"
)

type (
	ModuleOption interface {
		Apply(m *Module) error
	}

	ModuleOptionFunc func(m *Module) error

	Factory interface {
		New() SystemChecker
	}

	FactoryFunc func() SystemChecker
)

func (f ModuleOptionFunc) Apply(m *Module) error {
	return f(m)
}

func (f FactoryFunc) New() SystemChecker {
	return f()
}

func WithCheck(name string, factory Factory) ModuleOption {
	return ModuleOptionFunc(func(m *Module) error {
		return m.Register(name, factory)
	})
}

func NewModule(opts ...ModuleOption) (*Module, error) {
	m := &Module{
		knownChecks: make(map[string]Factory),
	}

	for i := range opts {
		if err := opts[i].Apply(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

type Module struct {
	lock        sync.RWMutex
	knownChecks map[string]Factory
}

func (m *Module) Lookup(c grammar.Check, srvLookup config.ServerLookup) (SystemChecker, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		factory Factory
		ok      bool
	)
	if factory, ok = m.knownChecks[strings.ToLower(c.Initiator.Name)]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchCheck, c.Initiator.Name)
	}

	chk := factory.New()
	if err := chk.UnmarshalCheck(c, srvLookup); err != nil {
		return nil, err
	}

	return chk, nil
}

func (m *Module) Register(name string, factory Factory) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	name = strings.ToLower(name)

	if _, ok := m.knownChecks[name]; ok {
		return fmt.Errorf("%w: %s", ErrConflictingCheck, name)
	}

	m.knownChecks[name] = factory

	return nil
}
