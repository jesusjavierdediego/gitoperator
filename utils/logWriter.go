package utils

import (
	"fmt"
	config "xqledger/gitoperator/configuration"
	"time"

	logger "github.com/sirupsen/logrus"
)

var Configuration config.Configuration

func getFormattedNow() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

func PrintLogError(err error, comp string, phase string, errorMessage string) bool {
	logger.WithFields(logger.Fields{
		"Time":      getFormattedNow(),
		"Component": comp,
		"Phase":     phase,
		"Error":     err,
	}).Error(errorMessage)
	return true
}

func PrintLogWarn(err error, comp string, phase string, errorMessage string) bool {
	logger.WithFields(logger.Fields{
		"Time":      getFormattedNow(),
		"Component": comp,
		"Phase":     phase,
		"Error":     err,
	}).Warn(errorMessage)
	return true
}

func PrintLogInfo(comp string, phase string, message string) bool {
	logger.WithFields(logger.Fields{
		"Time":      getFormattedNow(),
		"Component": comp,
		"Phase":     phase,
	}).Info(message)
	return true
}
