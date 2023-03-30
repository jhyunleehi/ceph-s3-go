package s3

// Include the necessary AWS SDK modules.
import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var ctx = context.Background()

// Object storage name.
const storeObjectStorageName = "store"

// Function returning a pointer to a variable with an accessKeyId value.
func getAccessKeyIdValue(objectStorageName string) *string {
	// Necessary environment variable name.
	const accessKeyId = "accessKeyId"
	value, found := os.LookupEnv(objectStorageName + "_" + accessKeyId)
	// If the environment variable has been found ...
	if found {
		return &value
	}
	return nil
}

// Function returning a pointer to a variable with a secretAccessKey value.
func getSecretAccessKeyValue(objectStorageName string) *string {
	// Necessary environment variable name.
	const secretAccessKey = "secretAccessKey"
	value, found := os.LookupEnv(objectStorageName + "_" + secretAccessKey)
	// If the environment variable has been found ...
	if found {
		return &value
	}
	return nil
}

// Function returning a pointer to a variable with a user credentials.
func getCredentials(conf Config ) *credentials.StaticCredentialsProvider {
	accessKeyIdValue := conf.AccessKey
	secretAccessKeyValue := conf.SecureKey
	if accessKeyIdValue != "" && secretAccessKeyValue !="" {
		credentials := credentials.NewStaticCredentialsProvider(accessKeyIdValue, secretAccessKeyValue, "")
		// A pointer to a variable with credentials is returned.
		return &credentials
	}
	// If any environment variable not found, return only the nil (it's compatible with a pointer).
	return nil
}

// Function declaration: Getting an S3 SDK client
func getS3Client(ctx context.Context, endpoint string, credentials *credentials.StaticCredentialsProvider) (*s3.Client, error) {
	// Necessary environment variable name.
	const apiUrl = "apiUrl"
	// Getting the environment variable value.
	
	// Obtaining the S3 SDK client configuration based on the passed parameters.
	cnf, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(credentials),
		// Zerops supports only the S3 default region for API calls.
		// It doesn't mean that the physical HW infrastructure is located there also.
		// All Zerops infrastructure is completely located in Europe/Prague.
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	// Create a new S3 SDK client instance.
	s3Client := s3.NewFromConfig(
		// Passing the S3 SDK client configuration created before.
		cnf,
		s3.WithEndpointResolver(
			// Applying of the Zerops Object Storage API URL endpoint.
			s3.EndpointResolverFromURL(endpoint),
		),
		func(opts *s3.Options) {
			// Zerops supports currently only S3 path-style addressing model.
			// The virtual-hosted style model will be supported in near future.
			opts.UsePathStyle = true
		},
	)
	if s3Client != nil {
		return s3Client, nil
	}
	return nil, errors.New("creating an S3 SDK client failed")
}
