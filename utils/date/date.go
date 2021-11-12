package date

import "time"

const (
	ApiDateLayout = time.RFC3339
)

func GetNowUtc() time.Time {
	return time.Now().UTC()
}

func GetNowUtcString() string {
	return GetNowUtc().Format(ApiDateLayout)
}
