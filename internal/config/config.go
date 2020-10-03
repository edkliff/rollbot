package config

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Server   string   `yaml:"server"`
	LogLevel LogLevel `yaml:"loglevel"`
	Language Lang     `yaml:"language"`

	Dice struct {
		Count int64 `yaml:"count"`
		Dice  int64 `yaml:"dice"`
		Adder int64 `yaml:"adder"`
	} `yaml:"dice"`

	VK struct {
		ConfirmationResponse string `yaml:"confirmation_response"`
		Token                string `yaml:"token"`
		VKServer             string `yaml:"vkserver"`
		APIVersion           string `yaml:"apiversion"`
	} `yaml:"vk"`

	DB DBConfig `yaml:"storage"`
}

func  LoadConfig(filename string) (*Config, error) {
	c := Config{}
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type DBConfig struct {
	Kind string `yaml:"kind"`
	Filename string `yaml:"filename"`
}

type LogLevel uint8

const (
	Errors LogLevel = iota
	Debugs
)

func (ll *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s := ""
	err := unmarshal(&s)
	if err !=nil {
		return err
	}
	s = strings.Replace(s, "\"", "", -1)
	s = strings.ToLower(s)
	switch s {
	case "error":
		z := Errors
		*ll = z
	case "debug":
		z := Debugs
		*ll = z
	default:
		z := Errors
		*ll = z
	}
	return nil
}

type Lang uint8

const (
	Ru Lang = iota
	En
)

func (ll *Lang) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s := ""
	unmarshal(&s)
	s = strings.Replace(s, "\"", "", -1)
	s = strings.ToLower(s)
	switch s {
	case "ru":
		z := Ru
		*ll = z
	case "en":
		z := En
		*ll = z
	default:
		z := Ru
		*ll = z
	}
	return nil
}
