package main

import (
	"context"
	"fmt"

	"io/ioutil"

	"path/filepath"

	"github.com/ceph/go-ceph/rgw/admin"
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

func NewCephS3(c Config) (*admin.API, error) {
	svc, err := admin.New(c.Url, c.AccessKey, c.SecureKey, nil)
	if err != nil {
		panic(err)
	}
	log.Debugf("created session successfully [%+v] ", svc)
	return svc, nil
}

func main() {
	conf, err := getConfig()
	if err != nil {
		log.Error(err)
		return
	}
	co, err := NewCephS3(conf)
	if err != nil {
		log.Error(err)
		return
	}
	user, err := co.GetUser(context.Background(), admin.User{ID: "admin"})
	if err != nil {
		log.Errorf("%+v", err)
	}
	// Print the user display name
	fmt.Println(user.DisplayName)
	fmt.Printf("%+v", user)

}
