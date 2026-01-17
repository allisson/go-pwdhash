package cast

import (
	"strconv"

	"github.com/ccoveille/go-safecast/v2"
)

// ConvertStringToUint32 parses a base-10 string and returns a uint32 value.
func ConvertStringToUint32(s string) (uint32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return uint32(0), err
	}

	return safecast.Convert[uint32](i)
}

// ConvertStringToUint8 parses a base-10 string and returns a uint8 value.
func ConvertStringToUint8(s string) (uint8, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return uint8(0), err
	}

	return safecast.Convert[uint8](i)
}

// ConvertIntToUint32 safely casts an int to uint32.
func ConvertIntToUint32(i int) (uint32, error) {
	return safecast.Convert[uint32](i)
}
