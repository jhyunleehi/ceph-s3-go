package main

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Url       string `yaml:"url"`
	AccessKey string `yaml:"accesskey"`
	SecureKey string `yaml:"securekey"`
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
	conf := Config{
		Url:       "http://192.168.57.11",
		AccessKey: "XMH5L0E9LMAX44PAR66G",
		SecureKey: "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG",
	}
	co, err := NewCephS3(conf)
	if err != nil {
		log.Error(err)
		return
	}
	user, err := co.GetUser(context.Background(), admin.User{ID: "myuser1"})
	if err != nil {
		log.Errorf("%+v", err)
	}
	// Print the user display name
	fmt.Println(user.DisplayName)
	fmt.Printf("%+v", user)

}
