package main

import "time"

type Benchmark struct {
	message   *[]byte
	duration  time.Duration
	algorithm string
}
