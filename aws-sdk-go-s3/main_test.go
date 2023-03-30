package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var endpointurl = "http://192.168.57.11"
var accesskey = "XMH5L0E9LMAX44PAR66G"
var securekey = "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG"

func TestListBuckets(t *testing.T) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String("local"),
		Credentials:      credentials.NewStaticCredentials(accesskey, securekey, ""),
		Endpoint:         aws.String(endpointurl),
		HTTPClient:       client,
		S3ForcePathStyle: aws.Bool(true),
	}),
	)

	result, err := GetAllBuckets(sess)
	if err != nil {
		fmt.Println("Got an error retrieving buckets:")
		fmt.Println(err)
		return
	}

	// snippet-start:[s3.go.list_buckets.imports.print]
	fmt.Println("Buckets:")

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name + ": " + bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}
}
