package config

import (
	"context"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Configs map[string]interface{} `yaml:"configs"`
}

var cfg config

// LoadConfig reads the configuration from a YAML file.
// The file path can be specified via a command-line flag "-config".
// If the flag is not provided, it defaults to "test".
// Returns an error if reading or parsing the file fails.
func LoadConfig(configFlag string) {
	cfg = config{
		Configs: make(map[string]interface{}),
	}

	filePath := "./config/config_" + configFlag + ".yml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Fatal(context.Background(), err)
	}
}

// GetValue fetches the configuration value associated with the given key.
// If the key does not exist, it returns an empty Value.
// To differentiate between missing and empty values, use LookupValue.
func GetValue(key string) Value {
	val, _ := LookupValue(key)
	return val
}

// LookupValue checks if the configuration contains the specified key.
// Returns the Value and true if found; otherwise, returns an empty Value and false.
func LookupValue(key string) (Value, bool) {
	val, exists := cfg.Configs[key]
	return value{
		value:  val,
		exists: exists,
	}, exists
}
