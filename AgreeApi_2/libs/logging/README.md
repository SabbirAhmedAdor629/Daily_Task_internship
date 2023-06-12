# Overview

This library provides some basic structured logging capabilities for programs written in Go. Each log event is sent to ``standard err``, ``standard output`` (or a specified file) in JSON format on a single line.

The library currently supports five log levels; ``NONE``, ``DEBUG``, ``INFO``, ``ERROR`` and ``FATAL``. 

``WARNING`` is not currently supported as warnings are usually ignored. And ``TRACE`` seems a little excessive.

To create a new Logger, an output channel and log level must be specified.

``` go
logger = NewLoggers(os.Stderr, log_level_info)
```

Creating a log event is as straightforward as calling the logger with a logging level; ``.Debug()``, ``.Info()``, ``.Error()`` or ``.Fatal()``.

``` go
// Do some stuff
logger.Info()

// Do some other stuff
logger.Debug()
```

This results in the following log entry. Note that the debug entry is not created as this exceeds ``log_level_info`` specified when the logger was created.

```
021/09/14 02:57:38 INFO: {}
```

The log level functions (``.Info()``, ``.Error()``, etc.) will output any specified parameters as part of the log entry.

``` go
// Create a new Logger
logger = NewLoggers(os.Stderr, log_level_debug)

// Do some stuff
logger.Info("Something normal happend here.")

// So some more stuff
logger.Debug("Something unexpected happened here.")
```

```
021/09/14 02:57:38 INFO: {"msg":{"note": "Something normal happend here."}}
...
021/09/14 02:57:38 DEBUG: {"msg":{"note": "Something unexpected happened here."}}
```

A logger can also accept a number of fields that, once set, are output on subsequent log entries. 

``` go
logger = NewLoggers(os.Stderr, log_level_info).
		SetFields(Fields{
			"app_name": "test_application", 
			"job_id":   "123456"})

// Suff happened
logger.Info("Thing #1")

// More stuff happened
logger.info("Thing #2")
```

```
021/09/14 02:57:38 INFO: {"app_name":"test_application","job_id":"123456","msg":{"note": "Thing #1"}}
...
021/09/14 02:57:38 INFO: {"app_name":"test_application","job_id":"123456","msg":{"note": "Thing #2"}}
```

It is also possible to create log "entries" once, and output these several times. This can be useful when creating log entries of various "types";

``` go
logger = NewLoggers(os.Stderr, log_level_info).
		SetFields(Fields{
			"app_name": "test_application",
			"job_id":   "123456"})
dbEntry := logger.NewEntry(log_database, "database log")

// Stuff working as expected
dbEntry.Info("UPDATE successful")

// Stuff all stuffed up
dblogEntry.Fatal("Database Connection Error")
```

```
021/09/14 02:57:38 INFO: {"app_name":"test_application","job_id":"123456","msg_code":105, "msg_text":"database log",msg":{"note": "UPDATE successful"}}
...
021/09/14 02:57:38 FATAL: {"app_name":"test_application","job_id":"123456","msg_code":105, "msg_text":"database log","msg":{"note": "Database Connection Error"}}
```

## Installation

1. Enter the following at the command line in your project's root directory. Where ``../logging`` is the relative path to the logging source code. 
``` 
$ replace influencemobile.com/logging => ../logging
$ go mod tidy
$ go get
```
