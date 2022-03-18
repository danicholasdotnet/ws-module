package wxws

import (
	"os"
	"time"
)

var (
	WriteWait, _   = time.ParseDuration("5s")
	MaxMessageSize = int64(4096)
	PongWait, _    = time.ParseDuration("10s")
	PingPeriod, _  = time.ParseDuration("5s")
	Host           = os.Getenv("WS_HOST")
	Port           = os.Getenv("WS_PORT")
)
