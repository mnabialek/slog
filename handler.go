package slog

import "io"

// FlushWriter is the interface satisfied by logging destinations.
type FlushWriter interface {
	Flush() error
	// Writer the output writer
	io.Writer
}

// FlushCloseWriter is the interface satisfied by logging destinations.
type FlushCloseWriter interface {
	Flush() error
	// WriteCloser the output writer
	io.WriteCloser
}

// FormatterWriterHandler interface
type FormatterWriterHandler interface {
	Handler
	// Formatter record formatter
	Formatter() Formatter
	// Writer the output writer
	Writer() io.Writer
}

//
// Handler interface
//

// Handler interface definition
type Handler interface {
	// Closer Close handler.
	// You should first call Flush() on close logic.
	// Refer the FileHandler.Close() handle
	io.Closer
	// Flush and sync logs to disk file.
	Flush() error
	// IsHandling Checks whether the given record will be handled by this handler.
	IsHandling(level Level) bool
	// Handle a log record.
	//
	// All records may be passed to this method, and the handler should discard
	// those that it does not want to handle.
	Handle(*Record) error
}

/********************************************************************************
 * Common parts for handler
 ********************************************************************************/

// LevelWithFormatter struct definition
//
// - support set log formatter
// - only support set one log level
type LevelWithFormatter struct {
	Formattable
	// Level for log message. if current level <= Level will log message
	Level Level
}

// NewLvFormatter create new instance
func NewLvFormatter(lv Level) *LevelWithFormatter {
	return &LevelWithFormatter{Level: lv}
}

// IsHandling Check if the current level can be handling
func (h *LevelWithFormatter) IsHandling(level Level) bool {
	return h.Level.ShouldHandling(level)
}

// LevelsWithFormatter struct definition
//
// - support set log formatter
// - support setting multi log levels
type LevelsWithFormatter struct {
	Formattable
	// Levels for log message
	Levels []Level
}

// NewLvsFormatter create new instance
func NewLvsFormatter(levels []Level) *LevelsWithFormatter {
	return &LevelsWithFormatter{Levels: levels}
}

// IsHandling Check if the current level can be handling
func (h *LevelsWithFormatter) IsHandling(level Level) bool {
	for _, l := range h.Levels {
		if l == level {
			return true
		}
	}
	return false
}