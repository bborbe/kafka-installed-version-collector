package cron

import (
	"context"
)

type Cron interface {
	Run(ctx context.Context) error
}

type action func(ctx context.Context) error
