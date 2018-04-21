//Package aws common
package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CredentialsConfig authentication
type CredentialsConfig struct {
	Region    string
	AccessKey string
	SecretKey string
	RoleARN   string
	Profile   string

	Filename string
	Token    string

	Endpoint   string
	DisableSSL bool
}

// CredentialsConfigEnv initialize with environmnetal variable aws_<your text>_xxx
func CredentialsConfigEnv(own string) *CredentialsConfig {
	config := CredentialsConfig{
		AccessKey: os.Getenv(fmt.Sprint("aws_", own, "_access_key")),
		SecretKey: os.Getenv(fmt.Sprint("aws_", own, "_secret_key")),
		Region:    os.Getenv(fmt.Sprint("aws_", own, "_region")),
		RoleARN:   os.Getenv(fmt.Sprint("aws_", own, "_role_arn")),
		Filename:  os.Getenv(fmt.Sprint("aws_", own, "_filename")),
		Profile:   os.Getenv(fmt.Sprint("aws_", own, "_profile")),
		Token:     os.Getenv(fmt.Sprint("aws_", own, "_token")),
	}
	return &config
}

// Credentials credentials AWS
func (c *CredentialsConfig) Credentials() *session.Session {
	userCredentials := c.userCredentials()
	// we use the user with role arn
	if c.RoleARN == "" {
		return userCredentials
	}
	return c.assumeRole(userCredentials)
}

// assumeRole assume the role with the user credentials
func (c *CredentialsConfig) assumeRole(userCredentials client.ConfigProvider) *session.Session {
	assumeConfig := &aws.Config{
		Region:      aws.String(c.Region),
		Credentials: stscreds.NewCredentials(userCredentials, c.RoleARN),
		Endpoint:    aws.String(c.Endpoint),
		DisableSSL:  aws.Bool(true),
	}
	return session.New(assumeConfig)
}

// userCredentials the user credentials
func (c *CredentialsConfig) userCredentials() *session.Session {
	config := &aws.Config{
		Region:     aws.String(c.Region),
		Endpoint:   aws.String(c.Endpoint),
		DisableSSL: aws.Bool(c.DisableSSL),
	}

	if c.AccessKey != "" || c.SecretKey != "" {
		config.WithCredentials(credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.Token))
	} else if c.Profile != "" || c.Filename != "" {
		config.WithCredentials(credentials.NewSharedCredentials(c.Filename, c.Profile))
	}
	return session.New(config)
}
