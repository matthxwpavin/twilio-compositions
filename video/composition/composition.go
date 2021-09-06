package composition

import (
	"github.com/matthxwpavin/twilio-compositions/video"
	"time"
)

type CompStatus string

const (
	StatusEnqueued   CompStatus = "enqueued"
	StatusProcessing CompStatus = "processing"
	StatusCompleted  CompStatus = "completed"
	StatusDeleted    CompStatus = "deleted"
	StatusFailed     CompStatus = "failed"
)

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

var retFormatFunc = func(s string) Format {
	return &s
}

type Format *string

var (
	MP4  = retFormatFunc("mp4")
	WebM = retFormatFunc("webm")
)

type ComposeParams struct {

	// The SID of the Group Room with the media tracks to be used as composition sources.
	RoomSid string `form:"RoomSid"`

	// An object that describes the video layout of the composition hook in terms of regions.
	// See Specifying Video Layouts for more info.
	VideoLayout *video.VideoLayout `form:"-"`

	// An array of track names from the same group rooms to merge into the compositions
	// created by the composition hook. Can include zero or more track names.
	// A composition triggered by the composition hook includes all audio sources
	// specified in audio_sources except those specified in audio_sources_excluded.
	// The track names in this parameter can include an asterisk as a wild card character,
	// which matches zero or more characters in a track name. For example,
	// student* includes tracks named student as well as studentTeam.
	AudioSources *string `form:"AudioSources,omitempty"`

	// An array of track names to exclude. A composition triggered by the composition hook
	// includes all audio sources specified in audio_sources except
	// for those specified in audio_sources_excluded. The track names in this parameter can include
	// an asterisk as a wild card character, which matches zero or more characters in a track name.
	// For example, student* excludes student as well as studentTeam. This parameter can also be empty.
	AudioSourcesExcluded *string `form:"AudioSourcesExcluded,omitempty"`

	// A string that describes the columns (width) and rows (height)
	// of the generated composed video in pixels.
	// Defaults to 640x480. The string's format is {width}x{height} where:
	// 16 <= {width} <= 1280
	// 16 <= {height} <= 1280
	// {width} * {height} <= 921,600
	Resolution *string `form:"Resolution,omitempty"`

	// The container format of the media files used by the compositions
	// created by the composition hook. Can be: mp4 or webm and the default is webm.
	// If mp4 or webm, audio_sources must have one or more tracks and/or a video_layout element
	// must contain a valid video_sources list, otherwise an error occurs.
	Format Format `form:"Format,omitempty"`

	// The URL we should call using the status_callback_method to send status information
	// to your application on every composition event. I
	// f not provided, status callback events will not be dispatched.
	StatusCallback *string `form:"StatusCallback,omitempty"`

	// The HTTP method we should use to call status_callback. Can be: POST or GET and the default is POST.
	StatusCallbackMethod *string `form:"StatusCallbackMethod,omitempty"`

	// Whether to clip the intervals where there is no active media in the Compositions
	// triggered by the composition hook. The default is true.
	// Compositions with trim enabled are shorter when the Room is created and no Participant joins for a while
	// as well as if all the Participants leave the rooms and join later,
	// because those gaps will be removed. See Specifying Video Layouts for more info.
	Trim *bool `form:"Trim,omitempty"`
}

func (p *ComposeParams) GetVideoLayout() *video.VideoLayout {
	return p.VideoLayout
}

func (p *ComposeParams) GetResolution() *string {
	return p.Resolution
}

type GetParams struct {
	// Read only Composition resources with this status.
	// Can be: enqueued, processing, completed, deleted, or failed.
	Status *CompStatus `form:"Status,omitempty"`

	// Read only Composition resources created on or after this ISO 8601 date-time with time zone.
	DateCreatedAfter *string `form:"DateCreatedAfter,omitempty"`

	// Read only Composition resources created before this ISO 8601 date-time with time zone.
	DateCreatedBefore *string `form:"DateCreatedBefore,omitempty"`

	// Read only Composition resources with this Room SID.
	RoomSid *string `form:"RoomSid,omitempty"`
}
