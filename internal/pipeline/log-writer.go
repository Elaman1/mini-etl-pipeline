package pipeline

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	accesslog "mini-etl-pipeline/internal/access-log"
	processlog "mini-etl-pipeline/internal/process-log"
	"sync"
)

type Writer struct {
	mu     sync.Mutex
	logs   *processlog.ProcessLog
	writer *bufio.Writer
	closer io.Closer
}

func NewLogWriter(writer io.WriteCloser, logs *processlog.ProcessLog) *Writer {
	return &Writer{
		logs:   logs,
		writer: bufio.NewWriter(writer),
		closer: writer,
	}
}

func (w *Writer) Write(entry *accesslog.AccessLogEntry) {
	w.mu.Lock()
	defer w.mu.Unlock()

	mapLog := map[string]interface{}{
		"path":      entry.Path,
		"method":    entry.Method,
		"status":    entry.Status,
		"timestamp": entry.Timestamp,
		"ip":        entry.IP,
		"userAgent": entry.UserAgent,
	}

	bytes, err := json.Marshal(mapLog)
	if err != nil {
		w.logs.AddErrorLine(processlog.ErrorLogLine{
			FullLine: entry.FullLine,
			Err:      fmt.Sprintf("ошибка JSON: %v", err),
		})
		log.Printf("JSON error: %v\n", err)
		return
	}

	if _, err = w.writer.WriteString(string(bytes) + "\n"); err != nil {
		w.logs.AddErrorLine(processlog.ErrorLogLine{
			FullLine: entry.FullLine,
			Err:      fmt.Sprintf("ошибка записи в файл: %v", err),
		})
		log.Printf("write error: %v\n", err)
		return
	}
}

func (w *Writer) CloseLog() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.writer.Flush(); err != nil {
		log.Printf("flush error: %v", err)
		return
	}

	err := w.closer.Close()
	if err != nil {
		log.Printf("close error: %v", err)
		return
	}
}
