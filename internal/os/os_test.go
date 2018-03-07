package os_test

import (
	"os"
	"testing"

	internalos "github.com/sjeandeaux/github-cmd-go/internal/os"
	"github.com/stretchr/testify/assert"
)

func TestGetenv(t *testing.T) {

	os.Setenv("set-internal-os-key", "in env")
	defer func() {
		os.Setenv("set-internal-os-key", "")
	}()

	var inputs = []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{
			key:          "set-internal-os-key",
			defaultValue: "defaultValue",
			expected:     "in env",
		},
		{
			key:          "not-set-internal-os-key",
			defaultValue: "defaultValue",
			expected:     "defaultValue",
		},
	}

	for _, data := range inputs {
		actual := internalos.Getenv(data.key, data.defaultValue)
		assert.Equal(t, data.expected, actual)
	}

}
