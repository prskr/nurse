package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultCheckTimeout = 500 * time.Millisecond
	defaultAttemptCount = 20
)

func ConfigureFlags(cfg *Nurse) *flag.FlagSet {
	set := flag.NewFlagSet("nurse", flag.ContinueOnError)

	set.DurationVar(
		&cfg.CheckTimeout,
		"check-timeout",
		LookupEnvOr("NURSE_CHECK_TIMEOUT", defaultCheckTimeout, time.ParseDuration),
		"Timeout when running checks",
	)

	set.UintVar(
		&cfg.CheckAttempts,
		"check-attempts",
		LookupEnvOr("NURSE_CHECK_ATTEMPTS", defaultAttemptCount, parseUint),
		"Number of attempts for a check",
	)

	return set
}

//nolint:ireturn // false positive
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

//nolint:ireturn // false positive
func Identity[T any](in T) (T, error) {
	return in, nil
}

func parseUint(val string) (uint, error) {
	parsed, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
