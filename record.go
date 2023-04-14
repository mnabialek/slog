package slog

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// Record a log record definition
type Record struct {
	logger *Logger

	// Time for record log
	Time time.Time
	// Level log level for record
	Level Level
	// level name cache from Level
	levelName string

	// Channel log channel name. eg: "order", "goods", "user"
	Channel string
	Message string

	// Ctx context.Context
	Ctx context.Context

	// Fields custom fields data.
	// Contains all the fields set by the user.
	Fields M
	// Data log context data
	Data M
	// Extra log extra data
	Extra M

	// Caller information
	Caller *runtime.Frame
	// CallerFlag value. default is equals to Logger.CallerFlag
	CallerFlag uint8
	// CallerSkip value. default is equals to Logger.CallerSkip
	CallerSkip int

	// Buffer Can use Buffer on formatter
	// Buffer *bytes.Buffer

	// cache the r.Time.Nanosecond() / 1000
	// microSecond int
	// stacks []byte
	// formatted []byte
	// field caches mapping for optimize performance. TODO use map[string][]byte ?
	// strMp map[string]string
}

func newRecord(logger *Logger) *Record {
	return &Record{
		logger:  logger,
		Channel: DefaultChannelName,
		// with some options
		CallerFlag: logger.CallerFlag,
		CallerSkip: logger.CallerSkip,
		// init map data field
		// Data:   make(M, 2),
		// Extra:  make(M, 0),
		// Fields: make(M, 0),
	}
}

//
// ---------------------------------------------------------------------------
// Copy record with something
// ---------------------------------------------------------------------------
//

// WithTime set the record time
func (r *Record) WithTime(t time.Time) *Record {
	nr := r.Copy()
	nr.Time = t
	return nr
}

// WithContext on record
func (r *Record) WithContext(ctx context.Context) *Record {
	nr := r.Copy()
	nr.Ctx = ctx
	return nr
}

// WithError on record
func (r *Record) WithError(err error) *Record {
	return r.WithFields(M{FieldKeyError: err})
}

// WithData on record
func (r *Record) WithData(data M) *Record {
	nr := r.Copy()
	// if nr.Data == nil {
	// 	nr.Data = data
	// }

	nr.Data = data
	return nr
}

// WithField with a new field to record
func (r *Record) WithField(name string, val any) *Record {
	return r.WithFields(M{name: val})
}

// WithFields with new fields to record
func (r *Record) WithFields(fields M) *Record {
	nr := r.Copy()
	if nr.Fields == nil {
		nr.Fields = make(M, len(fields))
	}

	for k, v := range fields {
		nr.Fields[k] = v
	}
	return nr
}

// Copy new record from old record
func (r *Record) Copy() *Record {
	dataCopy := make(M, len(r.Data))
	for k, v := range r.Data {
		dataCopy[k] = v
	}

	fieldsCopy := make(M, len(r.Fields))
	for k, v := range r.Fields {
		fieldsCopy[k] = v
	}

	extraCopy := make(M, len(r.Extra))
	for k, v := range r.Extra {
		extraCopy[k] = v
	}

	return &Record{
		logger:     r.logger,
		Channel:    r.Channel,
		Time:       r.Time,
		Level:      r.Level,
		levelName:  r.levelName,
		CallerFlag: r.CallerFlag,
		CallerSkip: r.CallerSkip,
		Message:    r.Message,
		Data:       dataCopy,
		Extra:      extraCopy,
		Fields:     fieldsCopy,
	}
}

//
// ---------------------------------------------------------------------------
// Direct set value to record
// ---------------------------------------------------------------------------
//

// SetCtx on record
func (r *Record) SetCtx(ctx context.Context) *Record { return r.SetContext(ctx) }

// SetContext on record
func (r *Record) SetContext(ctx context.Context) *Record {
	r.Ctx = ctx
	return r
}

// SetData on record
func (r *Record) SetData(data M) *Record {
	r.Data = data
	return r
}

// AddData on record
func (r *Record) AddData(data M) *Record {
	if r.Data == nil {
		r.Data = data
		return r
	}

	for k, v := range data {
		r.Data[k] = v
	}
	return r
}

// AddValue add Data value to record
func (r *Record) AddValue(key string, value any) *Record {
	if r.Data == nil {
		r.Data = make(M, 8)
		return r
	}

	r.Data[key] = value
	return r
}

// SetExtra information on record
func (r *Record) SetExtra(data M) *Record {
	r.Extra = data
	return r
}

// AddExtra information on record
func (r *Record) AddExtra(data M) *Record {
	if r.Extra == nil {
		r.Extra = data
		return r
	}

	for k, v := range data {
		r.Extra[k] = v
	}
	return r
}

// SetExtraValue on record
func (r *Record) SetExtraValue(k string, v any) {
	if r.Extra == nil {
		r.Extra = make(M, 8)
	}
	r.Extra[k] = v
}

// SetTime on record
func (r *Record) SetTime(t time.Time) *Record {
	r.Time = t
	return r
}

// AddField add new field to the record
func (r *Record) AddField(name string, val any) *Record {
	if r.Fields == nil {
		r.Fields = make(M, 8)
	}

	r.Fields[name] = val
	return r
}

// AddFields add new fields to the record
func (r *Record) AddFields(fields M) *Record {
	if r.Fields == nil {
		r.Fields = fields
		return r
	}

	for n, v := range fields {
		r.Fields[n] = v
	}
	return r
}

// SetFields to the record
func (r *Record) SetFields(fields M) *Record {
	r.Fields = fields
	return r
}

//
// ---------------------------------------------------------------------------
// Add log message with builder
// TODO r.Build(InfoLevel).Str().Int().Float().Msg()
// ---------------------------------------------------------------------------
//

// Object data on record TODO optimize performance
// func (r *Record) Obj(obj fmt.Stringer) *Record {
// 	r.Data = ctx
// 	return r
// }

// Object data on record TODO optimize performance
// func (r *Record) Any(v any) *Record {
// 	r.Data = ctx
// 	return r
// }

// func (r *Record) Str(message string) {
// 	r.logWrite(level, []byte(message))
// }

// func (r *Record) Int(val int) {
// 	r.logWrite(level, []byte(message))
// }

//
// ---------------------------------------------------------------------------
// Add log message with level
// ---------------------------------------------------------------------------
//

func (r *Record) log(level Level, args []any) {
	r.Level = level
	// will reduce memory allocation once
	// r.Message = strutil.Byte2str(formatArgsWithSpaces(args))
	r.Message = formatArgsWithSpaces(args)
	// r.logWrite(level)
	r.logger.writeRecord(level, r)
}

func (r *Record) logf(level Level, format string, args []any) {
	r.Level = level
	r.Message = fmt.Sprintf(format, args...)
	// r.logWrite(level)
	r.logger.writeRecord(level, r)
}

// Log a message with level
func (r *Record) Log(level Level, args ...any) { r.log(level, args) }

// Logf a message with level
func (r *Record) Logf(level Level, format string, args ...any) {
	r.logf(level, format, args)
}

// Info logs a message at level Info
func (r *Record) Info(args ...any) { r.log(InfoLevel, args) }

// Infof logs a message at level Info
func (r *Record) Infof(format string, args ...any) {
	r.logf(InfoLevel, format, args)
}

// Trace logs a message at level Trace
func (r *Record) Trace(args ...any) { r.log(TraceLevel, args) }

// Tracef logs a message at level Trace
func (r *Record) Tracef(format string, args ...any) {
	r.logf(TraceLevel, format, args)
}

// Error logs a message at level Error
func (r *Record) Error(args ...any) { r.log(ErrorLevel, args) }

// Errorf logs a message at level Error
func (r *Record) Errorf(format string, args ...any) {
	r.logf(ErrorLevel, format, args)
}

// Warn logs a message at level Warn
func (r *Record) Warn(args ...any) { r.log(WarnLevel, args) }

// Warnf logs a message at level Warn
func (r *Record) Warnf(format string, args ...any) {
	r.logf(ErrorLevel, format, args)
}

// Notice logs a message at level Notice
func (r *Record) Notice(args ...any) { r.log(NoticeLevel, args) }

// Noticef logs a message at level Notice
func (r *Record) Noticef(format string, args ...any) {
	r.logf(NoticeLevel, format, args)
}

// Debug logs a message at level Debug
func (r *Record) Debug(args ...any) { r.log(DebugLevel, args) }

// Debugf logs a message at level Debug
func (r *Record) Debugf(format string, args ...any) {
	r.logf(DebugLevel, format, args)
}

// Print logs a message at level Print
func (r *Record) Print(args ...any) { r.log(PrintLevel, args) }

// Println logs a message at level Print, will not append \n. alias of Print
func (r *Record) Println(args ...any) { r.log(PrintLevel, args) }

// Printf logs a message at level Print
func (r *Record) Printf(format string, args ...any) {
	r.logf(PrintLevel, format, args)
}

// Fatal logs a message at level Fatal
func (r *Record) Fatal(args ...any) { r.log(FatalLevel, args) }

// Fatalln logs a message at level Fatal
func (r *Record) Fatalln(args ...any) { r.log(FatalLevel, args) }

// Fatalf logs a message at level Fatal
func (r *Record) Fatalf(format string, args ...any) {
	r.logf(FatalLevel, format, args)
}

// Panic logs a message at level Panic
func (r *Record) Panic(args ...any) { r.log(PanicLevel, args) }

// Panicln logs a message at level Panic
func (r *Record) Panicln(args ...any) { r.log(PanicLevel, args) }

// Panicf logs a message at level Panic
func (r *Record) Panicf(format string, args ...any) {
	r.logf(PanicLevel, format, args)
}

// ---------------------------------------------------------------------------
// helper methods
// ---------------------------------------------------------------------------

// LevelName get
func (r *Record) LevelName() string { return r.levelName }

// GoString of the record
func (r *Record) GoString() string {
	return "slog: " + r.Message
}

func (r *Record) timestamp() string {
	s := strconv.FormatInt(r.Time.UnixMicro(), 10)
	return s[:10] + "." + s[10:]
}
