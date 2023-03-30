package s3info

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type S3Info struct {
	AwsS3Region    string
	AwsAccessKey   string
	AwsSecretKey   string
	AwsProfileName string
	BucketName     string
	EndPoint       string //http://192.168.57.11:80  aws.String() aws.Endpoint
	S3Client       *s3.Client
}

func NewS3Info(accesskey, securekey, endpointurl string) (*S3Info, error) {
	s3info := S3Info{
		AwsS3Region:  "local",
		AwsAccessKey: accesskey,
		AwsSecretKey: securekey,
		EndPoint:     endpointurl,
	}
	return &s3info, nil
}

// key를 활용해서 Client 생성
func (s *S3Info) SetConfigByKeyEndpoint() error {
	credentials := credentials.NewStaticCredentialsProvider(s.AwsAccessKey, s.AwsSecretKey, "")
	conf, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials),
		config.WithRegion("local"),
	)
	if err != nil {
		log.Error(err)
		return err
	}

	// Create a new S3 SDK client instance.
	s3Client := s3.NewFromConfig(
		conf,
		s3.WithEndpointResolver(
			s3.EndpointResolverFromURL(s.EndPoint),
		),
		func(opts *s3.Options) {
			opts.UsePathStyle = true
		},
	)

	if err != nil {
		log.Printf("error: %v", err)
		//panic(err)
		return errors.New(err.Error())
	}
	s.S3Client = s3Client
	return nil
}

func (s *S3Info) SetConfigByDefault() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s.AwsS3Region))
	if err != nil {
		log.Fatal(err)
		return errors.New(err.Error())
	}
	s.S3Client = s3.NewFromConfig(cfg)
	return nil
}

// profile Name을 활용해서 Client 생성
func (s *S3Info) SetConfigByProfile() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s.AwsS3Region),
		config.WithSharedConfigProfile(s.AwsProfileName))
	if err != nil {
		log.Fatal(err)
		return errors.New(err.Error())
	}
	s.S3Client = s3.NewFromConfig(cfg)
	return nil
}

// key를 활용해서 Client 생성
func (s *S3Info) SetConfigByKey() error {
	creds := credentials.NewStaticCredentialsProvider(s.AwsAccessKey, s.AwsSecretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(s.AwsS3Region),
	)
	if err != nil {
		log.Printf("error: %v", err)
		//panic(err)
		return errors.New(err.Error())
	}
	s.S3Client = s3.NewFromConfig(cfg)
	return nil
}

// 서버를 통해 파일을 받아왔을 때 사용
func (s *S3Info) UploadFile(file io.Reader, filename, preFix string) *manager.UploadOutput {
	uploader := manager.NewUploader(s.S3Client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return result
}

// 파일이름을 통해 파일을 불러와 서버에 업로드
func (s *S3Info) UploadFileByFileName(originalFilename, fileName, preFix string) *manager.UploadOutput {
	file, err := ioutil.ReadFile(originalFilename)
	if err != nil {
		panic(err)
	}
	uploader := manager.NewUploader(s.S3Client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(preFix + "/" + fileName),
		Body:   bytes.NewReader(file),
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return result
}

// 원하는 파일을 다운로드 받을 때 사용
func (s *S3Info) DownloadFile(targetDirectory, key string) error {
	// Create the directories in the path
	splitKeyArr := strings.Split(key, "/")
	file := filepath.Join(targetDirectory, splitKeyArr[len(splitKeyArr)-1])
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return err
	}

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	downloader := manager.NewDownloader(s.S3Client)
	_, err = downloader.Download(context.TODO(), fd,
		&s3.GetObjectInput{
			Bucket: &s.BucketName,
			Key:    &key,
		})
	return err
}

// 원하는 파일을 다운로드 받을 때 사용
func (s *S3Info) DownloadFile2(targetDirectory, key string) error {
	// Create the directories in the path
	splitKeyArr := strings.Split(key, "/")
	file := filepath.Join(targetDirectory, splitKeyArr[len(splitKeyArr)-1])
	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	downloader := manager.NewDownloader(s.S3Client)
	_, err = downloader.Download(context.TODO(), fd,
		&s3.GetObjectInput{
			Bucket: &s.BucketName,
			Key:    &key,
		})
	return err
}

// 버킷안에 있는 Objects 확인할때 사용
func (s *S3Info) GetItems(prefix string) {
	paginator := s3.NewListObjectsV2Paginator(s.S3Client,
		&s3.ListObjectsV2Input{
			Bucket: &s.BucketName,
			Prefix: aws.String(prefix),
		})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalln("error:", err)
		}
		for _, obj := range page.Contents {

			fmt.Println(aws.ToString(obj.Key))
		}
	}
}

// 버킷 리스트 확인
func (s *S3Info) GetBucketList() {
	output, err := s.S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	for _, bucket := range output.Buckets {
		log.Debugf("[%+v]", bucket)
		fmt.Println(*bucket.Name)
	}
}

// 버킷 생성하기
func (s *S3Info) CreateBucket(bucketName string, region types.BucketLocationConstraint) error {
	// output, err := s.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
	// 	Bucket: &bucketName,
	// 	CreateBucketConfiguration: &types.CreateBucketConfiguration{
	// 		LocationConstraint: region,
	// 	},
	// })
	output, err := s.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Println(output.Location)
	return nil
}
