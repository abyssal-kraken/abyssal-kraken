package abyssalkraken

import (
	"fmt"
)

var MinVersion = Version{Value: 0}

type Version struct {
	Value int64
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
			Message: fmt.Sprintf("Version cannot be less than %s", MinVersion.ToString()),
		}
	}
	return Version{Value: value}, nil
}

func (v Version) ToInt() int64 {
	return v.Value
}

func (v Version) ToString() string {
	return fmt.Sprintf("%d", v.Value)
}

func ToVersion(value int64) (Version, error) {
	return NewVersion(value)
}
