package utils

import "time"

// GetCurrentTimestamp returns current Unix timestamp in seconds
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentTimestampMillis returns current Unix timestamp in milliseconds
func GetCurrentTimestampMillis() int64 {
	return time.Now().UnixMilli()
}

// AddDuration adds a duration to the current time and returns Unix timestamp
func AddDuration(duration time.Duration) int64 {
	return time.Now().Add(duration).Unix()
}

// IsExpired checks if a Unix timestamp has expired
func IsExpired(timestamp int64) bool {
	if timestamp == 0 {
		return false
	}
	return time.Now().Unix() > timestamp
}

// TimeUntilExpiry returns the duration until a timestamp expires
// Returns 0 if already expired
func TimeUntilExpiry(timestamp int64) time.Duration {
	if timestamp == 0 {
		return 0
	}
	expiry := time.Unix(timestamp, 0)
	duration := time.Until(expiry)
	if duration < 0 {
		return 0
	}
	return duration
}

// FormatTimestamp formats a Unix timestamp to RFC3339 string
func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}
