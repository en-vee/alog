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
fileName = "myapp.log"
filePath = "/var/log/myapp"
# Valid Values = TRACE|DEBUG|INFO|WARN|ERROR|CRITICAL
logLevel = "INFO"
```
The config options in the above file are self-explanatory.

## How it Works
* At startup (in the package init function), it first looks for an alog.conf in the current directory.  
* If not found, it then checks if there is such a config file as indicated in the location in the environment variable ```ALOG_CONF_DIR```  
* Finally, if alog.conf is not found in any of the above locations, it uses STDOUT as the logger destination.  
* Once the package initialiazation is complete, alog provides methods to log at one of the desired levels as mentioned earlier. * * The method names follow the levels and accept arguments in Printf style.  
* For example : ```alog.Debug(msg string, i ...interface{})```  
* If the log level specified in the conf file is DEBUG, any messages of level lower than DEBUG will not be written to the log file.
* This is accomplished by use of an array of function values with each element in the array pointing to the log function if the level is higher/equal to the configured level OR to a No-Op function if the level is lower.   
* A snippet from the code should make it clearer :
```
var logFuncsSlice = []logFuncType{noOpLogMsg, noOpLogMsg, noOpLogMsg, noOpLogMsg, noOpLogMsg, logMsg}
```

where noOpLogMsg is a method which does nothing  
and logMsg is just a wrapper around the standard library log.Printf  
* The index of each element in the array represents the log level (0=TRACE,1=DEBUG, etc.)
* Thus, the methods exposed by the *alog* package are :
```shell
alog.Trace(string, ...interface{})
alog.Debug(string, ...interface{})
alog.Info(string, ...interface{})
alog.Warn(string, ...interface{})
alog.Error(string, ...interface{})
alog.Critical(string, ...interface{})
```
* Sample Log Message
```shell
018/11/07 18:03:25 [ERROR]      - This is an ERROR message.
```


## Other package(s) used
Hashicorp HCL
golang based library for levelled logging
