package main

import "time"

type Benchmark struct {
	messageSize int64
	n           int64
	message     []byte
	duration    time.Duration
	algorithm   string
}
