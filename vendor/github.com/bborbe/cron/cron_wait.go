package cron

import (
	"context"
	"time"

	"github.com/golang/glog"
)

type cronWait struct {
	action action
	wait   time.Duration
}

func NewWaitCron(
	wait time.Duration,
	action action,
) *cronWait {
	c := new(cronWait)
	c.action = action
	c.wait = wait
	return c
}

func (c *cronWait) Run(ctx context.Context) error {
	for {
		glog.V(4).Infof("run cron action started")
		if err := c.action(ctx); err != nil {
			glog.V(2).Infof("action failed -> exit")
			return err
		}
		select {
		case <-ctx.Done():
			glog.V(2).Infof("context done -> exit")
			return nil
		case <-c.sleep():
			glog.V(4).Infof("sleep completed")
		}
	}
	return nil
}

func (c *cronWait) sleep() <-chan time.Time {
	glog.V(0).Infof("sleep for %v", c.wait)
	return time.After(c.wait)
}
