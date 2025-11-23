package config

import (
	"fmt"
	"log" //nolint:depguard // using standard log before structured logging is initialized
	"os"
	"slices"
	"strconv"
	"time"
)

func getEnvWithFallback(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: environment variable %v not set, falling back to %v", key, fallback)
		return fallback
	}
	return val
}

func getEnvRequired(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %v is required but not set", key)
	}
	return val, nil
}

func getEnvWithFallbackAndValidOptions(key string, fallback string, validoptions ...string) string {
	val := os.Getenv(key)
	if !slices.Contains(validoptions, val) {
		log.Printf("WARNING: invalid value %v for env %v, falling back to %v", val, key, fallback)
		return fallback
	}
	return val
}

func getEnvIntWithFallback(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: env %v not set, falling back to %v", key, fallback)
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("WARNING: could not parse int from env %v, falling back to %v", key, fallback)
		return fallback
	}

	return intVal
}

// getEnvDurationWithFallback retrieves a duration from an environment variable,
// uses fallback if non-positive, not parseable, or not set.
func getEnvDurationWithFallback(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: env %v not set, falling back to %v", key, fallback)
		return fallback
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("WARNING: env %v is not a valid duration format (%v), falling back to %v", key, val, fallback)
		return fallback
	}

	if dur <= 0 {
		log.Printf("WARNING: env %v must be positive, got %v, falling back to %v", key, dur, fallback)
		return fallback
	}

	return dur
}

func getEnvWithFallbackAndCustomValidation(key string, fallback string, validateFunc func(string) bool) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("WARNING: env %v not set, falling back to %v", key, fallback)
		return fallback
	}

	if !validateFunc(val) {
		log.Printf("WARNING: invalid value %v for env %v, falling back to %v", val, key, fallback)
		return fallback
	}
	return val
}
