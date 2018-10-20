package run

import (
	"context"

	"sync"

	"github.com/bborbe/run/errors"
	"github.com/golang/glog"
)

type RunFunc func(context.Context) error

// CancelOnFirstFinish executes all given functions. After the first function finishes, any remaining functions will be canceled.
func CancelOnFirstFinish(ctx context.Context, runners ...RunFunc) error {
	if len(runners) == 0 {
		glog.V(2).Infof("nothing to run")
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errors := make(chan error)
	for _, runner := range runners {
		run := runner
		go func() {
			errors <- run(ctx)
		}()
	}
	select {
	case err := <-errors:
		return err
	case <-ctx.Done():
		glog.V(1).Infof("context canceled return")
		return nil
	}
}

// CancelOnFirstError executes all given functions. When a function encounters an error all remaining functions will be canceled.
func CancelOnFirstError(ctx context.Context, runners ...RunFunc) error {
	if len(runners) == 0 {
		glog.V(2).Infof("nothing to run")
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errors := make(chan error)
	done := make(chan struct{})
	var wg sync.WaitGroup
	for _, runner := range runners {
		wg.Add(1)
		run := runner
		go func() {
			defer wg.Done()
			if result := run(ctx); result != nil {
				errors <- result
			}
		}()
	}
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	select {
	case err := <-errors:
		return err
	case <-done:
		return nil
	case <-ctx.Done():
		glog.V(1).Infof("context canceled return")
		return nil
	}
}

// All executes all given functions. Errors are wrapped into one aggregate error.
func All(ctx context.Context, runners ...RunFunc) error {
	if len(runners) == 0 {
		glog.V(2).Infof("nothing to run")
		return nil
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errorChannel := make(chan error, len(runners))
	var errorWg sync.WaitGroup
	var runWg sync.WaitGroup
	var errs []error
	errorWg.Add(1)
	go func() {
		defer errorWg.Done()
		for err := range errorChannel {
			errs = append(errs, err)
		}
	}()
	for _, runner := range runners {
		run := runner
		runWg.Add(1)
		go func() {
			defer runWg.Done()
			if err := run(ctx); err != nil {
				errorChannel <- err
			}
		}()
	}
	glog.V(4).Infof("wait on runs")
	runWg.Wait()
	close(errorChannel)
	glog.V(4).Infof("wait on error collect")
	errorWg.Wait()
	glog.V(4).Infof("run all finished")
	if len(errs) > 0 {
		glog.V(4).Infof("found %d errors", len(errs))
		return errors.New(errs...)
	}
	glog.V(4).Infof("finished without errors")
	return nil
}

func Sequential(ctx context.Context, funcs ...RunFunc) (err error) {
	for _, fn := range funcs {
		select {
		case <-ctx.Done():
			glog.V(1).Infof("context canceled return")
			return nil
		default:
			if err = fn(ctx); err != nil {
				return
			}
		}
	}
	return
}
