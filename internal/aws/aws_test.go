//Package aws common
package aws

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
)

var (
	testConfigFilename = filepath.Join("testdata", "config")
)

const assumeRoleRespMsg = `
<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <AssumedRoleUser>
      <Arn>arn:aws:sts::account_id:assumed-role/role/session_name</Arn>
      <AssumedRoleId>AKID:session_name</AssumedRoleId>
    </AssumedRoleUser>
    <Credentials>
      <AccessKeyId>PATRICK_STAR_AKID</AccessKeyId>
      <SecretAccessKey>OncleBens_Secret</SecretAccessKey>
      <SessionToken>ONCLE_BENS_SESSION_TOKEN</SessionToken>
      <Expiration>%s</Expiration>
    </Credentials>
  </AssumeRoleResult>
  <ResponseMetadata>
    <RequestId>request-id</RequestId>
  </ResponseMetadata>
</AssumeRoleResponse>
`

var dataConfigCredentials = []struct {
	awsProvider credentials.Value
	ourConfig   *CredentialsConfig
}{
	{
		awsProvider: credentials.Value{
			AccessKeyID:     "AccessKey",
			ProviderName:    "StaticProvider",
			SecretAccessKey: "SecretKey",
			SessionToken:    "Token",
		},
		ourConfig: &CredentialsConfig{
			AccessKey: "AccessKey",
			SecretKey: "SecretKey",
			Region:    "us-west-2",
			Token:     "Token",
		},
	},

	{
		awsProvider: credentials.Value{
			AccessKeyID:     "boom-my-access-key-id",
			ProviderName:    "SharedCredentialsProvider",
			SecretAccessKey: "boom-my-secret-access-key",
			SessionToken:    "",
		},
		ourConfig: &CredentialsConfig{
			Profile:  "oncle-bens",
			Filename: testConfigFilename,
			Region:   "us-west-2",
		},
	},
	{
		awsProvider: credentials.Value{
			AccessKeyID:     "PATRICK_STAR_AKID",
			ProviderName:    "AssumeRoleProvider",
			SecretAccessKey: "OncleBens_Secret",
			SessionToken:    "ONCLE_BENS_SESSION_TOKEN",
		},
		ourConfig: &CredentialsConfig{
			AccessKey: "AccessKey",
			SecretKey: "SecretKey",
			Region:    "us-west-2",
			Token:     "Token",
			RoleARN:   "arn:aws:iam::123456789012:user/OncleBens",
		},
	},
}

func TestConfig_Credentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(fmt.Sprintf(assumeRoleRespMsg, time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))
	}))
	defer server.Close()

	for _, d := range dataConfigCredentials {
		d.ourConfig.Endpoint = server.URL
		d.ourConfig.DisableSSL = true

		actualCredential, err := (*(d.ourConfig.Credentials().Config).Credentials).Get()
		assert.Nil(t, err)
		assert.Equal(t, d.awsProvider, actualCredential)
	}
}
