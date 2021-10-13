package date_utils

import "time"

const (
	apiDateLayout = time.RFC3339
)

func GetNowUtc() time.Time {
	return time.Now().UTC()
}

func GetNowUtcString() string {
	return GetNowUtc().Format(apiDateLayout)
}
