package apilogger

import logger "github.com/sirupsen/logrus"

/*
APILogger print log using logger
*/
type APILogger struct {
	component string
	level     string
	service   string
}

var uniqueLogger = make(map[string]APILogger)

func GetLogger(level string, component string, service string) *APILogger {
	log, ok := uniqueLogger[service]
	if ok {
		return &log
	}
	typeLevel, _ := logger.ParseLevel(level)
	logger.SetLevel(typeLevel)

	newLogger := APILogger{component, level, service}
	uniqueLogger[service] = APILogger{component, level, service}
	return &newLogger
}

func (a *APILogger) PrintLogError(err error, phase string, errorMessage string) bool {
	doMessage(a, phase, err).Error(errorMessage)
	return true
}

func (a *APILogger) PrintLogWarn(err error, phase string, errorMessage string) bool {
	doMessage(a, phase, err).Warn(errorMessage)
	return true
}

func (a *APILogger) PrintLogInfo(phase string, message string) bool {
	doMessage(a, phase, nil).Info(message)
	return true
}

func doMessage(a *APILogger, phase string, err error) *logger.Entry {
	entry := logger.WithFields(logger.Fields{
		"Component": a.component,
		"Service":   a.service,
		"Phase":     phase,
		"Error":     err,
	})
	return entry
}
