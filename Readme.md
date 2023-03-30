# Ceph RGW Storage S3 API

## AWS SDK GO S3 모듈

* 기존 `github.com/aws/aws-sdk-go/aws` 라이브러리를 이용
* Ceph 스토리지 사용을 위해서는 스토리지에서 제공하는 EndpointUrl 설정이 필요함  
* 참조: <https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go/s3>

```go
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
```

* 실행

```log
C:\Gocode\src\ceph-s3-go\aws-sdk-go-s3>go run -mod vendor main.go
```

## AWS SDK GO V2 S3 모듈

* github.com/aws/aws-sdk-go-v2/service/s3  
* Ceph 스토리지 사용을 위해서는 스토리지에서 제공하는 EndpointUrl 설정이 필요함  
* 참조: <https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2/s3>

```go
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

```
* 실행
```log
C:\Gocode\src\ceph-s3-go\aws-sdk-go-v2>go run -mod vendor main.go
```

## go ceph admin

* github.com/ceph/go-ceph/rgw/admin
* admin portal, rdb, fs 등 다양한 접근 방법 제공
* 여기서는 rgw/admin 이용한 s3 API 테스트

```go
package main

import (
 "context"
 "fmt"

 "io/ioutil"

 "path/filepath"

 "github.com/ceph/go-ceph/rgw/admin"
 "gopkg.in/yaml.v2"

 log "github.com/sirupsen/logrus"
)

type Config struct {
 Url       string `yaml:"url"`
 AccessKey string `yaml:"accesskey"`
 SecureKey string `yaml:"securekey"`
}

func getConfig() (Config, error) {
 conf := Config{}
 filename, err := filepath.Abs("./config.yml")
 if err != nil {
  log.Error(err)
  return conf, err
 }
 yamlFile, err := ioutil.ReadFile(filename)
 if err != nil {
  log.Error(err)
  return conf, err
 }
 err = yaml.Unmarshal(yamlFile, &conf)
 if err != nil {
  panic(err)
 }
 return conf, nil
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
 user, err := co.GetUser(context.Background(), admin.User{ID: "admin"})
 if err != nil {
  log.Errorf("%+v", err)
 }
 // Print the user display name
 fmt.Println(user.DisplayName)
 fmt.Printf("%+v", user)

}
````

* 실행방법

```log
C:\Gocode\src\ceph-s3-go\go-ceph-admin>go run -mod vendor main.go
```
