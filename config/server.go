package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type ServerType uint

const (
	ServerTypeUnspecified ServerType = iota
	ServerTypeRedis

	//nolint:lll // cannot break regex
	// language=regexp
	urlSchemeRegex = `^(?P<scheme>\w+)://((?P<username>[\w-.]+)(:(?P<password>.*))?@)?(\{(?P<hosts>(.+:\d{1,5})(,(.+:\d{1,5}))*)}|(?P<host>.+:\d{1,5}))(?P<path>/.*$)?`
)

func (t ServerType) Scheme() string {
	switch t {
	case ServerTypeRedis:
		return "redis"
	case ServerTypeUnspecified:
		fallthrough
	default:
		return "unknown"
	}
}

var hostsRegexp = regexp.MustCompile(urlSchemeRegex)

type Credentials struct {
	Username string
	Password *string
}

var (
	_ json.Unmarshaler = (*Server)(nil)
	_ yaml.Unmarshaler = (*Server)(nil)
)

type marshaledServer struct {
	Type  ServerType
	URL   string
	Hosts []string
	Args  map[string]any
}

type Server struct {
	Type        ServerType
	Credentials *Credentials
	Hosts       []string
	Path        []string
	Args        map[string]any
}

func (s *Server) UnmarshalYAML(value *yaml.Node) error {
	s.Args = make(map[string]any)

	tmp := new(marshaledServer)

	if err := value.Decode(tmp); err != nil {
		return err
	}

	return s.mergedMarshaledServer(*tmp)
}

func (s *Server) UnmarshalJSON(bytes []byte) error {
	s.Args = make(map[string]any)

	tmp := new(marshaledServer)

	if err := json.Unmarshal(bytes, tmp); err != nil {
		return err
	}

	return s.mergedMarshaledServer(*tmp)
}

func (s *Server) UnmarshalURL(rawUrl string) error {
	rawPath, username, password, err := s.extractBaseProperties(rawUrl)
	if err != nil {
		return err
	}

	if username != "" {
		s.Credentials = &Credentials{
			Username: username,
		}
		if password != "" {
			s.Credentials.Password = &password
		}
	}

	if rawPath != "" {
		parsedUrl, err := url.Parse(fmt.Sprintf("%s://%s%s", s.Type.Scheme(), s.Hosts[0], rawPath))
		if err != nil {
			return err
		}
		if err := s.unmarshalPath(parsedUrl); err != nil {
			return err
		}
	} else {
		s.Args = make(map[string]any)
	}

	return nil
}

func (s *Server) extractBaseProperties(rawUrl string) (rawPath, username, password string, err error) {
	allMatches := hostsRegexp.FindAllStringSubmatch(rawUrl, -1)
	if matchLen := len(allMatches); matchLen != 1 {
		return "", "", "", fmt.Errorf("ambiguous server match: %d", matchLen)
	}

	match := allMatches[0]

	for i, name := range hostsRegexp.SubexpNames() {
		if i == 0 {
			continue
		}

		switch name {
		case "scheme":
			s.Type = SchemeToServerType(match[i])
		case "host":
			singleHost := match[i]
			if singleHost != "" {
				s.Hosts = []string{singleHost}
			}
		case "hosts":
			hosts := strings.Split(match[i], ",")
			for _, host := range hosts {
				s.Hosts = append(s.Hosts, strings.TrimSpace(host))
			}
		case "path":
			rawPath = match[i]
		case "username":
			username = match[i]
		case "password":
			password = match[i]
		}
	}

	return rawPath, username, password, nil
}

func (s *Server) unmarshalPath(u *url.URL) error {
	s.Path = strings.Split(strings.Trim(u.EscapedPath(), "/"), "/")

	q := u.Query()
	qm := map[string][]string(q)
	s.Args = make(map[string]any, len(qm))

	for k := range qm {
		var val any
		if err := json.Unmarshal([]byte(q.Get(k)), &val); err != nil {
			return err
		} else {
			s.Args[k] = val
		}
	}

	return nil
}

func (s *Server) mergedMarshaledServer(srv marshaledServer) error {
	if srv.URL != "" {
		if err := s.UnmarshalURL(srv.URL); err != nil {
			return err
		}

		if srv.Args != nil {
			maps.Copy(s.Args, srv.Args)
		}
		return nil
	}

	s.Type = srv.Type
	s.Hosts = srv.Hosts
	if srv.Args != nil {
		s.Args = srv.Args
	}

	return nil
}
