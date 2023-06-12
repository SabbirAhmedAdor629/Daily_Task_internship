package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"
)

const (
	Log_level_reserved = iota // 0
	Log_level_none            // 1
	Log_level_fatal           // 2
	Log_level_error           // 3
	Log_level_warning         // 4
	Log_level_info            // 5
	Log_level_debug           // 6
	Log_level_trace           // 7

	DEFAULT_CIRCUIT_BREAK = 100
)

var (
	LogLevelText = [8]string{
		Log_level_reserved: "RESERVED",
		Log_level_none:     "NONE",
		Log_level_fatal:    "FATAL",
		Log_level_error:    "ERROR",
		Log_level_warning:  "WARNING",
		Log_level_info:     "INFO",
		Log_level_debug:    "DEBUG",
		Log_level_trace:    "TRACE",
	}
)

type (
	LoggersArray []Logger
	Logger       struct{ *log.Logger }
	Fields       map[string]interface{}

	Loggers struct {
		LoggersArray // Contains individual Loggers for FATAL, ERROR, INFO, and DEBUG
		envLogLevel  int
		Entry        Fields // Contains the "default" fields for the logger
		sync.RWMutex
		CircuitBreak int
		LogCount     int
	}

	LogEntry struct {
		LoggersArray        // Contains individual Loggers for FATAL, ERROR, INFO, and DEBUG
		Entry        Fields // When a new event is created, the fields in Logger are copied and stored here.
		sync.RWMutex
	}
)

func validateCircuitBreak(num int) int { return num }
func (l *Loggers) incLogCount(num int) int {
	if l.CircuitBreak > 0 {
		l.LogCount++
	}
	return l.LogCount
}

func (l *Loggers) NewEntry(logCode int, logMsgText []string) *LogEntry {
	return l.WithFields(Fields{
		"MSG_CODE": logCode,
		"MSG_TEXT": logMsgText[logCode]})
}

func NewLoggers(out io.Writer, envLogLevel int, circuitBreak int) *Loggers {
	var (
		flags = log.LstdFlags | log.Lmsgprefix | log.LUTC
		fatal = ioutil.Discard
		error = ioutil.Discard
		info  = ioutil.Discard
		debug = ioutil.Discard
	)

	if envLogLevel >= 1 {
		fatal = out
	}
	if envLogLevel >= 2 {
		error = out
	}
	if envLogLevel >= 5 {
		info = out
	}
	if envLogLevel >= 6 {
		debug = out
	}
	return &Loggers{
		LoggersArray: LoggersArray{
			Log_level_fatal: Logger{Logger: log.New(fatal, LogLevelText[Log_level_fatal]+": ", flags)},
			Log_level_error: Logger{Logger: log.New(error, LogLevelText[Log_level_error]+": ", flags)},
			Log_level_info:  Logger{Logger: log.New(info, LogLevelText[Log_level_info]+": ", flags)},
			Log_level_debug: Logger{Logger: log.New(debug, LogLevelText[Log_level_debug]+": ", flags)}},
		envLogLevel:  envLogLevel,
		Entry:        Fields{},
		RWMutex:      sync.RWMutex{},
		CircuitBreak: validateCircuitBreak(circuitBreak),
	}
}

// Converts a struct to key/value pairs
func ConvertStruct(s interface{}) Fields {
	in, _ := json.Marshal(s)

	var inMap Fields
	_ = json.Unmarshal(in, &inMap)

	return inMap
}

//Converts variadic args into key/value pairs
func ConvertArgs(args ...interface{}) Fields {
	if len(args)%2 != 0 {
		panic("ConvertArgs: Parameters must be an array of ordered pairs")
	}
	m := make(Fields)
	for i := 0; i < len(args); i = i + 2 {
		k := fmt.Sprintf("%v", args[i])
		v := args[i+1]
		m[k] = v
	}
	return m
}

func (l *Loggers) SetFields(fields Fields) *Loggers {
	l.Lock()
	defer l.Unlock()

	l.Entry = Fields{}
	for k, v := range fields {
		l.Entry[k] = v
	}
	return l
}

// Creates new LogEntry with key/value pairs specified in f. Copies any existing
// fields out of Loggers
func (l *Loggers) WithFields(fields Fields) *LogEntry {
	l.Lock()
	defer l.Unlock()

	var entry = Fields{}
	for k, v := range l.Entry {
		entry[k] = v
	}
	for k, v := range fields {
		entry[k] = v
	}

	return &LogEntry{
		LoggersArray: l.LoggersArray,
		Entry:        entry,
		RWMutex:      sync.RWMutex{},
	}
}

// Adds key/values in fields to LogEntry
func (e *LogEntry) WithFields(fields Fields) *LogEntry {
	e.Lock()
	defer e.Unlock()

	for k, v := range fields {
		e.Entry[k] = v
	}
	return e
}

// Creates new LogEntry with key/value pairs from s.
func (l *Loggers) WithStruct(s interface{}) *LogEntry {
	l.Lock()
	defer l.Unlock()

	var entry = Fields{}
	for k, v := range l.Entry {
		entry[k] = v
	}
	for k, v := range ConvertStruct(s) {
		entry[k] = v
	}

	return &LogEntry{
		LoggersArray: l.LoggersArray,
		Entry:        entry,
		RWMutex:      sync.RWMutex{},
	}
}

// Converts struct to key/value pairs and adds these to LogEntry
func (l *LogEntry) WithStruct(s interface{}) *LogEntry {
	l.Lock()
	defer l.Unlock()

	for k, v := range ConvertStruct(s) {
		l.Entry[k] = v
	}
	return l
}

func WithDefaultFields(args ...interface{}) Fields {
	var fields = Fields{}
	fields["msg"] = Fields{}

	var note interface{}
	if len(args) == 1 {
		switch args[0].(type) {
		case string, int, float32, float64, bool, uint:
			note = Fields{"note": fmt.Sprintf("%v", args[0])}
		case interface{}:
			note = ConvertStruct(args[0])
		}
	}
	if len(args) > 1 {
		note = ConvertArgs(args...)
	}

	if note != nil {
		fields["msg"] = note
	}

	return fields
}

func MergeFields(entry Fields, args ...interface{}) Fields {
	if len(args) == 0 {
		return entry
	}

	var fields = Fields{}
	for k, v := range WithDefaultFields(args...) {
		fields[k] = v
	}

	for k, v := range entry {
		fields[k] = v
	}
	return fields
}

// Logs only what is specified in args
func (e *Loggers) Info(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_info]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
	e.incLogCount(1)
}

// Logs only what is specified in args
func (e *Loggers) Debug(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_debug]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
	e.incLogCount(1)
}

// Logs only what is specified in args
func (e *Loggers) Fatal(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_fatal]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
	e.incLogCount(1)
}

// Logs only what is specified in args
func (e *Loggers) Error(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_error]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
	e.incLogCount(1)
}

// Logs args AND existng key/values
func (e *LogEntry) Info(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_info]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
}

// Logs args AND existng key/values
func (e *LogEntry) Debug(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_debug]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
}

// Logs args AND existng key/values
func (e *LogEntry) Fatal(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_fatal]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
}

// Logs args AND existng key/values
func (e *LogEntry) Error(args ...interface{}) {
	e.Lock()
	defer e.Unlock()

	logger := e.LoggersArray[Log_level_error]
	message, _ := json.Marshal(MergeFields(e.Entry, args...))
	logger.Print(string(message))
}
