package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"mini-etl-pipeline/internal/pipeline"
	processlog "mini-etl-pipeline/internal/process-log"
	"os"
	"os/signal"
)

func RunServer() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("вышла ошибка при чтении конфига %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logFileName := os.Getenv("LOG_FILE_NAME")
	processLog := processlog.NewProcessLog()
	if err := pipeline.Init(ctx, logFileName, processLog); err != nil {
		return err
	}

	return nil
}
