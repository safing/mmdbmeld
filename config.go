package mmdbmeld

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the geoip build config.
type Config struct {
	Databases []DatabaseConfig `yaml:"databases"`
	Defaults  DefaultConfig    `yaml:"defaults"`
}

// DatabaseConfig holds the config for building one database.
type DatabaseConfig struct {
	Name     string            `yaml:"name"`
	MMDB     MMDBConfig        `yaml:"mmdb"`
	Types    map[string]string `yaml:"types"`
	Inputs   []DatabaseInput   `yaml:"inputs"`
	Output   string            `yaml:"output"`
	Optimize Optimizations     `yaml:"optimize"`
	Merge    MergeConfig       `yaml:"merge"`
}

// DefaultConfig holds a subset of DatabaseConfig fields to be used as a default config.
type DefaultConfig struct {
	Types    map[string]string `yaml:"types"`
	Optimize Optimizations     `yaml:"optimize"`
	Merge    MergeConfig       `yaml:"merge"`
}

// MMDBConfig holds mmdb specific config.
type MMDBConfig struct {
	IPVersion  int `yaml:"ipVersion"`
	RecordSize int `yaml:"recordSize"`
}

// DatabaseInput holds database input config.
type DatabaseInput struct {
	File     string            `yaml:"file"`
	Fields   []string          `yaml:"fields"`
	FieldMap map[string]string `yaml:"fieldMap"`
}

// Optimizations holds optimization config.
type Optimizations struct {
	FloatDecimals  int   `yaml:"floatDecimals"`
	ForceIPVersion *bool `yaml:"forceIPVersion"`
	MaxPrefix      int   `yaml:"maxPrefix"`
}

// ForceIPVersionEnabled reports whether ForceIPVersion is set and true.
func (o Optimizations) ForceIPVersionEnabled() bool {
	return o.ForceIPVersion != nil && *o.ForceIPVersion
}

// MergeConfig holds merge configuration.
type MergeConfig struct {
	ConditionalResets []ConditionalResetConfig `yaml:"conditionalResets"`
}

// ConditionalResetConfig defines a conditional reset merge config.
type ConditionalResetConfig struct {
	IfChanged []string `yaml:"ifChanged"`
	Reset     []string `yaml:"reset"`
}

// LoadConfig loads a configuration file.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	return config, yaml.Unmarshal(data, config)
}

// ApplyTo applies the default config to the given database config.
func (d DefaultConfig) ApplyTo(c *DatabaseConfig) {
	// Add all missing default types.
	if c.Types == nil {
		c.Types = make(map[string]string)
	}
	for k, v := range d.Types {
		_, ok := c.Types[k]
		if !ok {
			c.Types[k] = v
		}
	}

	// Apply Optimizations.
	if c.Optimize.FloatDecimals == 0 && d.Optimize.FloatDecimals != 0 {
		c.Optimize.FloatDecimals = d.Optimize.FloatDecimals
	}
	if c.Optimize.ForceIPVersion == nil && d.Optimize.ForceIPVersion != nil {
		c.Optimize.ForceIPVersion = d.Optimize.ForceIPVersion
	}
	if c.Optimize.MaxPrefix == 0 && d.Optimize.MaxPrefix != 0 {
		c.Optimize.MaxPrefix = d.Optimize.MaxPrefix
	}

	// Apply Merge Config.
	if len(c.Merge.ConditionalResets) == 0 && len(d.Merge.ConditionalResets) != 0 {
		c.Merge.ConditionalResets = d.Merge.ConditionalResets
	}
}
