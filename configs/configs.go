package configs

import (
	"fmt"
	"os"
  "github.com/goccy/go-yaml"
)

type Config struct {
	APIEndpoint struct {
		URL string `yaml:"url"`
		Origin string `yaml:"origin"`
	} `yaml:"api_endpoint"`
}

func LoadConfig(path string)(*Config, error){
    file, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("error reading config file: %v", err)
    }
    var config Config
    err = yaml.Unmarshal(file, &config)
    if err != nil {
        return nil, fmt.Errorf("error parsing config file: %v", err)
    }
    return &config, nil
}
