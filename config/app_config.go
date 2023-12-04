package config

import (
	"io"
	"os"
	"time"
)

type (
	Option interface {
		Extend(n *Nurse) error
	}
	OptionFunc func(n *Nurse) error
)

func (f OptionFunc) Extend(n *Nurse) error {
	return f(n)
}

func New(opts ...Option) (*Nurse, error) {
	var (
		inst = new(Nurse)
		err  error
	)

	for i := range opts {
		if err = opts[i].Extend(inst); err != nil {
			return nil, err
		}
	}

	return inst, nil
}

type Nurse struct {
	Servers       map[string]Server
	Endpoints     map[Route]EndpointSpec
	CheckTimeout  time.Duration
	CheckAttempts uint
	Insecure      bool
}

func (n *Nurse) ServerLookup() (*ServerRegister, error) {
	register := NewServerRegister()

	for name, srv := range n.Servers {
		if err := register.Register(name, srv); err != nil {
			return nil, err
		}
	}

	return register, nil
}

// Merge merges the current Nurse instance with another one
// giving the current instance precedence means no set value is overwritten
func (n *Nurse) Merge(other Nurse) {
	if n.CheckTimeout == 0 {
		n.CheckTimeout = other.CheckTimeout
	}

	if n.CheckAttempts == 0 {
		n.CheckAttempts = other.CheckAttempts
	}

	for name, srv := range other.Servers {
		if _, ok := n.Servers[name]; !ok {
			n.Servers[name] = srv
		}
	}
}

func (n *Nurse) Unmarshal(reader io.ReadSeeker) error {
	providers := []func(io.Reader) configDecoder{
		newYAMLDecoder,
		newJSONDecoder,
	}

	for i := range providers {
		if err := n.tryUnmarshal(reader, providers[i]); err == nil {
			return nil
		}
	}

	return nil
}

func (n *Nurse) ReadFromFile(configFilePath string) error {
	if file, err := os.Open(configFilePath); err != nil {
		return err
	} else {
		defer func() {
			_ = file.Close()
		}()

		return n.Unmarshal(file)
	}
}

func (n *Nurse) tryUnmarshal(seeker io.ReadSeeker, prov func(reader io.Reader) configDecoder) error {
	if _, err := seeker.Seek(0, 0); err != nil {
		return err
	}

	decoder := prov(seeker)
	return decoder.DecodeConfig(n)
}
