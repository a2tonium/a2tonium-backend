package config

import (
	"context"
	"flag"
	"os"

	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"gopkg.in/yaml.v3"
)

type config struct {
	Configs map[string]interface{} `yaml:"configs"`
}

var cfg config

func init() {
	// Initialize the config with an empty map
	cfg = config{
		Configs: make(map[string]interface{}),
	}

	// Load configuration from file
	if err := LoadConfig(); err != nil {
		logger.Fatal(context.Background(), err)
	}
}

// LoadConfig reads the configuration from a YAML file.
// The file path can be specified via a command-line flag "-config".
// If the flag is not provided, it defaults to "test".
// Returns an error if reading or parsing the file fails.
func LoadConfig() error {
	configFlag := flag.String("config", "test", "configuration file suffix")
	flag.Parse()

	filePath := "./config/config_" + *configFlag + ".yml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return nil
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
