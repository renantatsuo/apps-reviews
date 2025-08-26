package config

import (
	"log/slog"
	"strings"
	"time"

	"github.com/renantatsuo/envv"
)

type Config struct {
	LogLevel         slog.Level
	Port             int
	StoreDir         string
	AppIDs           []string
	ReviewsTimeLimit time.Duration
	PollingInterval  time.Duration
}

// LoadConfigFromEnv loads the config from the environment variables.
func LoadConfigFromEnv() (*Config, error) {
	port := envv.Get("PORT").Int().Default(8080).Parse()
	logLevelStr := envv.Get("LOG_LEVEL").String().Default("debug").Parse()
	// defaults to hevy and instagram
	appleAppIDsStr := envv.Get("APP_IDS").String().Default("1458862350,389801252").Parse()
	appleAppIDs := strings.Split(appleAppIDsStr, ",")
	storeDir := envv.Get("STORE_DIR").String().Default("data").Parse()
	reviewsTimeLimit := envv.Get("REVIEWS_TIME_LIMIT").Duration().Default(48 * time.Hour).Parse()
	pollingInterval := envv.Get("POLLING_INTERVAL").Duration().Default(30 * time.Second).Parse()

	logLevel, err := parseLogLevel(logLevelStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		LogLevel:         logLevel,
		Port:             port,
		StoreDir:         storeDir,
		AppIDs:           appleAppIDs,
		ReviewsTimeLimit: reviewsTimeLimit,
		PollingInterval:  pollingInterval,
	}, nil
}

func parseLogLevel(logLevel string) (l slog.Level, err error) {
	err = l.UnmarshalText([]byte(logLevel))
	return
}
