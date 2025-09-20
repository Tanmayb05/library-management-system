package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	setupLogger()
}

func setupLogger() {
	// Set log level based on environment
	level := strings.ToLower(getEnv("LOG_LEVEL", "info"))
	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn", "warning":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Logger.SetLevel(logrus.FatalLevel)
	case "panic":
		Logger.SetLevel(logrus.PanicLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	env := strings.ToLower(getEnv("APP_ENV", "development"))
	if env == "production" {
		// JSON format for production
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
	} else {
		// Text format for development
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set output
	Logger.SetOutput(os.Stdout)

	// Add default fields
	Logger.WithFields(logrus.Fields{
		"service": "library-management",
		"version": getEnv("APP_VERSION", "1.0.0"),
	}).Info("Logger initialized")
}

// GetLogger returns the configured logger instance
func GetLogger() *logrus.Logger {
	return Logger
}

// WithFields creates a new logger entry with the specified fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

// WithField creates a new logger entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

// WithError creates a new logger entry with an error field
func WithError(err error) *logrus.Entry {
	return Logger.WithError(err)
}

// Info logs a message at Info level
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof logs a formatted message at Info level
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn logs a message at Warn level
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf logs a formatted message at Warn level
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error logs a message at Error level
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf logs a formatted message at Error level
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal logs a message at Fatal level and exits
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf logs a formatted message at Fatal level and exits
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Debug logs a message at Debug level
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf logs a formatted message at Debug level
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}