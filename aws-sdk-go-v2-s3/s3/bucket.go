package s3

// Include the necessary AWS SDK modules.
import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

// Required bucket name.
const localBucketName = "records"

// Function declaration: Getting a unique bucket name based on <accessKeyId> value
func getUniqueBucketName(objectStorageName string, localBucketName string) (*string, error) {
	// Necessary environment variable names.
	const accessKeyId = "accessKeyId"
	// All bucket names in the Zerops shared object storage namespace have to be unique!
	// Getting the environment variable value that will be used as the unique prefix.
	value, found := os.LookupEnv(objectStorageName + "_" + accessKeyId)
	if !found {
		return nil, errors.New("non-existed accessKeyId environment variable")
	}
	// Unique bucket name preparation.
	bucketName := strings.ToLower(value + "." + localBucketName)
	return &bucketName, nil
}

// Function declaration: Getting an S3 SDK client
func createBucket(ctx context.Context, s3Client *s3.Client, localBucketName string) (*s3.CreateBucketOutput, error) {
	bucketName := localBucketName
	return s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &bucketName})
}

func listBuckets(ctx context.Context, s3Client *s3.Client) (*s3.ListBucketsOutput, error) {
	return s3Client.ListBuckets(ctx, nil)
}

// Function declaration: Getting an exited bucket ACL
func getBucketAcl(ctx context.Context, s3Client *s3.Client, localBucketName string) (*s3.GetBucketAclOutput, error) {
	bucketName := localBucketName
	_, err := s3Client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		return nil, err
	}
	return s3Client.GetBucketAcl(ctx, &s3.GetBucketAclInput{Bucket: &bucketName})
}

func Listbucket() error {
	conf, err := getConfig()
	if err != nil {
		log.Error(err)
		return err
	}
	//ctx := context.TODO()
	ctx := context.Background()
	// Calling the function: getCredentials
	storeCredentials := getCredentials(conf)
	if storeCredentials != nil {
		s3Client, err := getS3Client(ctx, conf.Url, storeCredentials)
		if err != nil {
			log.Error(err)
			return err
		}
		if s3Client != nil {
			// Calling the function: createBucket
			listBucketsOutput, err := listBuckets(ctx, s3Client)
			if err != nil {
				log.Error(err)
				return err
			}
			// Printing the existed bucket names
			for _, bucket := range listBucketsOutput.Buckets {
				fmt.Println(*bucket.Name)
			}
		}
	}
	return nil
}

func Createbucket(name string) error {
	//ctx := context.TODO()
	ctx := context.Background()
	conf, err := getConfig()
	if err != nil {
		log.Error(err)
		return err
	}
	// Calling the function: getCredentials
	storeCredentials := getCredentials(conf)
	if storeCredentials != nil {
		// Calling the function: getS3Client
		s3Client, err := getS3Client(ctx, conf.Url, storeCredentials)
		if err != nil {
			log.Error(err)
			return err
		}
		if s3Client != nil {
			// Calling the function: createBucket
			createBucketOutput, err := createBucket(ctx, s3Client, name)
			// Calling the function: listBuckets
			if err != nil {
				log.Error(err)
			}
			log.Debugf("%+v", createBucketOutput)
		}
	}
	return nil
}

func Getbucketacl(name string) error {
	ctx := context.Background()
	conf, err := getConfig()
	if err != nil {
		log.Error(err)
		return err
	}
	storeCredentials := getCredentials(conf)
	if storeCredentials != nil {
		s3Client, err := getS3Client(ctx, conf.Url, storeCredentials)
		if err != nil {
			log.Error(err)
			return err
		}
		if s3Client != nil {
			bucketAclOutput, err := getBucketAcl(ctx, s3Client, name)
			if err != nil {
				log.Error(err)
			}
			log.Debugf("%+v", bucketAclOutput)
			// if err == nil {
			// 	// Formatting and printing the returned value of the bucket ACL.
			// 	grantsJSON, err := json.MarshalIndent(bucketAclOutput.Grants, "", "  ")
			// 	fmt.Printf("%s\n", string(grantsJSON))
			// }
		}
	}
	return nil
}
