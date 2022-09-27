package config

import (
	"os"
	"path"
	"strings"
)

const (
	ServerKeyPrefix   = "NURSE_SERVER"
	EndpointKeyPrefix = "NURSE_ENDPOINT"
)

func ServersFromEnv() (map[string]Server, error) {
	servers := make(map[string]Server)
	for _, kv := range os.Environ() {
		key, value, valid := strings.Cut(kv, "=")
		if !valid {
			continue
		}

		if !strings.HasPrefix(key, ServerKeyPrefix) {
			continue
		}

		serverName := strings.ToLower(strings.Trim(strings.Replace(key, ServerKeyPrefix, "", -1), "_"))
		srv := Server{}
		if err := srv.UnmarshalURL(value); err != nil {
			return nil, err
		}

		servers[serverName] = srv
	}

	return servers, nil
}

func EndpointsFromEnv() (map[Route]EndpointSpec, error) {
	endpoints := make(map[Route]EndpointSpec)

	for _, kv := range os.Environ() {
		key, value, valid := strings.Cut(kv, "=")
		if !valid {
			continue
		}

		if !strings.HasPrefix(key, EndpointKeyPrefix) {
			continue
		}

		endpointRoute := path.Join(strings.Split(strings.ToLower(strings.Trim(strings.Replace(key, EndpointKeyPrefix, "", -1), "_")), "_")...)
		spec := EndpointSpec{}
		if err := spec.Parse(value); err != nil {
			return nil, err
		}

		endpoints[Route(endpointRoute)] = spec
	}

	return endpoints, nil
}
