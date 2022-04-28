package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrServerNameAlreadyRegistered = errors.New("server name is already registered")
	ErrNoSuchServer                = errors.New("no known server with given name")
	DefaultLookup                  = &ServerLookup{
		servers: make(map[string]Server),
	}
)

type ServerLookup struct {
	lock    sync.RWMutex
	servers map[string]Server
}

func (l *ServerLookup) Register(name string, srv Server) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	name = strings.ToLower(name)

	if _, ok := l.servers[name]; ok {
		return fmt.Errorf("%w: %s", ErrServerNameAlreadyRegistered, name)
	}

	l.servers[name] = srv
	return nil
}

func (l *ServerLookup) Lookup(name string) (*Server, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	name = strings.ToLower(name)

	match, ok := l.servers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchServer, name)
	}
	return &match, nil
}
