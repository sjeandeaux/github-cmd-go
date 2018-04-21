// Package information on project
package information_test

import (
	"testing"

	"github.com/sjeandeaux/toolators/information"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	assert.NotEmpty(t, information.Print())
}
