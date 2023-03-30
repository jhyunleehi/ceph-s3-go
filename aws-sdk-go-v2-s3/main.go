package main

import "aws-sdk-go-v2-s3/s3"

func main() {
	s3.Listbucket()
	s3.Createbucket("test")
	s3.Getbucketacl("test")

}
