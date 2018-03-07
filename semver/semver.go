// Package semver https://semver.org/
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

//Position Major, Minor and Patch
type Position int

//NamePosition the name and his position
type NamePosition map[string]Position

const (
	//PositionMajor position major
	PositionMajor = "major"
	//PositionMinor position minor
	PositionMinor = "minor"
	//PositionPatch position patch
	PositionPatch = "patch"
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

//Increment increment the version
func (v *Version) Increment(po string) (*Version, error) {
	switch po {
	case PositionMajor:
		return &Version{
			Major: v.Major + 1,
			Minor: 0,
			Patch: 0,
		}, nil
	case PositionMinor:
		return &Version{
			Major: v.Major,
			Minor: v.Minor + 1,
			Patch: 0,
		}, nil
	case PositionPatch:
		return &Version{
			Major: v.Major,
			Minor: v.Minor,
			Patch: v.Patch + 1,
		}, nil
	default:
		return nil, fmt.Errorf("%q is unknowndd", po)
	}
}

//String print X.Y.Z
func (v *Version) String() string {
	const format = "%d.%d.%d"
	return fmt.Sprintf(format, v.Major, v.Minor, v.Patch)
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
