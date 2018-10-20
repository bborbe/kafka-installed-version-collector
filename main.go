// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The kafka-version-collector collects available versions of software and publish it to a topic.
package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/kafka-installed-version-collector/version"
	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := &version.App{}
	flag.DurationVar(&app.Wait, "wait", time.Hour, "time to wait before next version collect")
	flag.IntVar(&app.Port, "port", 9002, "port to listen")
	flag.StringVar(&app.AppName, "app-name", "", "app name")
	flag.StringVar(&app.AppRegex, "app-regex", "", "app regex")
	flag.StringVar(&app.AppUrl, "app-url", "", "app url")
	flag.StringVar(&app.KafkaBrokers, "kafka-brokers", "", "kafka brokers")
	flag.StringVar(&app.KafkaTopic, "kafka-topic", "", "kafka topic")
	flag.StringVar(&app.SchemaRegistryUrl, "kafka-schema-registry-url", "", "kafka schema registry url")

	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.V(0).Infof("Parameter AppName: %s", app.AppName)
	glog.V(0).Infof("Parameter AppRegex: %s", app.AppRegex)
	glog.V(0).Infof("Parameter AppUrl: %s", app.AppUrl)
	glog.V(0).Infof("Parameter KafkaBrokers: %s", app.KafkaBrokers)
	glog.V(0).Infof("Parameter KafkaSchemaRegistryUrl: %s", app.SchemaRegistryUrl)
	glog.V(0).Infof("Parameter KafkaTopic: %s", app.KafkaTopic)
	glog.V(0).Infof("Parameter Port: %d", app.Port)
	glog.V(0).Infof("Parameter Wait: %v", app.Wait)

	if err := app.Validate(); err != nil {
		glog.Exitf("validate app failed: %v", err)
	}

	ctx := contextWithSig(context.Background())

	glog.V(1).Infof("app started")
	if err := app.Run(ctx); err != nil {
		glog.Exitf("app failed: %+v", err)
	}
	glog.V(1).Infof("app finished")
}

func contextWithSig(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}
