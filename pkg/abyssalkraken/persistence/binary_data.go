package persistence

import (
	"bytes"
	"fmt"
)

type BinaryData struct {
	Version int64
	Data    []byte
}

func (b *BinaryData) Equals(other *BinaryData) bool {
	if other == nil {
		return false
	}
	return b.Version == other.Version && bytes.Equal(b.Data, other.Data)
}

func (b *BinaryData) HashCode() int {
	hash := int(b.Version)
	for _, v := range b.Data {
		hash = 31*hash + int(v)
	}
	return hash
}

func (b *BinaryData) ToString() string {
	return fmt.Sprintf("BinaryData(version=%d)", b.Version)
}
