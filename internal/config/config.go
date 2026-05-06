package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// SSHConfig holds SSH connection parameters.
type SSHConfig struct {
	User    string        `yaml:"user"`
	KeyFile string        `yaml:"key_file"`
	Timeout time.Duration `yaml:"timeout"`
}

// FilterConfig defines include/exclude glob patterns for file selection.
type FilterConfig struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

// Config is the top-level configuration structure for confsnap.
type Config struct {
	Hosts    []string     `yaml:"hosts"`
	Files    []string     `yaml:"files"`
	SSH      SSHConfig    `yaml:"ssh"`
	Filter   FilterConfig `yaml:"filter"`
	StoreDir string       `yaml:"store_dir"`
	Output   string       `yaml:"output"` // "text" or "json"
}

// defaults applies sensible default values to fields that were not set.
func (c *Config) defaults() {
	if c.SSH.User == "" {
		c.SSH.User = "root"
	}
	if c.SSH.Timeout == 0 {
		c.SSH.Timeout = 10 * time.Second
	}
	if c.StoreDir == "" {
		c.StoreDir = ".confsnap"
	}
	if c.Output == "" {
		c.Output = "text"
	}
}

// validate checks that the configuration contains the minimum required fields.
func (c *Config) validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("config: at least one host is required")
	}
	if len(c.Files) == 0 && len(c.Filter.Include) == 0 {
		return errors.New("config: at least one file or include filter is required")
	}
	if c.Output != "text" && c.Output != "json" {
		return errors.New("config: output must be \"text\" or \"json\"")
	}
	return nil
}

// Load reads and parses a YAML configuration file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.defaults()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
