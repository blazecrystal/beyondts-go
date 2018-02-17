package utils

import "time"

func Sleep(milliSec int) {
	time.Sleep(time.Duration(milliSec) * time.Millisecond)
}
