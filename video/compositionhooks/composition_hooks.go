package compositionhooks

import (
	"twilio-compositions/video"
	"encoding/json"
	"errors"
	"github.com/ajg/form"
	"net/url"
	"time"
)

// More info https://www.twilio.com/docs/video/api/composition-hooks

type CompositionHooks struct {
	AccountSid           string                 `json:"account_sid"`
	Sid                  string                 `json:"sid"`
	FriendlyName         string                 `json:"friendly_name"`
	Enabled              bool                   `json:"enabled"`
	DateCreated          time.Time              `json:"date_created"`
	DateUpdated          *time.Time             `json:"date_updated"`
	AudioSources         []string               `json:"audio_sources"`
	AudioSourcesExcluded []string               `json:"audio_sources_excluded"`
	VideoLayout          map[string]interface{} `json:"video_layout"`
	Format               string                 `json:"format"`
	Trim                 bool                   `json:"trim"`
	URL                  string                 `json:"url"`
	Resolution           string                 `json:"resolution"`
	StatusCallbackMethod string                 `json:"status_callback_method"`
	StatusCallback       string                 `json:"status_callback"`
}

type CompositionHooksList struct {
	CompositionHooks []CompositionHooks `json:"composition_hooks"`
}

const (
	HD   = "1280x720"
	PAL  = "1024x576"
	VGA  = "640x480"
	CIF  = "320x240"
)

var retFormatFunc = func(s string) Format {
	return &s
}

type Format *string

var (
	MP4  = retFormatFunc("mp4")
	WebM = retFormatFunc("webm")
)

type CreateParams struct {
	// A descriptive string that you create to describe the resource.
	// It can be up to 100 characters long and it must be unique within the account.
	FriendlyName string `form:"FriendlyName"`

	// Whether the composition hook is active. When true,
	// the composition hook will be triggered for every completed Group Room in the account.
	// When false, the composition hook will never be triggered.
	Enabled *bool `form:"Enabled,omitempty"`

	// An object that describes the video layout of the composition hook in terms of regions.
	// See Specifying Video Layouts for more info.
	VideoLayout *video.VideoLayout `form:"-"`

	// An array of track names from the same group room to merge into the compositions
	// created by the composition hook. Can include zero or more track names.
	// A composition triggered by the composition hook includes all audio sources
	// specified in audio_sources except those specified in audio_sources_excluded.
	// The track names in this parameter can include an asterisk as a wild card character,
	// which matches zero or more characters in a track name. For example,
	// student* includes tracks named student as well as studentTeam.
	AudioSources []string `form:"AudioSources,omitempty"`

	// An array of track names to exclude. A composition triggered by the composition hook
	// includes all audio sources specified in audio_sources except
	// for those specified in audio_sources_excluded. The track names in this parameter can include
	// an asterisk as a wild card character, which matches zero or more characters in a track name.
	// For example, student* excludes student as well as studentTeam. This parameter can also be empty.
	AudioSourcesExcluded []string `form:"AudioSourcesExcluded,omitempty"`

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
	StatusCallBack *string `form:"StatusCallBack,omitempty"`

	// The HTTP method we should use to call status_callback. Can be: POST or GET and the default is POST.
	StatusCallBackMethod *string `form:"StatusCallBackMethod,omitempty"`

	// Whether to clip the intervals where there is no active media in the Compositions
	// triggered by the composition hook. The default is true.
	// Compositions with trim enabled are shorter when the Room is created and no Participant joins for a while
	// as well as if all the Participants leave the room and join later,
	// because those gaps will be removed. See Specifying Video Layouts for more info.
	Trim *bool `form:"Trim,omitempty"`
}

func (p *CreateParams) FormValues() (url.Values, error) {
	var regionBytes []byte
	hasVideolayout := p.VideoLayout != nil
	if hasVideolayout {
		regionMap := make(map[string]interface{})
		for _, r := range p.VideoLayout.GetRegions() {
			if r == nil {
				return nil, errors.New("Error, the region is nil.")
			}
			if r.Prop == nil {
				return nil, errors.New("Error, the region must have properties.")
			}

			propBytes, err := json.Marshal(r.Prop)
			if err != nil {
				return nil, err
			}

			propObj := make(map[string]interface{})
			if err := json.Unmarshal(propBytes, &propObj); err != nil {
				return nil, err
			}

			regionMap[r.Name] = propObj
		}

		var err error
		regionBytes, err = json.Marshal(regionMap)
		if err != nil {
			return nil, err
		}
	}

	url, err := form.EncodeToValues(p)
	if err != nil {
		return nil, err
	}

	if hasVideolayout {
		url.Set("VideoLayout", string(regionBytes))
	}

	return url, nil
}
