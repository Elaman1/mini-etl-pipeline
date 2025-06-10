package pipeline

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	accesslog "mini-etl-pipeline/internal/access-log"
	processlog "mini-etl-pipeline/internal/process-log"
	"os"
	"sync"
)

func LogReader(ctx context.Context, wg *sync.WaitGroup, file *os.File, processLog *processlog.ProcessLog) <-chan string {
	output := make(chan string, 10)

	go func() {
		defer wg.Done()
		defer close(output)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				// Если даже завершили контекст все равно записываем в логи что пропустили
				processLog.AddErrorLine(processlog.ErrorLogLine{
					FullLine: scanner.Text(),
					Err:      fmt.Sprintf("не отработано при завершении контекста %v", ctx.Err()),
				})
				continue
			default:
			}

			line := scanner.Text()
			output <- line
			log.Printf("Отправлено сообщение %s", line)
		}

		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				return
			}

			processLog.AddErrorLine(processlog.ErrorLogLine{
				FullLine: scanner.Text(),
				Err:      fmt.Sprintf("ошибка %v", err),
			})
		}
	}()

	return output
}

func LogProcessor(ctx context.Context, wg *sync.WaitGroup, ch <-chan string, processLog *processlog.ProcessLog) <-chan *accesslog.AccessLogEntry {
	output := make(chan *accesslog.AccessLogEntry, 10)

	go func() {
		defer wg.Done()
		defer close(output)
		for line := range ch {
			select {
			case <-ctx.Done():
				log.Println("Канал закрыт, завершаем Process pipeline")
				return
			default:
			}

			accessLog, err := accesslog.ParseAccessLogLine(line)
			if err != nil {
				processLog.AddErrorLine(processlog.ErrorLogLine{
					FullLine: line,
					Err:      fmt.Sprintf("Вышла ошибка при парсинге %s", err),
				})
				continue
			}

			if !accessLog.FailedStatus() {
				log.Println("Пришел успешный ответ, просто пропускаем")
				continue
			}

			output <- accessLog
		}
	}()

	return output
}

func LogWriter(ctx context.Context, wg *sync.WaitGroup, ch <-chan *accesslog.AccessLogEntry, writer *Writer) {
	defer wg.Done()

	for accessLog := range ch {
		select {
		case <-ctx.Done():
			writer.Write(&accesslog.AccessLogEntry{
				FullLine: accessLog.FullLine,
			})
			continue
		default:
		}

		writer.Write(accessLog)
	}
}
