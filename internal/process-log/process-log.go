package process_log

import "sync"

type ProcessLog struct {
	mu            sync.Mutex
	DoneCountLine int64          `json:"done_count_line"`
	SkippedLines  []ErrorLogLine `json:"skipped_line"`
}

func (pl *ProcessLog) AddDoneCountLine() {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.DoneCountLine += 1
}

func (pl *ProcessLog) AddErrorLine(errLogLine ErrorLogLine) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.SkippedLines = append(pl.SkippedLines, errLogLine)
}

type ErrorLogLine struct {
	FullLine string `json:"full_line"`
	Err      string `json:"err"`
}

func NewProcessLog() *ProcessLog {
	return &ProcessLog{
		SkippedLines: make([]ErrorLogLine, 0), // Остальное можно не заполнять
	}
}
