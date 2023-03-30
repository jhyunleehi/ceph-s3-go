package s3

import (
	
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"	

	log "github.com/sirupsen/logrus"
)


type Config struct {
	Url       string `yaml:"url"`
	AccessKey string `yaml:"accesskey"`
	SecureKey string `yaml:"securekey"`
}

func getConfig() (Config, error) {
	conf := Config{}
	filename, err := filepath.Abs("./config.yml")
	if err != nil {
		log.Error(err)
		return conf, err
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return conf, err
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
	return conf, nil
}