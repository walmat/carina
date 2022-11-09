package main

import (
	"golang.org/x/sys/windows"
)

var (
	winNm               = windows.NewLazyDLL("winmm.dll")
	procTimeBeginPeriod = winNm.NewProc("timeBeginPeriod")
)

func main() {
	_, _, _ = procTimeBeginPeriod.Call(1)
	sharedMain()
}
