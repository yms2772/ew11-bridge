package ew11

type DeviceBaseConfig struct {
	Platform            string         `json:"platform"`
	Name                string         `json:"name"`
	UniqueID            string         `json:"unique_id"`
	AvailabilityTopic   string         `json:"availability_topic"`
	PayloadAvailable    string         `json:"payload_available"`
	PayloadNotAvailable string         `json:"payload_not_available"`
	Device              map[string]any `json:"device"`
}
