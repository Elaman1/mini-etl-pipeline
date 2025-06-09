package pipeline

import (
	"context"
	"fmt"
	processlog "mini-etl-pipeline/internal/process-log"
	"os"
	"sync"
)

func Init(ctx context.Context, logFileName string, processLog *processlog.ProcessLog) error {
	if logFileName == "" {
		return fmt.Errorf("не указано название файла")
	}

	logFile, err := os.Open(logFileName)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %v", err)
	}

	// Будем создавать 3 цепочки пайплайна
	wg := &sync.WaitGroup{}

	file, err := os.Create("result.txt")
	if err != nil {
		return fmt.Errorf("вышла ошибка при чтении файла %v", err)
	}
	defer file.Close()

	// Первый обработчик
	wg.Add(1)
	readerCh := LogReader(ctx, wg, logFile, processLog)

	// Второй обработчик
	wg.Add(1)
	processCh := LogProcessor(ctx, wg, readerCh, processLog)

	// Третий обработчик
	wg.Add(1)

	writer := NewLogWriter(file, processLog)
	defer writer.CloseLog()
	LogWriter(ctx, wg, processCh, writer)

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		// Если даже завершили вручную, все равно ждем пока отработают или положат в лог
		fmt.Println("Завершение по сигналу контекста")
		<-done // ждём завершения всех воркеров
		return nil
	case <-done:
		fmt.Println("Пайплайн завершился сам")
		return nil
	}
}
