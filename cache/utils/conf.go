package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	URLList          []string `yaml:"URLs"`
	MinTimeout       int      `yaml:"MinTimeout"`
	MaxTimeout       int      `yaml:"MaxTimeout"`
	NumberOfRequests int      `yaml:"NumberOfRequests"`
}

func GetConf(filename string) *Conf {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("[ERROR] Read conf yaml file %s: %s", filename, err)
	}
	var c Conf
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("[ERROR] Unmarshal file %s: %v", filename, err)
	}

	return &c
}
