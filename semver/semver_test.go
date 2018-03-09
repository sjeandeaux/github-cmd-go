package semver_test

import (
	"os/exec"
	"testing"

	"github.com/sjeandeaux/github-cmd-go/semver"
	"github.com/stretchr/testify/assert"
)

func TestNewVersionOk(t *testing.T) {
	var inputs = []struct {
		input    string
		expected *semver.Version
	}{
		{
			input:    "333.666.999",
			expected: &semver.Version{Major: 333, Minor: 666, Patch: 999},
		},

		{
			input:    "333.666.999-beta",
			expected: &semver.Version{Major: 333, Minor: 666, Patch: 999},
		},
	}

	for _, data := range inputs {
		actual, iWantNil := semver.NewVersion(data.input)
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
		_, actualErr := semver.NewVersion(data.input)
		assert.EqualError(t, actualErr, data.expectedError)
	}
}

func TestIncrementOk(t *testing.T) {
	var inputs = []struct {
		value    *semver.Version
		expected *semver.Version
		position string
	}{
		{
			position: semver.PositionMajor,
			value:    &semver.Version{Major: 332, Minor: 666, Patch: 999},
			expected: &semver.Version{Major: 333, Minor: 0, Patch: 0},
		},

		{
			position: semver.PositionMinor,
			value:    &semver.Version{Major: 333, Minor: 665, Patch: 999},
			expected: &semver.Version{Major: 333, Minor: 666, Patch: 0},
		},

		{
			position: semver.PositionPatch,
			value:    &semver.Version{Major: 333, Minor: 666, Patch: 998},
			expected: &semver.Version{Major: 333, Minor: 666, Patch: 999},
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
		value         *semver.Version
		expectedError string
		position      string
	}{
		{
			position:      "unknown",
			value:         &semver.Version{Major: 332, Minor: 666, Patch: 999},
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
		value    *semver.Version
		expected string
	}{
		{

			value:    &semver.Version{Major: 333, Minor: 666, Patch: 999},
			expected: "333.666.999",
		},
	}

	for _, data := range inputs {
		actual := data.value.String()
		assert.Equal(t, data.expected, actual)
	}
}

func TestNewGitVersionOK(t *testing.T) {
	//TODO https://golang.org/src/os/exec/exec_test.go
	fakeExectuor := func() semver.Executor {
		return func(name string, arg ...string) *exec.Cmd {
			switch arg[0] {
			case "rev-list":
				return exec.Command("echo", "response-rev-list")
			case "describe":
				assert.Equal(t, "response-rev-list", arg[2])
				return exec.Command("echo", "6.6.6")
			}

			return nil
		}
	}

	semver.SetExecutor(fakeExectuor())
	defer semver.SetExecutorWithDefault()

	actual, iWantNil := semver.NewGitVersion()
	assert.Nil(t, iWantNil)
	assert.Equal(t, &semver.Version{Major: 6, Minor: 6, Patch: 6}, actual)
}

func TestNewGitVersionWithoutTag(t *testing.T) {
	//TODO https://golang.org/src/os/exec/exec_test.go
	fakeExectuor := func() semver.Executor {
		return func(name string, arg ...string) *exec.Cmd {
			switch arg[0] {
			case "rev-list":
				return exec.Command("toutestko", "désenchantée")
			}
			return nil
		}
	}

	semver.SetExecutor(fakeExectuor())
	defer semver.SetExecutorWithDefault()

	actual, iWantNil := semver.NewGitVersion()
	assert.Nil(t, iWantNil)
	assert.Equal(t, &semver.Version{Major: 0, Minor: 0, Patch: 0}, actual)
}

func TestNewGitVersionKo(t *testing.T) {
	//TODO https://golang.org/src/os/exec/exec_test.go
	fakeExectuor := func() semver.Executor {
		return func(name string, arg ...string) *exec.Cmd {
			switch arg[0] {
			case "rev-list":
				return exec.Command("echo", "response-rev-list")
			case "describe":
				assert.Equal(t, "response-rev-list", arg[2])
				return exec.Command("toutestko", "désenchantée")
			}
			return nil
		}
	}

	semver.SetExecutor(fakeExectuor())
	defer semver.SetExecutorWithDefault()

	actual, iWantNil := semver.NewGitVersion()
	assert.NotNil(t, iWantNil)
	assert.Nil(t, actual)
}
