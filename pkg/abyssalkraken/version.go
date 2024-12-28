package abyssalkraken

import (
	"fmt"
)

var MinVersion = Version{value: 0}

type Version struct {
	value int64
}

type InvalidVersionError struct {
	Version int64
	Message string
}

func (e *InvalidVersionError) Error() string {
	return fmt.Sprintf("invalid version: %d - %s", e.Version, e.Message)
}

func NewVersion(value int64) (Version, error) {
	if value < MinVersion.ToInt() {
		return Version{}, &InvalidVersionError{
			Version: value,
			Message: "Version cannot be negative",
		}
	}
	return Version{value: value}, nil
}

func (v Version) ToInt() int64 {
	return v.value
}

func (v Version) ToString() string {
	return fmt.Sprintf("%d", v.value)
}

func ToVersion(value int64) (Version, error) {
	return NewVersion(value)
}
