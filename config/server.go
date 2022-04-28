package config

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"

	"github.com/baez90/nurse/internal/values"
)

type ServerType uint

const (
	ServerTypeUnspecified ServerType = iota
	ServerTypeRedis
)

var hostsRegexp = regexp.MustCompile(`^{(.+:\d{1,5})(,(.+:\d{1,5}))*}|(.+:\d{1,5})$`)

func ParseFromURL(url *url.URL) (*Server, error) {
	srv := &Server{
		Type:  SchemeToServerType(url.Scheme),
		Hosts: hostsRegexp.FindAllString(url.Host, -1),
	}

	if user := url.User; user != nil {
		srv.Credentials = &Credentials{
			Username: user.Username(),
		}
		if pw, ok := user.Password(); ok {
			srv.Credentials.Password = values.StringP(pw)
		}
	}

	srv.Path = strings.Split(strings.Trim(url.EscapedPath(), "/"), "/")

	q := url.Query()
	qm := map[string][]string(q)
	srv.Args = make(map[string]any, len(qm))

	for k := range qm {
		var val any
		if err := json.Unmarshal([]byte(q.Get(k)), &val); err != nil {
			return nil, err
		} else {
			srv.Args[k] = val
		}
	}

	return srv, nil
}

type Credentials struct {
	Username string
	Password *string
}

type Server struct {
	Type        ServerType
	Credentials *Credentials
	Hosts       []string
	Path        []string
	Args        map[string]any
}
