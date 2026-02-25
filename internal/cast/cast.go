// Package cast provides narrow numeric conversion helpers.
package cast

import (
	"github.com/ccoveille/go-safecast/v2"
)

// ConvertStringToUint32 parses a base-10 string and returns a uint32 value.
func ConvertStringToUint32(s string) (uint32, error) {
	return safecast.Parse[uint32](s)
}

// ConvertStringToUint8 parses a base-10 string and returns a uint8 value.
func ConvertStringToUint8(s string) (uint8, error) {
	return safecast.Parse[uint8](s)
}

// ConvertIntToUint32 safely casts an int to uint32.
func ConvertIntToUint32(i int) (uint32, error) {
	return safecast.Convert[uint32](i)
}
