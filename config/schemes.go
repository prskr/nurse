package config

import "strings"

var scheme2ServerType = map[string]ServerType{
	"redis": ServerTypeRedis,
}

func SchemeToServerType(scheme string) ServerType {
	if match, ok := scheme2ServerType[strings.ToLower(scheme)]; ok {
		return match
	}
	return ServerTypeUnspecified
}
