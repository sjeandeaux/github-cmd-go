package semver

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewVersionOk(t *testing.T) {
	var inputs = []struct {
		input    string
		expected *Version
	}{
		{
			input:    "333.666.999",
			expected: &Version{Major: 333, Minor: 666, Patch: 999},
		},

		{
			input:    "333.666.999-beta",
			expected: &Version{Major: 333, Minor: 666, Patch: 999},
		},
	}

	for _, data := range inputs {
		actual, iWantNil := NewVersion(data.input)
		assert.Equal(t, data.expected, actual)
		assert.Nil(t, iWantNil)
	}
}

func TestNewVersionKo(t *testing.T) {
	var inputs = []struct {
		input         string
		expectedError string
	}{
		{
			input:         "",
			expectedError: "\"\" is bad",
		},

		{
			input:         "a.666.666",
			expectedError: "\"a\" is bad",
		},

		{
			input:         "666.b.666",
			expectedError: "\"b\" is bad",
		},

		{
			input:         "c",
			expectedError: "\"c\" is bad",
		},
	}

	for _, data := range inputs {
		_, actualErr := NewVersion(data.input)
		assert.EqualError(t, actualErr, data.expectedError)
	}
}
