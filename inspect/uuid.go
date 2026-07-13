package inspect

import "github.com/google/uuid"

// IsUUID returns true if the subject looks like a UUID.
func IsUUID(subject any) bool {
	switch subject := subject.(type) {
	case *uuid.UUID, uuid.UUID:
		return true
	case string:
		return isUUIDBytes([]byte(subject))
	case []byte:
		return isUUIDBytes(subject)
	default:
		return false
	}
}

// isUUIDBytes returns true if the supplied byte sequence looks like a UUID.
func isUUIDBytes(subject []byte) bool {
	size := len(subject)
	if size != 36 && size != 32 {
		return false
	}
	for x, ch := range subject {
		if size == 36 && (x == 8 || x == 13 || x == 18 || x == 23) && ch == '-' {
			continue
		}
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) { //nolint:staticcheck
			return false
		}
	}
	return true
}
