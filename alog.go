// Package log provides levelled logging functionality.
// It also allows one to configure the destination of the logs. The default is stdout.
// It looks for a logger configuration file alog.conf first in the current directory and then in the directory defined by ALOG_CONF_DIR
// If it does not find a logger configuration file in any of these locations, then it uses STDOUT as the logging destination
package alog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/en-vee/aconf"
)

/*
log package provides functions to log messages at various log levels :
TRACE
DEBUG
INFO
WARNING
ERROR
CRITICAL
*/

var (
	loggerConfigFileName = "alog.conf"
	logDestination       = os.Stdout
	logLevel             LogLevel
)

// LogLevel is the type used to specify the log level
type LogLevel uint8

// The log level constants specify the log levels which can be accepted
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
	// Is alog.conf present in local folder ?
	// 	If yes,
	// 	Instantiate an io.Reader using the file name alog.conf
	// If not, then check if ALOG_CONF_DIR environment variable has been defined.
	// 		If yes, then attempt to create an io.Reader from alog.conf.
	// If reader is still nil, then just set destination output to stdout

	var ok bool
	configParser := &aconf.HoconParser{}
	alogConfig := &struct {
		Alog struct {
			FileName string `hocon:"fileName"`
			LogLevel string `hocon:"logLevel"`
		} `hocon:"alog"`
	}{}

	// Select logger config file, giving priority to local alog.conf
	if logConfDir, ok := os.LookupEnv("ALOG_CONF_DIR"); ok && !fileExists("alog.conf") {
		loggerConfigFileName = fmt.Sprintf("%s%c%s", logConfDir, os.PathSeparator, "alog.conf")
	}

	if reader, err := os.Open(loggerConfigFileName); err == nil {
		if err := configParser.Parse(reader, alogConfig); err == nil {
			if len(alogConfig.Alog.FileName) != 0 {
				logDestination, err = os.OpenFile(alogConfig.Alog.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Fprintf(os.Stderr, "alog: unable to open log file : "+alogConfig.Alog.FileName+". Error : "+err.Error()+"\n")
					fmt.Fprintf(os.Stderr, "alog: using STDOUT for logging\n")
					logDestination = os.Stdout
				}
			}

			if logLevel, ok = logStringToIntLevelMap[alogConfig.Alog.LogLevel]; !ok {
				fmt.Println("alog: invalid log level specified :", alogConfig.Alog.LogLevel, "Using default level of TRACE")
			}
		}
	}

	SetLogLevel(logLevel)
	log.SetOutput(logDestination)
	//log.SetPrefix(logLevelIntToStringMap[logLevel] + " - ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

func setLogLevel(level LogLevel) {
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

func SetLogLevel(level LogLevel) error {

	if level > CRITICAL {
		return &InvalidLogLevelError{level}
	}

	setLogLevel(level)

	return nil
}

func SetLogDestination(w io.Writer) {
	singleTon.Do(func(){
		log.SetOutput(w)
	})
}

// noOpLogMsg is just an empty (No Operation) implementation which does nothing.
// It is needed with full signature so that it can be set into a function value which is compatible with the actual log.Printf method
func noOpLogMsg(level LogLevel, msg string, objs ...interface{}) {}

// logMsg performs actual logging to a destination when used as a function value for a specific log level
func logMsg(level LogLevel, msg string, objs ...interface{}) {

	var sb strings.Builder

	sb.WriteString("- ")
	sb.WriteString(logLevelIntToStringMap[level])
	sb.WriteString("- ")
	sb.WriteString(msg)

	m := sb.String()

	if len(objs) > 0 {
		log.Printf(m, objs...)
	} else {
		log.Printf(m)
	}
	//log.Printf("%-12s - %s\n", logLevelIntToStringMap[level], msg)
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
