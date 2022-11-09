package userpref

type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type UserPreferences struct {
	Language  LabelValue `json:"language"`
	Currency  LabelValue `json:"currency"`
	Sounds    Sounds     `json:"sounds"`
	Behaviors `json:"behaviors"`
}

type Behaviors struct {
	Notifications bool `json:"notifications"`

	RetryCheckouts bool `json:"retryCheckouts"`
	MonitorPooling bool `json:"monitorPooling"`

	AutoClickHarvester   bool `json:"autoClick"`
	AutoLaunchCollective bool `json:"autoLaunch"`
}
