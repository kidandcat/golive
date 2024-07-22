package frontend

import (
	"time"
)

type Function struct {
	Stack   string
	Status  string
	Time    time.Duration
	Started map[int64]time.Time
}
