package slog

import (
	"io"
	"os"
)

// SugaredLogger definition.
// Is a fast and usable Logger, which already contains
// the default formatting and handling capabilities
type SugaredLogger struct {
	*Logger
	// Formatter log message formatter. default use TextFormatter
	Formatter Formatter
	// Output writer
	Output io.Writer
	// Level for log handling. if log record level <= Level, it will be record.
	Level Level
}

// NewSugaredLogger create new SugaredLogger
func NewSugaredLogger(output io.Writer, level Level) *SugaredLogger {
	sl := &SugaredLogger{
		Level:  level,
		Output: output,
		Logger: New(),
		// default value
		Formatter: NewTextFormatter(),
	}

	// NOTICE: use self as an log handler
	sl.AddHandler(sl)

	return sl
}

// NewJSONSugared create new SugaredLogger with JSONFormatter
func NewJSONSugared(out io.Writer, level Level) *SugaredLogger {
	sl := NewSugaredLogger(out, level)
	sl.Formatter = NewJSONFormatter()

	return sl
}

// Configure current logger
func (sl *SugaredLogger) Configure(fn func(sl *SugaredLogger)) *SugaredLogger {
	fn(sl)
	return sl
}

// Reset the logger
func (sl *SugaredLogger) Reset() {
	*sl = *NewSugaredLogger(os.Stdout, DebugLevel)
}

// IsHandling Check if the current level can be handling
func (sl *SugaredLogger) IsHandling(level Level) bool {
	return sl.Level.ShouldHandling(level)
}

// Handle log record
func (sl *SugaredLogger) Handle(record *Record) error {
	bts, err := sl.Formatter.Format(record)
	if err != nil {
		return err
	}

	_, err = sl.Output.Write(bts)
	return err
}

// Close all log handlers
func (sl *SugaredLogger) Close() error {
	sl.Logger.VisitAll(func(handler Handler) error {
		if _, ok := handler.(*SugaredLogger); !ok {
			_ = handler.Close()
		}
		return nil
	})
	return nil
}

// Flush all logs. alias of the FlushAll()
func (sl *SugaredLogger) Flush() error {
	return sl.FlushAll()
}

// FlushAll all logs
func (sl *SugaredLogger) FlushAll() error {
	sl.Logger.VisitAll(func(handler Handler) error {
		if _, ok := handler.(*SugaredLogger); !ok {
			_ = handler.Flush()
		}
		return nil
	})
	return nil
}