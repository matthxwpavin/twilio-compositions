package recording

import "time"

type RecordingInstance struct {
	AccountSid      string    `json:"account_sid"`
	Status          string    `json:"status"`
	DateCreated     time.Time `json:"date_created"`
	Sid             string    `json:"sid"`
	SourceSid       string    `json:"source_sid"`
	Size            int       `json:"size"`
	URL             string    `json:"url"`
	Type            string    `json:"type"`
	Duration        int       `json:"duration"`
	ContainerFormat string    `json:"container_format"`
	Codec           string    `json:"codec"`
	TrackName       string    `json:"track_name"`
	Offset          int       `json:"offset"`
	GroupingSids    struct {
		RoomSid string `json:"room_sid"`
	} `json:"grouping_sids"`
	Links struct {
		Media string `json:"media"`
	} `json:"links"`
}

type RecordingList struct {
	Recording []RecordingInstance `json:"recording"`
	Meta      struct {
		Page            int     `json:"page"`
		PageSize        int     `json:"page_size"`
		FirstPageUrl    string  `json:"first_page_url"`
		PreviousPageUrl *string `json:"previous_page_url"`
		Url             string  `json:"url"`
		NextPageUrl     *string `json:"next_page_url"`
		Key             string  `json:"key"`
	} `json:"meta"`
}
