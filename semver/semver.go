// Package semver https://semver.org/
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

//Version the version X.Y.Z. (TODO prebuild)
type Version struct {
	// Major information
	Major int64
	// Minor information
	Minor int64
	// Patch information
	Patch int64
}

//NewVersion parse the value in Version
func NewVersion(value string) (*Version, error) {
	return parse(value)
}

//parse value X.Z.Y
func parse(value string) (*Version, error) {
	const (
		expectedSize     = 3
		separatorVersion = "."
		separatorPatch   = "-"
	)

	maMiPa := strings.SplitN(value, separatorVersion, expectedSize)
	if len(maMiPa) < expectedSize {
		return nil, fmt.Errorf("%q is bad", value)
	}
	major, err := convert(maMiPa[0])
	if err != nil {
		return nil, err
	}

	minor, err := convert(maMiPa[1])
	if err != nil {
		return nil, err
	}

	patchWithMeta := strings.Split(maMiPa[2], separatorPatch)
	patch, err := convert(patchWithMeta[0])
	if err != nil {
		return nil, err
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func convert(value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%q is bad", value)
	}
	return parsed, nil
}
