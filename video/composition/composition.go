package composition

import "time"

type Composition struct {
	AccountSid           string     `json:"account_sid"`
	AudioSources         []string   `json:"audio_sources"`
	AudioSourcesExcluded []string   `json:"audio_sources_excluded"`
	Bitrate              int        `json:"bitrate"`
	DateCompleted        time.Time  `json:"date_completed"`
	DateCreated          time.Time  `json:"date_created"`
	DateDeleted          *time.Time `json:"date_deleted"`
	Duration             int        `json:"duration"`
	Format               string     `json:"format"`
	Links                struct {
		Media string `json:"media"`
	} `json:"links"`
	Resolution  string                 `json:"resolution"`
	RoomSid     string                 `json:"room_sid"`
	Sid         string                 `json:"sid"`
	Size        int                    `json:"size"`
	Status      string                 `json:"status"`
	Trim        bool                   `json:"trim"`
	URL         string                 `json:"url"`
	VideoLayout map[string]interface{} `json:"video_layout"`
}

type CompositionList struct {
	Compositions []Composition `json:"compositions"`
	Meta         Meta          `json:"meta"`
}

type Meta struct {
	FirstPageUrl    string  `json:"first_page_url"`
	Key             string  `json:"key"`
	NextPageUrl     *string `json:"next_page_url"`
	Page            uint64  `json:"page"`
	PageSize        uint64  `json:"page_size"`
	PreviousPageUrl *string `json:"previous_page_url"`
	URL             string  `json:"url"`
}

const (
	StatusEnqueued   = "enqueued"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusDeleted    = "deleted"
	StatusFailed     = "failed"
)
