package access_log

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var logPattern = regexp.MustCompile(
	`^(?P<ip>\S+) - - \[(?P<time>[^\]]+)] "(?P<method>\S+)\s(?P<path>\S+)\s(?P<protocol>\S+)" (?P<status>\d{3}) (?P<size>\d+) "(?P<referer>[^"]*)" "(?P<agent>[^"]*)"$`,
)

const timeLayout = "02/Jan/2006:15:04:05 -0700"

func ParseAccessLogLine(line string) (*AccessLogEntry, error) {
	match := logPattern.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("line doesn't match expected format")
	}

	names := logPattern.SubexpNames()
	data := make(map[string]string)
	for i, name := range names {
		if i != 0 {
			data[name] = match[i]
		}
	}

	timestamp, err := time.Parse(timeLayout, data["time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	status, err := strconv.Atoi(data["status"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse status code: %v", err)
	}

	size, err := strconv.Atoi(data["size"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse size: %v", err)
	}

	return &AccessLogEntry{
		IP:        data["ip"],
		Timestamp: timestamp,
		Method:    data["method"],
		Path:      data["path"],
		Protocol:  data["protocol"],
		Status:    status,
		Size:      size,
		Referer:   data["referer"],
		UserAgent: data["agent"],
		FullLine:  line,
	}, nil
}
