package walker

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strconv"
	"time"
)

func Init() (*Consumer, *Producer, error) {
	logger, err := createLogger()
	if err != nil {
		log.Fatalf("Could not create logger %v", err)
		return nil, nil, err
	}

	pool, err := createDbPool()
	repo := UrlRepository{
		Pool:   pool,
		Logger: logger,
	}
	if err != nil {
		log.Fatalf("could not create db pool %v", err)
		return nil, nil, err
	}
	w := Walker{
		Logger: logger,
	}
	pageChan := make(chan string)
	killChan := make(chan bool)
	producer := Producer{
		Logger:     logger,
		PageChan:   pageChan,
		Repository: repo,
		KillChan:   killChan,
	}
	consumer := Consumer{
		PageChan: pageChan,
		KillChan: killChan,
		Repo:     repo,
		Logger:   logger,
		Walker:   w,
		Producer: &producer,
		Ticker:   time.NewTicker(2 * time.Second),
	}

	return &consumer, &producer, nil
}

func createLogger() (*zap.Logger, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("could not create uuid to log %v", err)
		return nil, err
	}
	stringifiedId := id.String()
	logger := InitLogger(zap.Stringp("id", &stringifiedId))
	logger.Info("This is the first string")
	logger.Info("This is the second string")
	return logger, nil
}

func InitLogger(fields ...zap.Field) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder
	consoleEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel))
	return zap.New(core).With(fields...)
}

func createDbPool() (*sql.DB, error) {
	port, errPort := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if errPort != nil {
		return nil, errPort
	}
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		user, pass, database, host, port))
	if err != nil {
		return nil, err
	}

	return db, nil

}
