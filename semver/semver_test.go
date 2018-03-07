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
			input:         "666.666.c",
			expectedError: "\"c\" is bad",
		},

		{
			input:         "d",
			expectedError: "\"d\" is bad",
		},
	}

	for _, data := range inputs {
		_, actualErr := NewVersion(data.input)
		assert.EqualError(t, actualErr, data.expectedError)
	}
}

func TestIncrementOk(t *testing.T) {
	var inputs = []struct {
		value    *Version
		expected *Version
		position string
	}{
		{
			position: PositionMajor,
			value:    &Version{Major: 332, Minor: 666, Patch: 999},
			expected: &Version{Major: 333, Minor: 0, Patch: 0},
		},

		{
			position: PositionMinor,
			value:    &Version{Major: 333, Minor: 665, Patch: 999},
			expected: &Version{Major: 333, Minor: 666, Patch: 0},
		},

		{
			position: PositionPatch,
			value:    &Version{Major: 333, Minor: 666, Patch: 998},
			expected: &Version{Major: 333, Minor: 666, Patch: 999},
		},
	}

	for _, data := range inputs {
		actual, iWantNil := data.value.Increment(data.position)
		assert.Equal(t, data.expected, actual)
		assert.Nil(t, iWantNil)
	}
}

func TestIncrementKo(t *testing.T) {
	var inputs = []struct {
		value         *Version
		expectedError string
		position      string
	}{
		{
			position:      "unknown",
			value:         &Version{Major: 332, Minor: 666, Patch: 999},
			expectedError: "\"unknown\" is unknowndd",
		},
	}

	for _, data := range inputs {
		_, actualErr := data.value.Increment(data.position)
		assert.EqualError(t, actualErr, data.expectedError)
	}
}

func TestString(t *testing.T) {
	var inputs = []struct {
		value    *Version
		expected string
	}{
		{

			value:    &Version{Major: 333, Minor: 666, Patch: 999},
			expected: "333.666.999",
		},
	}

	for _, data := range inputs {
		actual := data.value.String()
		assert.Equal(t, data.expected, actual)
	}
}
