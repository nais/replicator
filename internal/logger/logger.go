package logger

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func SetupLogrus(level log.Level) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	log.SetFormatter(formatter)
	log.SetLevel(level)
}
