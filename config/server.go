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
	ServerTypePostgres
	ServerTypeMysql

	//nolint:lll // cannot break regex
	// language=regexp
	urlSchemeRegex = `^(?P<scheme>\w+)://((?P<username>[\w-.]+)(:(?P<password>.*))?@)?(\{(?P<hosts>(.+:\d{1,5})(,(.+:\d{1,5}))*)}|(?P<host>.+:\d{1,5}))(?P<path>/.*$)?`
)

//nolint:goconst
func (t ServerType) Scheme() string {
	switch t {
	case ServerTypeRedis:
		return "redis"
	case ServerTypePostgres:
		return "postgres"
	case ServerTypeMysql:
		return "mysql"
	case ServerTypeUnspecified:
		fallthrough
	default:
		return "unknown"
	}
}

func (t ServerType) Driver() string {
	switch t {
	case ServerTypeRedis:
		return "redis"
	case ServerTypePostgres:
		return "pgx"
	case ServerTypeMysql:
		return "mysql"
	case ServerTypeUnspecified:
		fallthrough
	default:
		return ""
	}
}

var hostsRegexp = regexp.MustCompile(urlSchemeRegex)

type Credentials struct {
	Username string
	Password *string
}

func (c *Credentials) appendToBuilder(builder *strings.Builder) {
	if c == nil {
		return
	}

	_, _ = builder.WriteString(c.Username)
	if c.Password != nil {
		_, _ = builder.WriteString(":")
		_, _ = builder.WriteString(*c.Password)
	}
	_, _ = builder.WriteString("@")
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

func (s *Server) Query() string {
	if s.Args == nil || len(s.Args) < 1 {
		return ""
	}

	var builder strings.Builder

	for k, v := range s.Args {
		if builder.Len() > 0 {
			_, _ = builder.WriteString("&")
		}
		_, _ = builder.WriteString(k)
		_, _ = builder.WriteString("=")
		_, _ = builder.WriteString(fmt.Sprintf("%v", v))
	}

	return builder.String()
}

func (s *Server) DSNs() (result []string) {
	var builder strings.Builder
	result = make([]string, 0, len(s.Hosts))
	joinedPath := strings.Join(s.Path, "/")

	for i := range s.Hosts {
		builder.Reset()
		s.Credentials.appendToBuilder(&builder)

		_, _ = builder.WriteString("tcp(")
		_, _ = builder.WriteString(s.Hosts[i])
		_, _ = builder.WriteString(")")

		if len(s.Path) > 0 {
			_, _ = builder.WriteString("/")

			_, _ = builder.WriteString(joinedPath)
		}

		if query := s.Query(); query != "" {
			_, _ = builder.WriteString("?")
			_, _ = builder.WriteString(query)
		}

		result = append(result, builder.String())
	}

	return result
}

func (s *Server) ConnectionStrings() (result []string) {
	var builder strings.Builder
	result = make([]string, 0, len(s.Hosts))
	joinedPath := strings.Join(s.Path, "/")

	for i := range s.Hosts {
		builder.Reset()

		_, _ = builder.WriteString(s.Type.Scheme())
		_, _ = builder.WriteString("://")

		s.Credentials.appendToBuilder(&builder)

		_, _ = builder.WriteString(s.Hosts[i])
		if len(s.Path) > 0 {
			_, _ = builder.WriteString("/")

			_, _ = builder.WriteString(joinedPath)
		}

		if query := s.Query(); query != "" {
			_, _ = builder.WriteString("?")
			_, _ = builder.WriteString(query)
		}

		result = append(result, builder.String())
	}

	return result
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

func (s *Server) UnmarshalURL(rawURL string) error {
	rawPath, username, password, err := s.extractBaseProperties(rawURL)
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
		parsedURL, err := url.Parse(fmt.Sprintf("%s://%s%s", s.Type.Scheme(), s.Hosts[0], rawPath))
		if err != nil {
			return err
		}
		if err := s.unmarshalPath(parsedURL); err != nil {
			return err
		}
	} else {
		s.Args = make(map[string]any)
	}

	return nil
}

func (s *Server) extractBaseProperties(rawURL string) (rawPath, username, password string, err error) {
	allMatches := hostsRegexp.FindAllStringSubmatch(rawURL, -1)
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
