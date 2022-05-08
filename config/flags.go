package config

import (
	"flag"
	"os"
	"time"
)

const defaultCheckTimeout = 500 * time.Millisecond

func ConfigureFlags(cfg *Nurse) *flag.FlagSet {
	set := flag.NewFlagSet("nurse", flag.ContinueOnError)

	set.DurationVar(
		&cfg.CheckTimeout,
		"check-timeout",
		LookupEnvOr("NURSE_CHECK_TIMEOUT", defaultCheckTimeout, time.ParseDuration),
		"Timeout when running checks",
	)

	return set
}

func LookupEnvOr[T any](envKey string, fallback T, parse func(envVal string) (T, error)) T {
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

func Identity[T any](in T) (T, error) {
	return in, nil
}
