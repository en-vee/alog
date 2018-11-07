// Package log provides levelled logging functionality.
// It also allows one to configure the destination of the logs. The default is stdout.
// It looks for a logger configuration file alog.conf first in the current directory and then in the directory defined by AXLRATE_LOGGER_CONF_DIR
// If it does not find a logger configuration file in any of these locations, then it uses STDOUT as the logging destination
package alog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/mapstructure"
)

/*
log package provides functions to log messages at various log levels :
TRACE
DEBUG
INFO
WARNING
MINOR
MAJOR
CRITICAL
*/

var (
	loggerConfigFileName       = "alog.conf"
	logDestination             = os.Stdout
	logLevel             uint8 = 1
)

type LogLevel uint8

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
)

// Logging Function type
type logFuncType func(LogLevel, string, ...interface{})
type StringInterfaceMap map[string]interface{}

// Array containing function values which perform the actual logging
// Initialized with NoOp logger for all log levels except CRITICAL
var logFuncsSlice = []logFuncType{noOpLogMsg, noOpLogMsg, noOpLogMsg, noOpLogMsg, noOpLogMsg, logMsg}

var logLevelIntToStringMap = map[LogLevel]string{
	TRACE:    "[TRACE] ",
	DEBUG:    "[DEBUG] ",
	INFO:     "[INFO] ",
	WARN:     "[WARN] ",
	ERROR:    "[ERROR] ",
	CRITICAL: "[CRITICAL] ",
}

var logStringToIntLevelMap = map[string]LogLevel{
	"TRACE":    0,
	"DEBUG":    1,
	"INFO":     2,
	"WARN":     3,
	"ERROR":    4,
	"CRITICAL": 5,
}

type loggerConf struct {
	fileName string
	filePath string
	logLevel LogLevel
}

var theConfig loggerConf

func init() {

	//log.SetFlags(0)

	var useLocal bool

	// Read in axlrate-logger.conf from current folder
	// Slurp read whole file into buffer
	if fileContents, err := ioutil.ReadFile(loggerConfigFileName); err == nil {
		decodeConfFile(fileContents)
		useLocal = true
		logFileName := fmt.Sprintf("%s%c%s", theConfig.filePath, os.PathSeparator, theConfig.fileName)
		var err error
		logDestination, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Error creating log destination handle %v", err)
		}
	}
	// If not present, then attempt to read from Environment Variable AXLRATE_LOGGER_CONF_DIR
	if !useLocal {
		if logConfDir, ok := os.LookupEnv("AXLRATE_LOGGER_CONF_DIR"); !ok {
			if fileContents, err := ioutil.ReadFile(fmt.Sprintf("%s%c%s", logConfDir, os.PathSeparator, loggerConfigFileName)); err == nil {
				decodeConfFile(fileContents)
				logFileName := fmt.Sprintf("%s%c%s", theConfig.filePath, os.PathSeparator, theConfig.fileName)
				logDestination, _ = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			}
		}
	}
	SetLogLevel(theConfig.logLevel)
	log.SetOutput(logDestination)

}

func decodeConfFile(fileContents []byte) {
	var v interface{}
	//var theConfig loggerConf
	if err := hcl.Unmarshal(fileContents, &v); err == nil {
		if err := mapstructure.Decode(v, &theConfig); err == nil {
			if val, ok := v.(map[string]interface{})["fileName"].(string); ok {
				theConfig.fileName = val
			}
			if val, ok := v.(map[string]interface{})["filePath"].(string); ok {
				theConfig.filePath = val
			}
			if val, ok := (v.(map[string]interface{})["logLevel"]).(string); ok {
				theConfig.logLevel = logStringToIntLevelMap[val]
				//fmt.Println(theConfig.logLevel)
			}
		}
	}
	//fmt.Println(theConfig)
}

var singleTon sync.Once

// InvalidLogLevelError is used to indicate invalid log level
type InvalidLogLevelError struct {
	got LogLevel
}

// Stringer interface method(s)
func (ie *InvalidLogLevelError) String() string {
	return fmt.Sprintf("%d", ie.got)
}

// error interface method
func (ie *InvalidLogLevelError) Error() string {
	return fmt.Sprintf("Invalid Log Level : %v. Valid Values are TRACE|DEBUG|INFO|WARN|ERROR|CRITICAL", ie.got)
}

func SetLogLevel(level LogLevel) error {

	if level > CRITICAL {
		return &InvalidLogLevelError{level}
	}

	setLogLevel := func() {
		if level > CRITICAL {
			level = CRITICAL
		}

		for i := range logFuncsSlice {
			logFuncsSlice[i] = noOpLogMsg
		}

		// Level     => 0 1 2 3 4 5
		// Set/Unset => O O O X X X
		// For example, If level = 0, which is TRACE, then select slice from 0 through len(logFuncs)
		// If level = 1, which is DEBUG, then select slice from 1 through len(logFuncs)
		p := logFuncsSlice[level:]

		for i := range p {
			p[i] = logMsg
		}
	}

	singleTon.Do(setLogLevel)

	return nil
}

// noOpLogMsg is just an empty (No Operation) implementation which does nothing.
// It is needed with full signature so that it can be set into a function value which is compatible with the actual log.Printf method
func noOpLogMsg(level LogLevel, mmsg string, objs ...interface{}) {}

// logMsg performs actual logging to a destination when used as a function value for a specific log level
func logMsg(level LogLevel, msg string, objs ...interface{}) {
	//log.SetPrefix(logLevelIntToStringMap[level])
	log.Printf("%-12s - %s\n", logLevelIntToStringMap[level], msg)
}

func Trace(msg string, objs ...interface{}) {
	var level LogLevel = TRACE
	// Select Function based on level
	logFunc := logFuncsSlice[level]
	logFunc(level, msg, objs...)
}

func Debug(msg string, objs ...interface{}) {
	var level = DEBUG
	// Select Function based on slice
	logFunc := logFuncsSlice[level]
	logFunc(level, msg, objs...)
}

func Info(msg string, objs ...interface{}) {
	var level LogLevel = INFO
	// Select Function based on slice
	logFunc := logFuncsSlice[level]
	logFunc(level, msg, objs...)
}

func Warn(msg string, objs ...interface{}) {
	var level LogLevel = WARN
	// Select Function based on slice
	logFunc := logFuncsSlice[level]
	logFunc(level, msg, objs...)
}

func Error(msg string, objs ...interface{}) {
	var level LogLevel = ERROR
	// Select Function based on slice
	logFunc := logFuncsSlice[level]
	logFunc(level, msg, objs...)
}

func Critical(msg string, objs ...interface{}) {
	var level LogLevel = CRITICAL
	// Select Function based on slice
	logFunc := logFuncsSlice[level]
	logFunc(level, msg, objs...)
}
