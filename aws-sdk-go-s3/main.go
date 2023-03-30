package main

import (
	"crypto/tls"
	"fmt"	
	"net/http"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"	
)


func main() {
	endpointurl := "http://192.168.57.11"
	accesskey := "XMH5L0E9LMAX44PAR66G"
	securekey := "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG"
	//ctx := context.TODO()
	//ctx := context.Background()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String("local"),
		Credentials:      credentials.NewStaticCredentials(accesskey,securekey, ""),
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
	// snippet-end:[s3.go.list_buckets.imports.print]
}

// snippet-end:[s3.go.list_buckets]
func GetAllBuckets(sess *session.Session) (*s3.ListBucketsOutput, error) {
    // snippet-start:[s3.go.list_buckets.imports.call]
    svc := s3.New(sess)

    result, err := svc.ListBuckets(&s3.ListBucketsInput{})
    // snippet-end:[s3.go.list_buckets.imports.call]
    if err != nil {
        return nil, err
    }
    return result, nil
}