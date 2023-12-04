package check

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrModuleNameConflict = errors.New("module name conflict")
	ErrNoSuchModule       = errors.New("no module of given name known")
)

func NewRegistry() *Registry {
	return &Registry{
		mods: make(map[string]*Module),
	}
}

type (
	Registry struct {
		lock sync.RWMutex
		mods map[string]*Module
	}
)

func (r *Registry) Register(modules ...*Module) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, mod := range modules {
		modName := strings.ToLower(mod.Name())

		if _, ok := r.mods[modName]; ok {
			return fmt.Errorf("%w: %s", ErrModuleNameConflict, modName)
		}

		r.mods[modName] = mod
	}

	return nil
}

//nolint:ireturn // required to implement interface
func (r *Registry) Lookup(modName string) (CheckerLookup, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	modName = strings.ToLower(modName)

	if mod, ok := r.mods[modName]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchModule, modName)
	} else {
		return mod, nil
	}
}
