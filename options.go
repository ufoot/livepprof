// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
	"fmt"
	"time"
)

const (
	// defaultDelay when doing profiles
	defaultDelay = time.Minute
	// defaultLimit to not keep every single item in memory, only the
	// ones with the biggest numbers are kept (sorted before filtering out).
	defaultLimit = 20
)

type opts struct {
	filter      string
	errHandler  func(err error)
	delay       time.Duration
	limit       int
	disabled    bool
	enabledFunc func() bool
}

var defaultOpts = opts{
	delay: defaultDelay,
	limit: defaultLimit,
}

func (o *opts) enabled() bool {
	if o.enabledFunc == nil {
		// If no callback func is given, use the default
		// static flag. By default -> it's enabled.
		return !o.disabled
	}
	return o.enabledFunc()
}

// Option passed when creating the live profiler.
type Option func(o *opts) error

// WithFilter helps you spot where time is spent within your code
// by reporting functions which are in packages containing a given
// filter in their name. Typically if you are coding something in
// github.com/thisisme/supergolib/package1 you might want to
// give this "supergolib/package1" as you want functions in that
// part of the code to be reported. Not doing this, you might not
// have insightful reports as you may get a lot of entries about
// small "leaf" part of the code. If you do this the live profiler
// will aggregate all data until it finds a common parent in
// "supergolib/package1".
func WithFilter(filter string) Option {
	return func(o *opts) error {
		o.filter = filter
		return nil
	}
}

// WithErrorHandler allows custom handling of errors.
// This is useful as live profiler does thing in the background, the
// instanciation and start can not return all possible errors, so
// they need to be handled later, in a separate goroutine, as they happen.
func WithErrorHandler(errHandler func(err error)) Option {
	return func(o *opts) error {
		o.errHandler = errHandler
		return nil
	}
}

// WithDelay allows a custom delay to be used. Default is one minute.
func WithDelay(delay time.Duration) Option {
	return func(o *opts) error {
		if delay <= 0 {
			return fmt.Errorf("invalid delay: %s", delay.String())
		}
		o.delay = delay
		return nil
	}
}

// WithLimit allows a custom limit of displayed funcs to be used. Default is 20.
func WithLimit(limit int) Option {
	return func(o *opts) error {
		if limit <= 0 {
			return fmt.Errorf("invalid limit: %d", limit)
		}
		o.limit = limit
		return nil
	}
}

// WithEnabled allows you to enable/disable the profiler. If enabled is false,
// no profiling fill be done, even if the profiler is started.
func WithEnabled(enabled bool) Option {
	return func(o *opts) error {
		// internally, we use a disabled flag as this way, the default
		// is to have it enabled. But the option is "Enabled"
		// because "not disabled" is harder to understand in a public API.
		o.disabled = !enabled
		return nil
	}
}

// WithEnabledFunc allows you to enable/disable the profiler with a callback.
// This is useful if you want to enable/disable it on-the-fly without explicitly
// stopping or starting it. One use-case is if you have a dynamic configuration.
// This will actually poll the dynamic configuration, and enable/disable it.
// Typically interesting if you are concerned with the CPU the profiling is consuming,
// and/or if you want to get rid of any side effect.
func WithEnabledFunc(enabledFunc func() bool) Option {
	return func(o *opts) error {
		o.enabledFunc = enabledFunc
		return nil
	}
}
