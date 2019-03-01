# alog
* **alog** is a golang implementation of levelled logging.  
* It wraps the log package from the standard library and provides methods for logging at the following levels (which do not get written to the log if the log level is lower than the configured value in alog.conf) :
- TRACE
- DEBUG
- INFO
- WARN
- ERROR
- CRITICAL
* The biggest advantage this package offers is that there is no per-call expensive check of the level of logging. This is all determined in the *init* function of the package which creates pointers to the correct function which is to be used for logging

## Getting It
go get -u "github.com/en-vee/alog"

## Configuration (alog.conf)
```shell
alog {
    fileName = "C://Temp//axlrate1.log" # Name, including the full path, of the file to which the log is to be written
    logLevel = "TRACE" # Valid Values = TRACE|DEBUG|INFO|WARN|ERROR|CRITICAL
}
```
* The config options in the above file are self-explanatory

## Usage
* Import the alog package
```go
import "github.com/en-vee/alog"
```
* Create a configuration file as shown above
* Log at the desired level
```go
alog.Info("This is an INFO message")
```

## How it Works
* At startup (in the package init function), it first looks for an alog.conf in the current directory.  
* If not found, it then checks if there is such a config file as indicated in the location in the environment variable ```ALOG_CONF_DIR```  
* Finally, if alog.conf is not found in any of the above locations, it uses STDOUT as the logger destination.  
* Once the package initialiazation is complete, alog provides methods to log at one of the desired levels as mentioned earlier. * * The method names follow the levels and accept arguments in Printf style.  
* For example : ```alog.Debug(msg string, i ...interface{})```  
* If the log level specified in the conf file is DEBUG, any messages of level lower than DEBUG will not be written to the log file.

* Thus, the methods exposed by the *alog* package are :
```go
alog.Trace(string, ...interface{})
alog.Debug(string, ...interface{})
alog.Info(string, ...interface{})
alog.Warn(string, ...interface{})
alog.Error(string, ...interface{})
alog.Critical(string, ...interface{})
```
* Sample Log Message
```shell
2018/11/07 18:03:25 [ERROR]      - This is an ERROR message.
```


## Other package(s) used
github.com/en-vee/aconf - golang based library for parsing/unmarshaling HOCON files
