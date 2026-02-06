package connection

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

// SetLogger dipanggil dari main.go
func SetLogger(logger *zap.Logger) {
	Log = logger
}
