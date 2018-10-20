package cron

import (
	"context"

	"github.com/golang/glog"
	robfig_cron "github.com/robfig/cron"
)

type cronExpression struct {
	ext    string
	action action
}

func NewExpressionCron(
	ext string,
	action action,
) *cronExpression {
	c := new(cronExpression)
	c.ext = ext
	c.action = action
	return c
}

func (c *cronExpression) Run(ctx context.Context) error {
	glog.V(4).Infof("register cron actions")
	errChan := make(chan error)
	cron := robfig_cron.New()
	cron.Start()
	defer cron.Stop()
	cron.AddFunc(c.ext, func() {
		glog.V(4).Infof("run cron action started")
		if err := c.action(ctx); err != nil {
			glog.V(2).Infof("action failed -> exit")
			errChan <- err
		}
		glog.V(4).Infof("run cron action finished")
	})
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}
