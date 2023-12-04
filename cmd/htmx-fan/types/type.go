package types

import "time"

type State struct {
	Temperature float32
	FanSpeed    uint8
	Timestamp   string
	Realtime    time.Time
}
