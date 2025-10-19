package config

import (
	"log"
	"os"
	"strings"
)

// Config holds runtime configuration for the server
// Values are read once at startup via FromEnv.
type Config struct {
	Port           string   // HTTP port, e.g. ":8080"
	AllowedOrigins []string // CORS/WS allowed origins; ["*"] to allow all
	EnablePprof    bool     // enable /debug/pprof endpoints
	LogLevel       string   // debug, info, warn, error
}

// FromEnv loads configuration from environment variables with sensible defaults.
// PORT: string, default "8080"
// ALLOWED_ORIGIN: string, default "*"
func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Backward compat: ALLOWED_ORIGIN (single) or ALLOWED_ORIGINS (comma separated)
	allowedSingle := os.Getenv("ALLOWED_ORIGIN")
	allowedMulti := os.Getenv("ALLOWED_ORIGINS")
	var allowed []string
	if allowedMulti != "" {
		// split by comma
		for _, v := range strings.Split(allowedMulti, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				allowed = append(allowed, v)
			}
		}
	} else if allowedSingle != "" {
		allowed = []string{allowedSingle}
	} else {
		allowed = []string{"*"}
	}
	enablePprof := false
	if v := os.Getenv("ENABLE_PPROF"); v == "1" || v == "true" || v == "TRUE" {
		enablePprof = true
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" { logLevel = "info" }
	log.Printf("Config: PORT=%s ALLOWED_ORIGINS=%v ENABLE_PPROF=%v LOG_LEVEL=%s", port, allowed, enablePprof, logLevel)
	return Config{
		Port:           ":" + port,
		AllowedOrigins: allowed,
		EnablePprof:    enablePprof,
		LogLevel:       logLevel,
	}
}
