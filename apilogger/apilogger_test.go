package apilogger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	var logger = GetLogger("debug", "API-endpoint-light", "service")

	logger.PrintLogError(fmt.Errorf("error"), "phase", "message")

	logger.PrintLogWarn(fmt.Errorf("error"), "phase", "message")

	logger.PrintLogInfo("phase", "message")

	logger = GetLogger("debug", "API-endpoint-light", "service")
}
