// Package subtle provides wrappers for constant-time operations.
package subtle

import "crypto/subtle"

// ConstantTimeCompare matches slices while avoiding timing leaks.
func ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}
