package config

import (
	"log/slog"
	"time"

	"github.com/renantatsuo/envv"
)

type Config struct {
	LogLevel         slog.Level
	Port             int
	StoreDir         string
	ReviewsTimeLimit time.Duration
	PollingInterval  time.Duration
	DatabaseConnStr  string
	QueueConnStr     string
}

// LoadConfigFromEnv loads the config from the environment variables.
func LoadConfigFromEnv() (Config, error) {
	port := envv.Get("PORT").Int().Default(8080).Parse()
	logLevelStr := envv.Get("LOG_LEVEL").String().Default("debug").Parse()
	storeDir := envv.Get("STORE_DIR").String().Default("data").Parse()
	reviewsTimeLimit := envv.Get("REVIEWS_TIME_LIMIT").Duration().Default(48 * time.Hour).Parse()
	pollingInterval := envv.Get("POLLING_INTERVAL").Duration().Default(30 * time.Second).Parse()
	databaseConnStr := envv.Get("DATABASE_CONN_STR").String().Default("data/database.db").Parse()
	queueConnStr := envv.Get("QUEUE_CONN_STR").String().Default("data/queue.db").Parse()

	logLevel, err := parseLogLevel(logLevelStr)
	if err != nil {
		return Config{}, err
	}

	return Config{
		LogLevel:         logLevel,
		Port:             port,
		StoreDir:         storeDir,
		ReviewsTimeLimit: reviewsTimeLimit,
		PollingInterval:  pollingInterval,
		DatabaseConnStr:  databaseConnStr,
		QueueConnStr:     queueConnStr,
	}, nil
}

func parseLogLevel(logLevel string) (l slog.Level, err error) {
	err = l.UnmarshalText([]byte(logLevel))
	return
}
