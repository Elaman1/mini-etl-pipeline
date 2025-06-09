package access_log

import (
	"time"
)

type AccessLogEntry struct {
	IP        string
	Timestamp time.Time
	Method    string
	Path      string
	Protocol  string
	Status    int
	Size      int
	Referer   string
	UserAgent string
	FullLine  string
}

func (a AccessLogEntry) FailedStatus() bool {
	switch a.Status / 100 {
	case 4:
	case 5:
		return true
	default:
		return false
	}

	return true
}
