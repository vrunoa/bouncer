package config

import (
	"errors"
	"github.com/vrunoa/bouncer/internal/unit"
	"gopkg.in/yaml.v2"
	"os"
)

type Configuration struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Image      DockerImage `yaml:"image"`
}

type DockerImage struct {
	Name   string      `yaml:"name"`
	Policy GuardPolicy `yaml:"policy"`
}

type GuardPolicy struct {
	Deny []DenyPolicy `yaml:"deny"`
}

type DenyPolicy struct {
	Desc string `yaml:"desc"`
	Size string `yaml:"size"`
}

func ReadYaml(filePath string) (*Configuration, error) {
	if filePath == "" {
		return &Configuration{}, errors.New("empty file name")
	}
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return &Configuration{}, err
	}
	var c Configuration
	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return &Configuration{}, err
	}
	return &c, nil
}

func (c *Configuration) Validate() error {
	if c.Image.Name == "" {
		return errors.New("missing image name")
	}
	if len(c.Image.Policy.Deny) == 0 {
		return errors.New("missing deny policies")
	}
	for _, deny := range c.Image.Policy.Deny {
		if _, _, err := unit.ParseSize(deny.Size); err != nil {
			return err
		}
	}
	return nil
}
