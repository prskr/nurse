package main

import (
	"os"
	"time"
)

var checkTimeout time.Duration

func main() {
}

func lookupEnvOr[T any](envKey string, fallback T, parse func(envVal string) (T, error)) T {
	envVal := os.Getenv(envKey)
	if envVal == "" {
		return fallback
	}

	if parsed, err := parse(envVal); err != nil {
		return fallback
	} else {
		return parsed
	}
}

func identity[T any](in T) (T, error) {
	return in, nil
}
