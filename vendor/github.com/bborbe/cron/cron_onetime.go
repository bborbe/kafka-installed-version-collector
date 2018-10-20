package cron

import (
	"context"

	"github.com/golang/glog"
)

type cronOneTime struct {
	action action
}

func NewOneTimeCron(
	action action,
) *cronOneTime {
	c := new(cronOneTime)
	c.action = action
	return c
}

func (c *cronOneTime) Run(ctx context.Context) error {
	glog.V(4).Infof("run cron action started")
	if err := c.action(ctx); err != nil {
		glog.V(2).Infof("action failed -> exit")
		return err
	}
	return nil
}
