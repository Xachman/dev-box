package main

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	VolumeDir string `yaml:"volumeDir"`
	Namespace string
}

func GetConfig() Config {
	config := Config{}
	fileContents, err := ioutil.ReadFile(fmt.Sprintf("../data/config.yml"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	yErr := yaml.Unmarshal(fileContents, &config)
	if yErr != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("Volume dir: %s \n", config.VolumeDir)
	return config
}

func (c *Config) GetVolumeDir() string {
	return c.VolumeDir
}
func (c *Config) GetNamespace() string {
	return c.Namespace
}
