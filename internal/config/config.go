package config

import (
	"errors"
	"fmt"
	"github.com/vrunoa/bouncer/internal/unit"
	yaml "gopkg.in/yaml.v2"
	"os"
	"regexp"
	"strconv"
)

type Configuration struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Image      Image  `yaml:"image"`
}

type Image struct {
	Name   string      `yaml:"name"`
	Policy GuardPolicy `yaml:"policy"`
}

type GuardPolicy struct {
	Deny []DenyPolicy `yaml:"deny"`
}

type DenyPolicy struct {
	Desc      string `yaml:"desc"`
	Size      string `yaml:"size"`
	FloatSize float64
	Unit      unit.Unit
}

func ReadYaml(filePath string) (*Configuration, error) {
	if filePath == "" {
		return &Configuration{}, errors.New("empty file name")
	}
	f, err := os.Open(filePath)
	if err != nil {
		return &Configuration{}, err
	}
	defer f.Close()

	var c Configuration
	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return &Configuration{}, err
	}
	if err = c.Validate(); err != nil {
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
	for ii, deny := range c.Image.Policy.Deny {
		sizeFloat, sizeUnit, err := ParseSize(deny.Size)
		if err != nil {
			return err
		}
		c.Image.Policy.Deny[ii].FloatSize = sizeFloat
		c.Image.Policy.Deny[ii].Unit = sizeUnit
	}
	return nil
}

// ParseSize validates and parse string size
func ParseSize(size string) (float64, unit.Unit, error) {
	mi := string(unit.Mega)
	gi := string(unit.Giga)
	regex := fmt.Sprintf("^(?P<Size>[0-9]+(\\.[0-9])?)(?P<Unit>(%s|%s))", mi, gi)
	reg := regexp.MustCompile(regex)
	if ok := reg.MatchString(size); !ok {
		return 0, unit.Unsupported, errors.New("invalid size")
	}
	match := reg.FindStringSubmatch(size)
	paramsMap := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	sizeStr := paramsMap["Size"]
	s, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0, unit.Unsupported, err
	}
	unitStr := paramsMap["Unit"]
	unt := unit.Unsupported
	if unitStr == string(unit.Giga) {
		unt = unit.Giga
	} else if unitStr == string(unit.Mega) {
		unt = unit.Mega
	}
	return s, unt, nil
}
