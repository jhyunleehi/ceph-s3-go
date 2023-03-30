package main

import (
	"aws-sdk-go-v2-s3info/s3info"

	log "github.com/sirupsen/logrus"
)

func main() {
	AwsAccessKey := "XMH5L0E9LMAX44PAR66G"
	AwsSecretKey := "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG"
	EndPointUrl := "http://192.168.57.11"

	s3info, err := s3info.NewS3Info(AwsAccessKey, AwsSecretKey, EndPointUrl)
	if err != nil {
		log.Error(err)
	}
	s3info.SetConfigByKeyEndpoint()
	s3info.GetBucketList()
	s3info.CreateBucket("newbucket1","")
		
	s3info.BucketName = "mybucket"
	s3info.GetItems("")
	
}
