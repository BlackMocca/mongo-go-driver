// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build logrus

package main

import (
	"context"
	"log"

	"github.com/BlackMocca/mongo-go-driver/bson"
	"github.com/BlackMocca/mongo-go-driver/mongo"
	"github.com/BlackMocca/mongo-go-driver/mongo/options"
	"github.com/bombsimon/logrusr/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a new logrus logger instance.
	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.DebugLevel)

	// Create a new sink for logrus using "logrusr".
	sink := logrusr.New(logger).GetSink()

	// Create a client with our logger options.
	loggerOptions := options.
		Logger().
		SetSink(sink).
		SetMaxDocumentLength(25).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	clientOptions := options.
		Client().
		ApplyURI("mongodb://localhost:27017").
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	defer client.Disconnect(context.TODO())

	// Make a database request to test our logging solution.
	coll := client.Database("test").Collection("test")

	_, err = coll.InsertOne(context.TODO(), bson.D{{"Alice", "123"}})
	if err != nil {
		log.Fatalf("InsertOne failed: %v", err)
	}
}
