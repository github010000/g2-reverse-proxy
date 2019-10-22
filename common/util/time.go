package util

import "time"

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type clock struct{}

func NewClock() Clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}

func (c *clock) Since(tm time.Time) time.Duration {
	return time.Since(tm)
}

type brokenClock struct {
	margin time.Duration
}

func NewBrokenClock(margin time.Duration) Clock {
	return &brokenClock{margin}
}

func (c *brokenClock) Now() time.Time {
	now := time.Now()
	return now.Add(c.margin)
}

func (c *brokenClock) Since(tm time.Time) time.Duration {
	since := time.Since(tm)
	return since + c.margin
}