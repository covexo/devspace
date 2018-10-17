package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// DiscardLogger just discards every log statement
type DiscardLogger struct{}

// Debug implements logger interface
func (d *DiscardLogger) Debug(args ...interface{}) {}

// Debugf implements logger interface
func (d *DiscardLogger) Debugf(format string, args ...interface{}) {}

// Info implements logger interface
func (d *DiscardLogger) Info(args ...interface{}) {}

// Infof implements logger interface
func (d *DiscardLogger) Infof(format string, args ...interface{}) {}

// Warn implements logger interface
func (d *DiscardLogger) Warn(args ...interface{}) {}

// Warnf implements logger interface
func (d *DiscardLogger) Warnf(format string, args ...interface{}) {}

// Error implements logger interface
func (d *DiscardLogger) Error(args ...interface{}) {}

// Errorf implements logger interface
func (d *DiscardLogger) Errorf(format string, args ...interface{}) {}

// Fatal implements logger interface
func (d *DiscardLogger) Fatal(args ...interface{}) {
	os.Exit(1)
}

// Fatalf implements logger interface
func (d *DiscardLogger) Fatalf(format string, args ...interface{}) {
	os.Exit(1)
}

// Panic implements logger interface
func (d *DiscardLogger) Panic(args ...interface{}) {
	panic(fmt.Sprint(args...))
}

// Panicf implements logger interface
func (d *DiscardLogger) Panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

// Done implements logger interface
func (d *DiscardLogger) Done(args ...interface{}) {}

// Donef implements logger interface
func (d *DiscardLogger) Donef(format string, args ...interface{}) {}

// Fail implements logger interface
func (d *DiscardLogger) Fail(args ...interface{}) {}

// Failf implements logger interface
func (d *DiscardLogger) Failf(format string, args ...interface{}) {}

// Print implements logger interface
func (d *DiscardLogger) Print(level logrus.Level, args ...interface{}) {}

// Printf implements logger interface
func (d *DiscardLogger) Printf(level logrus.Level, format string, args ...interface{}) {}

// StartWait implements logger interface
func (d *DiscardLogger) StartWait(message string) {}

// StopWait implements logger interface
func (d *DiscardLogger) StopWait() {}

// With implements logger interface
func (d *DiscardLogger) With(obj interface{}) *LoggerEntry {
	return &LoggerEntry{
		logger: d,
		context: map[string]interface{}{
			"context-1": obj,
		},
	}
}

// WithKey implements logger interface
func (d *DiscardLogger) WithKey(key string, obj interface{}) *LoggerEntry {
	return &LoggerEntry{
		logger: d,
		context: map[string]interface{}{
			key: obj,
		},
	}
}

// SetLevel implements logger interface
func (d *DiscardLogger) SetLevel(level logrus.Level) {}

func (d *DiscardLogger) printWithContext(fnType logFunctionType, contextFields map[string]interface{}, args ...interface{}) {
}

func (d *DiscardLogger) printWithContextf(fnType logFunctionType, contextFields map[string]interface{}, format string, args ...interface{}) {
}

// Write implements logger interface
func (d *DiscardLogger) Write(message []byte) (int, error) {
	return len(message), nil
}

// Close implements logger interface
func (d *DiscardLogger) Close() error {
	return nil
}
