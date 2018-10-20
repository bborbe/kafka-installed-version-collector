// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bborbe/cron"
	"github.com/bborbe/run"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seibert-media/go-kafka/schema"
)

type App struct {
	AppName           string
	AppRegex          string
	AppUrl            string
	KafkaBrokers      string
	KafkaTopic        string
	Port              int
	SchemaRegistryUrl string
	Wait              time.Duration
}

func (a *App) Validate() error {
	if a.Port <= 0 {
		return errors.New("Port missing")
	}
	if a.Wait <= 0 {
		return errors.New("Wait missing")
	}
	if a.KafkaBrokers == "" {
		return errors.New("KafkaBrokers missing")
	}
	if a.KafkaTopic == "" {
		return errors.New("KafkaTopic missing")
	}
	if a.SchemaRegistryUrl == "" {
		return errors.New("SchemaRegistryUrl missing")
	}
	if a.AppName == "" {
		return errors.New("AppName missing")
	}
	if a.AppRegex == "" {
		return errors.New("AppRegex missing")
	}
	if a.AppUrl == "" {
		return errors.New("AppUrl missing")
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	return run.CancelOnFirstFinish(
		ctx,
		a.runCron,
		a.runHttpServer,
	)
}

func (a *App) runCron(ctx context.Context) error {
	fetcher := &Fetcher{
		HttpClient: http.DefaultClient,
		Url:        a.AppUrl,
		App:        a.AppName,
		Regex:      a.AppRegex,
	}
	sender := &Sender{
		KafkaTopic:   a.KafkaTopic,
		KafkaBrokers: a.KafkaBrokers,
		SchemaRegistry: &schema.Registry{
			HttpClient:        http.DefaultClient,
			SchemaRegistryUrl: a.SchemaRegistryUrl,
		},
	}

	cronJob := cron.NewWaitCron(
		a.Wait,
		func(ctx context.Context) error {
			applicationVersionInstalled, err := fetcher.Fetch(ctx)
			if err != nil {
				return err
			}
			return sender.Send(ctx, *applicationVersionInstalled)
		},
	)
	return cronJob.Run(ctx)
}

func (a *App) runHttpServer(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: promhttp.Handler(),
	}
	go func() {
		select {
		case <-ctx.Done():
			server.Shutdown(ctx)
		}
	}()
	return server.ListenAndServe()
}
