package main

import (
	"ceph-s3/aws-sdk-go-v2-s3info/s3info"
	"testing"
)

var AwsAccessKey = "XMH5L0E9LMAX44PAR66G"
var AwsSecretKey = "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG"
var EndPointUrl = "http://192.168.57.11"

func TestListBuckets(t *testing.T) {
	s3info, err := s3info.NewS3Info(AwsAccessKey, AwsSecretKey, EndPointUrl)
	if err != nil {
		t.Error(err)
	}
	s3info.SetConfigByKeyEndpoint()
	s3info.GetBucketList()
	s3info.CreateBucket("newbucket1", "")
	//	return s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &bucketName})
	s3info.BucketName = "mybucket"
	s3info.GetItems("")
}
