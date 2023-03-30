package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/ceph/go-ceph/rgw/admin"

	log "github.com/sirupsen/logrus"
)

func TestListBuckets(t *testing.T) {
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
	rtn, err := co.ListBuckets(context.Background())
	if err != nil {
		log.Errorf("%+v", err)
	}
	for i, b := range rtn {
		fmt.Printf("[%d][%s]\n", i, b)
	}
	fmt.Printf("%+v", rtn)

	user, err := co.GetUser(context.Background(), admin.User{ID: "admin"})
	if err != nil {
		log.Errorf("%+v", err)
	}
	// Print the user display name
	fmt.Println(user.DisplayName)
	fmt.Printf("%+v", user)

}
