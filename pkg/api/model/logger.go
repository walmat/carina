package model

import "cloud.google.com/go/logging"

type Severity logging.Severity

const (
	// Default means the log entry has no assigned severity level.
	Default = Severity(logging.Default)
	// Debug means debug or trace information.
	Debug = Severity(logging.Debug)
	// Info means routine information, such as ongoing status or performance.
	Info = Severity(logging.Info)
	// Notice means normal but significant events, such as start up, shut down, or configuration.
	Notice = Severity(logging.Notice)
	// Warning means events that might cause problems.
	Warning = Severity(logging.Warning)
	// Error means events that are likely to cause problems.
	Error = Severity(logging.Error)
	// Critical means events that cause more severe problems or brief outages.
	Critical = Severity(logging.Critical)
	// Alert means a person must take an action immediately.
	Alert = Severity(logging.Alert)
	// Emergency means one or more systems are unusable.
	Emergency = Severity(logging.Emergency)
)

type LoggerEntry struct {
	Severity Severity          `json:"severity"`
	Payload  interface{}       `json:"payload"`
	Labels   map[string]string `json:"labels"`
	File     *LoggerFile       `json:"file"`
}

type LoggerFile struct {
	Name string `json:"file"`
	Line int    `json:"line"`
}
