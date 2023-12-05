package mmdbmeld

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the geoip build config.
type Config struct {
	Databases []DatabaseConfig `yaml:"databases"`
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
	FloatDecimals  int  `yaml:"floatDecimals"`
	ForceIPVersion bool `yaml:"forceIPVersion"`
	MaxPrefix      int  `yaml:"maxPrefix"`
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
