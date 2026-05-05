package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level confsnap configuration.
type Config struct {
	Defaults Defaults  `yaml:"defaults"`
	Hosts    []Host    `yaml:"hosts"`
	Files    []string  `yaml:"files"`
}

// Defaults holds default SSH connection settings.
type Defaults struct {
	User    string        `yaml:"user"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
	KeyFile string        `yaml:"key_file"`
}

// Host represents a single remote host to snapshot.
type Host struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	User    string `yaml:"user,omitempty"`
	Port    int    `yaml:"port,omitempty"`
	KeyFile string `yaml:"key_file,omitempty"`
}

// EffectiveUser returns the host-specific user or falls back to the default.
func (h Host) EffectiveUser(defaults Defaults) string {
	if h.User != "" {
		return h.User
	}
	return defaults.User
}

// EffectivePort returns the host-specific port or falls back to the default.
func (h Host) EffectivePort(defaults Defaults) int {
	if h.Port != 0 {
		return h.Port
	}
	if defaults.Port != 0 {
		return defaults.Port
	}
	return 22
}

// EffectiveKeyFile returns the host-specific key file or falls back to the default.
func (h Host) EffectiveKeyFile(defaults Defaults) string {
	if h.KeyFile != "" {
		return h.KeyFile
	}
	return defaults.KeyFile
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("at least one host must be defined")
	}
	if len(c.Files) == 0 {
		return fmt.Errorf("at least one file path must be defined")
	}
	for i, h := range c.Hosts {
		if h.Address == "" {
			return fmt.Errorf("host[%d] missing address", i)
		}
	}
	return nil
}
