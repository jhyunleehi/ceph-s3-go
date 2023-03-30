package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

func main() {
	endpointurl := "http://192.168.57.11"
	accesskey := "XMH5L0E9LMAX44PAR66G"
	securekey := "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG"
	//ctx := context.TODO()
	ctx := context.Background()

	credentials := credentials.NewStaticCredentialsProvider(accesskey, securekey, "")
	cnf, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(credentials),
		config.WithRegion("local"),
	)
	if err != nil {
		log.Error(err)
		return
	}
	// Create a new S3 SDK client instance.
	s3Client := s3.NewFromConfig(
		cnf,
		s3.WithEndpointResolver(
			s3.EndpointResolverFromURL(endpointurl),
		),
		func(opts *s3.Options) {
			opts.UsePathStyle = true
		},
	)
	if s3Client != nil {		
		listBucketsOutput, err := s3Client.ListBuckets(ctx, nil)
		if err != nil {
			log.Error(err)
			return
		}		
		for _, bucket := range listBucketsOutput.Buckets {
			fmt.Println(*bucket.Name)
		}
	}
}
