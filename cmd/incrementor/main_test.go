package main

import (
	"testing"

	"github.com/sjeandeaux/github-cmd-go/semver"
	"github.com/stretchr/testify/assert"
)

func TestCommandLineIncrementOk(t *testing.T) {
	var inputs = []struct {
		input    commandLineArgs
		expected *semver.Version
	}{
		{
			input: commandLineArgs{
				position: semver.PositionMajor,
				version:  "0.1.0",
			},
			expected: &semver.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
	}

	for _, data := range inputs {
		actual, _ := data.input.increment()
		assert.Equal(t, data.expected, actual)
	}
}

func TestCommandLineIncrementKo(t *testing.T) {

	var inputs = []struct {
		input       commandLineArgs
		expectedErr string
	}{
		{
			input: commandLineArgs{
				position: semver.PositionMajor,
				version:  "bad version",
			},
			expectedErr: "",
		},
		{
			input: commandLineArgs{
				position: "bad position",
				version:  "0.1.0",
			},
			expectedErr: "non",
		},
	}

	for _, data := range inputs {
		_, actual := data.input.increment()
		assert.Error(t, actual, data.expectedErr)
	}

}
