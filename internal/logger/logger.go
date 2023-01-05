package logger

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func SetupLogrus() {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.InfoLevel)
}
