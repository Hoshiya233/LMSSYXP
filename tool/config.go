package tool

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server `yaml:"server"`
	// Database Database `yaml:"database"`
	MMDpaths []string `yaml:"mmdpaths"`
}

type Server struct {
	Port string `yaml:"port"`
}

// type Database struct {
// 	Host     string `yaml:"host"`
// 	Port     string `yaml:"port"`
// 	User     string `yaml:"user"`
// 	Password string `yaml:"password"`
// }

var MMDPaths []string

func ParseConfig(path string) *Config {
	var config *Config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	MMDPaths = config.MMDpaths
	return config
}
