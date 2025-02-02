// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package unified

import (
	"github.com/BlackMocca/mongo-go-driver/internal/logger"
)

// orderedLogMessage is a logMessage with an "order" field representing the
// order in which the log message was observed.
type orderedLogMessage struct {
	*logMessage
	order int
}

// Logger is the Sink used to captured log messages for logger verification in
// the unified spec tests.
type Logger struct {
	lastOrder int
	logQueue  chan orderedLogMessage
	bufSize   int
}

func newLogger(olm *observeLogMessages, bufSize int) *Logger {
	if olm == nil {
		return nil
	}

	return &Logger{
		lastOrder: 1,
		logQueue:  make(chan orderedLogMessage, bufSize),
		bufSize:   bufSize,
	}
}

// Info implements the logger.Sink interface's "Info" method for printing log
// messages.
func (log *Logger) Info(level int, msg string, args ...interface{}) {
	if log.logQueue == nil {
		return
	}

	defer func() { log.lastOrder++ }()

	// If the order is greater than the buffer size, we must return. This
	// would indicate that the logQueue channel has been closed.
	if log.lastOrder > log.bufSize {
		return
	}

	// Add the Diff back to the level, as there is no need to create a
	// logging offset.
	level = level + logger.DiffToInfo

	logMessage, err := newLogMessage(level, msg, args...)
	if err != nil {
		panic(err)
	}

	// Send the log message to the "orderedLogMessage" channel for
	// validation.
	log.logQueue <- orderedLogMessage{
		order:      log.lastOrder + 1,
		logMessage: logMessage,
	}

	// If the order has reached the buffer size, then close the channe.
	if log.lastOrder == log.bufSize {
		close(log.logQueue)
	}
}

// Error implements the logger.Sink interface's "Error" method for printing log
// errors. In this case, if an error occurs we will simply treat it as
// informational.
func (log *Logger) Error(err error, msg string, args ...interface{}) {
	args = append(args, "error", err)
	log.Info(int(logger.LevelInfo), msg, args)
}
