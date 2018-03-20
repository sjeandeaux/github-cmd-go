package semver_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/sjeandeaux/github-cmd-go/semver"
	"github.com/stretchr/testify/assert"
)

func goWantHelperProcessArgs() []string {
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}
	return args
}
func fakeExecutor(t *testing.T, functName string) semver.Executor {
	return func(name string, arg ...string) *exec.Cmd {

		base := filepath.Base(os.Args[0])
		dir := filepath.Dir(os.Args[0])
		if dir == "." {
			t.Skip("skipping; running test at root somehow")
		}
		parentDir := filepath.Dir(dir)
		dirBase := filepath.Base(dir)
		if dirBase == "." {
			t.Skipf("skipping; unexpected shallow dir of %q", dir)
		}
		testArgs := make([]string, 3)
		testArgs[0] = fmt.Sprint("-test.run=", functName)
		testArgs[1] = "--"
		testArgs[2] = name
		cmd := exec.Command(filepath.Join(dirBase, base), append(testArgs, arg...)...)
		cmd.Dir = parentDir
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}

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
			expectedError: "\"unknown\" is unknown",
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
	semver.SetExecutor(fakeExecutor(t, "TestNewGitVersionOKHelper"))
	defer semver.SetExecutorWithDefault()
	actual, iWantNil := semver.NewGitVersion()
	assert.Nil(t, iWantNil)
	assert.Equal(t, &semver.Version{Major: 6, Minor: 6, Patch: 6}, actual)
}

func TestNewGitVersionOKHelper(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := goWantHelperProcessArgs()
	cmdGit := args[1]
	if cmdGit == "rev-list" {
		fmt.Println("commit-sha")
		os.Exit(0)
	} else {
		fmt.Println("6.6.6")
		os.Exit(0)
	}

}

func TestNewGitVersionWithoutTag(t *testing.T) {
	semver.SetExecutor(fakeExecutor(t, "TestNewGitVersionWithoutTagHelper"))
	defer semver.SetExecutorWithDefault()

	actual, iWantNil := semver.NewGitVersion()
	assert.Nil(t, iWantNil)
	assert.Equal(t, &semver.Version{Major: 0, Minor: 0, Patch: 0}, actual)
}

func TestNewGitVersionWithoutTagHelper(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	os.Exit(-1)
}

func TestNewGitVersionKo(t *testing.T) {
	semver.SetExecutor(fakeExecutor(t, "TestNewGitVersionKoHelper"))
	defer semver.SetExecutorWithDefault()

	actual, iWantNil := semver.NewGitVersion()
	assert.NotNil(t, iWantNil)
	assert.Nil(t, actual)
}

func TestNewGitVersionKoHelper(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := goWantHelperProcessArgs()
	cmdGit := args[1]
	if cmdGit == "rev-list" {
		fmt.Println("commit-sha")
		os.Exit(0)
	}
	os.Exit(-2)

}
