package adm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ceph/go-ceph/rgw/admin"
	tsuite "github.com/stretchr/testify/suite"
)

var (
	errNoEndpoint  = errors.New("endpoint not set")
	errNoAccessKey = errors.New("access key not set")
	errNoSecretKey = errors.New("secret key not set")
)

var (
	errMissingUserID          = errors.New("missing user ID")
	errMissingSubuserID       = errors.New("missing subuser ID")
	errMissingUserAccessKey   = errors.New("missing user access key")
	errMissingUserDisplayName = errors.New("missing user display name")
	errMissingUserCap         = errors.New("missing user capabilities")
	errMissingBucketID        = errors.New("missing bucket ID")
	errMissingBucket          = errors.New("missing bucket")
	errMissingUserBucket      = errors.New("missing bucket")
	errUnsupportedKeyType     = errors.New("unsupported key type")
)

func TestRadosGWTestSuite(t *testing.T) {
	tsuite.Run(t, new(RadosGWTestSuite))
}

type RadosGWTestSuite struct {
	tsuite.Suite
	endpoint       string
	accessKey      string
	secretKey      string
	bucketTestName string
}

type debugHTTPClient struct {
	client admin.HTTPClient
}

func newDebugHTTPClient(client admin.HTTPClient) *debugHTTPClient {
	return &debugHTTPClient{client}
}

func (c *debugHTTPClient) Do(req *http.Request) (*http.Response, error) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n%s\n", string(dump))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n%s\n", string(dump))

	return resp, nil
}

// S3Agent wraps the s3.S3 structure to allow for wrapper methods
type S3Agent struct {
	Client *s3.S3
}

func newS3Agent(accessKey, secretKey, endpoint string, debug bool) (*S3Agent, error) {
	const cephRegion = "us-east-1"

	logLevel := aws.LogOff
	if debug {
		logLevel = aws.LogDebug
	}
	client := http.Client{
		Timeout: time.Second * 15,
	}
	sess, err := session.NewSession(
		aws.NewConfig().
			WithRegion(cephRegion).
			WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")).
			WithEndpoint(endpoint).
			WithS3ForcePathStyle(true).
			WithMaxRetries(5).
			WithDisableSSL(true).
			WithHTTPClient(&client).
			WithLogLevel(logLevel),
	)
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	return &S3Agent{
		Client: svc,
	}, nil
}

func (s *S3Agent) createBucket(name string) error {
	bucketInput := &s3.CreateBucketInput{
		Bucket: &name,
	}
	_, err := s.Client.CreateBucket(bucketInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				return nil
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				return nil
			}
		}
		return fmt.Errorf("failed to create bucket %q. %w", name, err)
	}
	return nil
}

func (suite *RadosGWTestSuite) SetupConnection() {
	suite.accessKey = "XMH5L0E9LMAX44PAR66G"
	suite.secretKey = "rCKebrRAwKAjanRiwzJEM2oJkmPzOXjmLEj9RaiG"
	hostname := os.Getenv("HOSTNAME")
	endpoint := "192.168.57.11"
	if hostname != "test_ceph_aio" {
		endpoint = "192.168.57.11"
	}
	suite.endpoint = "http://" + endpoint
	suite.bucketTestName = "test-bucket"
}

func TestNew(t *testing.T) {
	type args struct {
		endpoint  string
		accessKey string
		secretKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *admin.API
		wantErr error
	}{
		{"no endpoint", args{}, nil, errNoEndpoint},
		{"no accessKey", args{endpoint: "http://192.168.0.1"}, nil, errNoAccessKey},
		{"no secretKey", args{endpoint: "http://192.168.0.1", accessKey: "foo"}, nil, errNoSecretKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := admin.New(tt.args.endpoint, tt.args.accessKey, tt.args.secretKey, nil)
			if tt.wantErr != err {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *RadosGWTestSuite) TestBucket() {
	suite.SetupConnection()
	co, err := admin.New(suite.endpoint, suite.accessKey, suite.secretKey, newDebugHTTPClient(http.DefaultClient))
	assert.NoError(suite.T(), err)

	s3, err := newS3Agent(suite.accessKey, suite.secretKey, suite.endpoint, true)
	assert.NoError(suite.T(), err)

	err = s3.createBucket(suite.bucketTestName)
	assert.NoError(suite.T(), err)

	suite.T().Run("list buckets", func(t *testing.T) {
		buckets, err := co.ListBuckets(context.Background())
		assert.NoError(suite.T(), err)
		fmt.Print( len(buckets))
	})

	suite.T().Run("info non-existing bucket", func(t *testing.T) {
		_, err := co.GetBucketInfo(context.Background(), admin.Bucket{Bucket: "foo"})
		assert.Error(suite.T(), err)
		assert.True(suite.T(), errors.Is(err, admin.ErrNoSuchBucket), err)
	})

	suite.T().Run("info existing bucket", func(t *testing.T) {
		_, err := co.GetBucketInfo(context.Background(), admin.Bucket{Bucket: suite.bucketTestName})
		assert.NoError(suite.T(), err)
	})

	suite.T().Run("remove bucket", func(t *testing.T) {
		err := co.RemoveBucket(context.Background(), admin.Bucket{Bucket: suite.bucketTestName})
		assert.NoError(suite.T(), err)
	})

	suite.T().Run("list bucket is now zero", func(t *testing.T) {
		buckets, err := co.ListBuckets(context.Background())
		assert.NoError(suite.T(), err)
		fmt.Print(len(buckets))
	})

	suite.T().Run("remove non-existing bucket", func(t *testing.T) {
		err := co.RemoveBucket(context.Background(), admin.Bucket{Bucket: "foo"})
		assert.Error(suite.T(), err)
	})
}

var testBucketQuota = 1000000

func (suite *RadosGWTestSuite) TestBucketQuota() {
	suite.SetupConnection()
	co, err := admin.New(suite.endpoint, suite.accessKey, suite.secretKey, newDebugHTTPClient(http.DefaultClient))
	assert.NoError(suite.T(), err)

	s3, err := newS3Agent(suite.accessKey, suite.secretKey, suite.endpoint, true)
	assert.NoError(suite.T(), err)

	err = s3.createBucket(suite.bucketTestName)
	assert.NoError(suite.T(), err)

	suite.T().Run("set bucket quota but no user is specified", func(t *testing.T) {
		err := co.SetIndividualBucketQuota(context.Background(), admin.QuotaSpec{})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserID.Error())
	})

	suite.T().Run("set bucket quota but no bucket is specified", func(t *testing.T) {
		err := co.SetIndividualBucketQuota(context.Background(), admin.QuotaSpec{UID: "admin"})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserBucket.Error())
	})

	suite.T().Run("set bucket quota", func(t *testing.T) {
		err := co.SetIndividualBucketQuota(context.Background(), admin.QuotaSpec{UID: "admin", Bucket: suite.bucketTestName, MaxSizeKb: &testBucketQuota})
		assert.NoError(suite.T(), err)

		bucketInfo, err := co.GetBucketInfo(context.Background(), admin.Bucket{Bucket: suite.bucketTestName})
		assert.NoError(suite.T(), err)

		assert.Equal(suite.T(), &testBucketQuota, bucketInfo.BucketQuota.MaxSizeKb)
	})
}

func (suite *RadosGWTestSuite) TestQuota() {
	suite.SetupConnection()
	co, err := admin.New(suite.endpoint, suite.accessKey, suite.secretKey, newDebugHTTPClient(http.DefaultClient))
	assert.NoError(suite.T(), err)

	suite.T().Run("set quota to user but user ID is empty", func(t *testing.T) {
		err := co.SetUserQuota(context.Background(), admin.QuotaSpec{})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserID.Error())
	})

	suite.T().Run("get user quota but no user is specified", func(t *testing.T) {
		_, err := co.GetUserQuota(context.Background(), admin.QuotaSpec{})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserID.Error())

	})
}

func (suite *RadosGWTestSuite) TestUsage() {
	suite.SetupConnection()
	co, err := admin.New(suite.endpoint, suite.accessKey, suite.secretKey, newDebugHTTPClient(http.DefaultClient))
	assert.NoError(suite.T(), err)

	suite.T().Run("get usage", func(t *testing.T) {
		pTrue := true
		usage, err := co.GetUsage(context.Background(), admin.Usage{ShowSummary: &pTrue})
		assert.NoError(suite.T(), err)
		assert.NotEmpty(t, usage)
	})

	suite.T().Run("trim usage", func(t *testing.T) {
		pFalse := false
		_, err := co.GetUsage(context.Background(), admin.Usage{RemoveAll: &pFalse})
		assert.NoError(suite.T(), err)
	})
}

const (
	testAk = "HDNEZQXZAA6NIWOBOL0U"
)

func (suite *RadosGWTestSuite) TestKeys() {
	suite.SetupConnection()
	co, err := admin.New(suite.endpoint, suite.accessKey, suite.secretKey, newDebugHTTPClient(http.DefaultClient))
	assert.NoError(suite.T(), err)

	var keys *[]admin.UserKeySpec

	suite.T().Run("create keys but user ID and SubUser is empty", func(t *testing.T) {
		_, err := co.CreateKey(context.Background(), admin.UserKeySpec{})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserID.Error())
	})

	suite.T().Run("create swift keys but SubUser is empty", func(t *testing.T) {
		_, err := co.CreateKey(context.Background(), admin.UserKeySpec{KeyType: "swift"})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingSubuserID.Error())
	})

	suite.T().Run("create some unknown key type", func(t *testing.T) {
		_, err := co.CreateKey(context.Background(), admin.UserKeySpec{KeyType: "some-key-type"})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errUnsupportedKeyType.Error())
	})

	suite.T().Run("create keys without ak or sk provided", func(t *testing.T) {
		keys, err = co.CreateKey(context.Background(), admin.UserKeySpec{UID: "myuser1"})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 2, len(*keys))
	})

	suite.T().Run("create keys with ak provided", func(t *testing.T) {
		keys, err = co.CreateKey(context.Background(), admin.UserKeySpec{UID: "myuser1", AccessKey: testAk})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 3, len(*keys))
	})

	suite.T().Run("remove keys but user ID and SubUser is empty", func(t *testing.T) {
		err := co.RemoveKey(context.Background(), admin.UserKeySpec{})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserID.Error())
	})

	suite.T().Run("remove s3 keys but ak is empty", func(t *testing.T) {
		err := co.RemoveKey(context.Background(), admin.UserKeySpec{UID: "myuser1"})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserAccessKey.Error())
	})

	suite.T().Run("remove s3 key", func(t *testing.T) {
		for _, key := range *keys {
			if key.AccessKey != suite.accessKey {
				err := co.RemoveKey(context.Background(), admin.UserKeySpec{UID: "myuser1", AccessKey: key.AccessKey})
				assert.NoError(suite.T(), err)
			}
		}
	})
}

// mockClient is the mock of the HTTP Client
// It can be used to mock HTTP request/response from the rgw admin ops API
type mockClient struct {
	// mockDo is a type that mock the Do method from the HTTP package
	mockDo mockDoType
}

// mockDoType is a custom type that allows setting the function that our Mock Do func will run instead
type mockDoType func(req *http.Request) (*http.Response, error)

// Do is the mock client's `Do` func
func (m *mockClient) Do(req *http.Request) (*http.Response, error) { return m.mockDo(req) }

var (
	fakeUserResponse = []byte(`
{
  "tenant": "",
  "user_id": "dashboard-admin",
  "display_name": "dashboard-admin",
  "email": "",
  "suspended": 0,
  "max_buckets": 1000,
  "subusers": [
     {
        "id": "dashboard-admin:swift",
        "permissions": "read"
     }
  ],
  "keys": [
    {
      "user": "dashboard-admin",
      "access_key": "4WD1FGM5PXKLC97YC0SZ",
      "secret_key": "YSaT5bEcJTjBJCDG5yvr2NhGQ9xzoTIg8B1gQHa3"
    }
  ],
  "swift_keys": [
    {
      "user": "dashboard-admin:swift",
      "secret_key": "VERY_SECRET"
    }
  ],
  "caps": [],
  "op_mask": "read, write, delete",
  "system": "true",
  "admin": "false",
  "default_placement": "",
  "default_storage_class": "",
  "placement_tags": [],
  "bucket_quota": {
    "enabled": false,
    "check_on_raw": false,
    "max_size": -1,
    "max_size_kb": 0,
    "max_objects": -1
  },
  "user_quota": {
    "enabled": false,
    "check_on_raw": false,
    "max_size": -1,
    "max_size_kb": 0,
    "max_objects": -1
  },
  "temp_url_keys": [],
  "type": "rgw",
  "mfa_ids": []
}`)
)

func TestUnmarshal_1(t *testing.T) {
	u := &admin.User{}
	err := json.Unmarshal(fakeUserResponse, &u)
	assert.NoError(t, err)
}

func (suite *RadosGWTestSuite) TestUser_1() {
	suite.SetupConnection()
	co, err := admin.New(suite.endpoint, suite.accessKey, suite.secretKey, newDebugHTTPClient(http.DefaultClient))
	assert.NoError(suite.T(), err)

	suite.T().Run("fail to create user since no UID provided", func(t *testing.T) {
		_, err = co.CreateUser(context.Background(), admin.User{Email: "leseb@example.com"})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserID.Error())
	})

	suite.T().Run("fail to create user since no no display name provided", func(t *testing.T) {
		_, err = co.CreateUser(context.Background(), admin.User{ID: "leseb", Email: "leseb@example.com"})
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, errMissingUserDisplayName.Error())
	})

	suite.T().Run("user creation success", func(t *testing.T) {
		usercaps := "users=read"
		user, err := co.CreateUser(context.Background(), admin.User{ID: "leseb", DisplayName: "This is leseb", Email: "leseb@example.com", UserCaps: usercaps})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "leseb@example.com", user.Email)
	})

	suite.T().Run("get user leseb by uid", func(t *testing.T) {
		user, err := co.GetUser(context.Background(), admin.User{ID: "leseb"})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "leseb@example.com", user.Email)
		assert.Equal(suite.T(), "users", user.Caps[0].Type)
		assert.Equal(suite.T(), "read", user.Caps[0].Perm)
		os.Setenv("LESEB_ACCESS_KEY", user.Keys[0].AccessKey)
	})

	suite.T().Run("get user leseb by key", func(t *testing.T) {
		user, err := co.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: os.Getenv("LESEB_ACCESS_KEY")}}})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "leseb@example.com", user.Email)
		os.Unsetenv("LESEB_ACCESS_KEY")
	})

	suite.T().Run("modify user email", func(t *testing.T) {
		user, err := co.ModifyUser(context.Background(), admin.User{ID: "leseb", Email: "leseb@leseb.com"})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "leseb@leseb.com", user.Email)
	})

	suite.T().Run("modify user max bucket", func(t *testing.T) {
		maxBuckets := -1
		user, err := co.ModifyUser(context.Background(), admin.User{ID: "leseb", MaxBuckets: &maxBuckets})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "leseb@leseb.com", user.Email)
		assert.Equal(suite.T(), -1, *user.MaxBuckets)
	})

	suite.T().Run("user already exists", func(t *testing.T) {
		_, err := co.CreateUser(context.Background(), admin.User{ID: "myuser1", DisplayName: "Admin user"})
		assert.Error(suite.T(), err)
		fmt.Sprintf("%+v", err)
	})

	suite.T().Run("get users", func(t *testing.T) {
		users, err := co.GetUsers(context.Background())
		assert.NoError(suite.T(), err)
		fmt.Printf("%+v",*users)
	})

	suite.T().Run("set user quota", func(t *testing.T) {
		quotaEnable := true
		maxObjects := int64(100)
		err := co.SetUserQuota(context.Background(), admin.QuotaSpec{QuotaType: "user", UID: "leseb", MaxObjects: &maxObjects, Enabled: &quotaEnable})
		assert.NoError(suite.T(), err)
	})

	suite.T().Run("get user quota", func(t *testing.T) {
		q, err := co.GetUserQuota(context.Background(), admin.QuotaSpec{QuotaType: "user", UID: "leseb"})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(100), *q.MaxObjects)
	})

	suite.T().Run("get user stat", func(t *testing.T) {
		statEnable := true
		user, err := co.GetUser(context.Background(), admin.User{ID: "leseb", GenerateStat: &statEnable})
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), user.Stat.Size)
	})

	suite.T().Run("create a subuser", func(t *testing.T) {
		err := co.CreateSubuser(context.Background(), admin.User{ID: "leseb"}, admin.SubuserSpec{Name: "foo", Access: admin.SubuserAccessReadWrite})
		assert.NoError(suite.T(), err)

		user, err := co.GetUser(context.Background(), admin.User{ID: "leseb"})
		assert.NoError(suite.T(), err)
		if err == nil {
			assert.Equal(suite.T(), user.Subusers[0].Name, "leseb:foo")
			// Note: the returned values are not equal to the input values ...
			assert.Equal(suite.T(), user.Subusers[0].Access, admin.SubuserAccess("read-write"))
		}
	})

	suite.T().Run("modify a subuser", func(t *testing.T) {
		err := co.ModifySubuser(context.Background(), admin.User{ID: "leseb"}, admin.SubuserSpec{Name: "foo", Access: admin.SubuserAccessRead})
		assert.NoError(suite.T(), err)

		user, err := co.GetUser(context.Background(), admin.User{ID: "leseb"})
		assert.NoError(suite.T(), err)
		if err == nil {
			assert.Equal(suite.T(), user.Subusers[0].Name, "leseb:foo")
			assert.Equal(suite.T(), user.Subusers[0].Access, admin.SubuserAccess("read"))
		}
	})

	suite.T().Run("remove a subuser", func(t *testing.T) {
		err := co.RemoveSubuser(context.Background(), admin.User{ID: "leseb"}, admin.SubuserSpec{Name: "foo"})
		assert.NoError(suite.T(), err)

		user, err := co.GetUser(context.Background(), admin.User{ID: "leseb"})
		assert.NoError(suite.T(), err)
		if err == nil {
			assert.Equal(suite.T(), len(user.Subusers), 0)
		}
	})

	suite.T().Run("remove user", func(t *testing.T) {
		err = co.RemoveUser(context.Background(), admin.User{ID: "leseb"})
		assert.NoError(suite.T(), err)
	})
}

func TestGetUserMockAPI(t *testing.T) {
	t.Run("test simple api mock", func(t *testing.T) {
		api, err := admin.New("127.0.0.1", "accessKey", "secretKey", returnMockClient())
		assert.NoError(t, err)
		u, err := api.GetUser(context.TODO(), admin.User{ID: "dashboard-admin"})
		assert.NoError(t, err)
		assert.Equal(t, "dashboard-admin", u.DisplayName, u)
	})
	t.Run("test get user with access key", func(t *testing.T) {
		api, err := admin.New("127.0.0.1", "accessKey", "secretKey", returnMockClient())
		assert.NoError(t, err)
		u, err := api.GetUser(context.TODO(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: "4WD1FGM5PXKLC97YC0SZ"}}})
		assert.NoError(t, err)
		assert.Equal(t, "dashboard-admin", u.DisplayName, u)
	})
	t.Run("test get user with nothing", func(t *testing.T) {
		api, err := admin.New("127.0.0.1", "accessKey", "secretKey", returnMockClient())
		assert.NoError(t, err)
		_, err = api.GetUser(context.TODO(), admin.User{})
		assert.Error(t, err)
		assert.EqualError(t, err, "missing user ID")
	})
	t.Run("test get user with missing correct key", func(t *testing.T) {
		api, err := admin.New("127.0.0.1", "accessKey", "secretKey", returnMockClient())
		assert.NoError(t, err)
		_, err = api.GetUser(context.TODO(), admin.User{Keys: []admin.UserKeySpec{{SecretKey: "4WD1FGM5PXKLC97YC0SZ"}}})
		assert.Error(t, err)
		assert.EqualError(t, err, "missing user access key")
	})
}

func returnMockClient() *mockClient {
	r := ioutil.NopCloser(bytes.NewReader(fakeUserResponse))
	return &mockClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodGet && req.URL.Path == "127.0.0.1/admin/user" {
				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			}
			return nil, fmt.Errorf("unexpected request: %q. method %q. path %q", req.URL.RawQuery, req.Method, req.URL.Path)
		},
	}
}
