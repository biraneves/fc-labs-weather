package config

import "time"

func ParseDuration(raw string) (time.Duration, error) {
	return parseDuration(raw)
}
