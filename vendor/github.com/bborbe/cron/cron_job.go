package cron

import (
	"context"
	"time"

	"github.com/golang/glog"
)

type cronJob struct {
	oneTime    bool
	expression string
	wait       time.Duration
	action     func(ctx context.Context) error
}

func NewCronJob(
	oneTime bool,
	expression string,
	wait time.Duration,
	action func(ctx context.Context) error,
) *cronJob {
	return &cronJob{
		oneTime:    oneTime,
		expression: expression,
		wait:       wait,
		action:     action,
	}
}

func (c *cronJob) Run(ctx context.Context) error {
	var runner Cron
	if c.oneTime {
		glog.V(2).Infof("create one-time cron")
		runner = NewOneTimeCron(c.action)
	} else if len(c.expression) > 0 {
		glog.V(2).Infof("create cron with expression %s", c.expression)
		runner = NewExpressionCron(
			c.expression,
			c.action,
		)
	} else {
		glog.V(2).Infof("create cron with wait %v", c.wait)
		runner = NewWaitCron(
			c.wait,
			c.action,
		)
	}
	return runner.Run(ctx)
}
