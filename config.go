package bugfixes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const DefaultServer = "https://api.bugfix.es/v1"

const DefaultTimeout = 10 * time.Second

var defaultHTTPClient = &http.Client{Timeout: DefaultTimeout}

type Config struct {
	Server      string
	AgentKey    string
	AgentSecret string
	LogLevel    string
	LocalOnly   bool
	HTTPClient  *http.Client
}

var (
	defaultConfigMu sync.RWMutex
	defaultConfig   *Config
)

func LoadConfigFromEnv() Config {
	localOnlyStr := strings.TrimSpace(os.Getenv("BUGFIXES_LOCAL_ONLY"))
	localOnly, err := strconv.ParseBool(localOnlyStr)
	if err != nil && localOnlyStr != "" {
		_, _ = fmt.Fprintf(os.Stderr, "bugfixes: invalid BUGFIXES_LOCAL_ONLY value %q, defaulting to false\n", localOnlyStr)

	}

	return Config{
		Server:      valueOrDefault(os.Getenv("BUGFIXES_SERVER"), DefaultServer),
		AgentKey:    os.Getenv("BUGFIXES_AGENT_KEY"),
		AgentSecret: os.Getenv("BUGFIXES_AGENT_SECRET"),
		LogLevel:    os.Getenv("BUGFIXES_LOG_LEVEL"),
		LocalOnly:   localOnly,
	}
}

func GetDefaultConfig() Config {
	defaultConfigMu.RLock()
	cfg := defaultConfig
	defaultConfigMu.RUnlock()
	if cfg != nil {
		return cfg.normalized()
	}

	return LoadConfigFromEnv().normalized()
}

func SetDefaultConfig(cfg Config) {
	cfg = cfg.normalized()

	defaultConfigMu.Lock()
	defer defaultConfigMu.Unlock()

	defaultConfig = &cfg
}

func ResetDefaultConfig() {
	defaultConfigMu.Lock()
	defer defaultConfigMu.Unlock()

	defaultConfig = nil
}

func (c Config) Merge(override Config) Config {
	merged := c

	if override.Server != "" {
		merged.Server = override.Server
	}
	if override.AgentKey != "" {
		merged.AgentKey = override.AgentKey
	}
	if override.AgentSecret != "" {
		merged.AgentSecret = override.AgentSecret
	}
	if override.LogLevel != "" {
		merged.LogLevel = override.LogLevel
	}
	if override.LocalOnly {
		merged.LocalOnly = true
	}
	if override.HTTPClient != nil {
		merged.HTTPClient = override.HTTPClient
	}

	return merged.normalized()
}

// GetHTTPClient returns the configured HTTP client, or a default client
// with a 10-second timeout.
func (c Config) GetHTTPClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return defaultHTTPClient
}

func (c Config) LogEndpoint() string {
	return strings.TrimRight(c.normalized().Server, "/") + "/log"
}

func (c Config) BugEndpoint() string {
	return strings.TrimRight(c.normalized().Server, "/") + "/bug"
}

func (c Config) normalized() Config {
	if c.Server == "" {
		c.Server = DefaultServer
	}

	return c
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
