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
	Recordings []RecordingInstance `json:"recordings"`
	Meta       struct {
		Page            int     `json:"page"`
		PageSize        int     `json:"page_size"`
		FirstPageUrl    string  `json:"first_page_url"`
		PreviousPageUrl *string `json:"previous_page_url"`
		Url             string  `json:"url"`
		NextPageUrl     *string `json:"next_page_url"`
		Key             string  `json:"key"`
	} `json:"meta"`
}

type Media struct {
	RedirectTo string `json:"redirect_to"`
}

type Container string

const (
	ContainerMka Container = "mka"
	ContainerMkv Container = "mkv"
)

type Codec string

const (
	CodecPCMU Codec = "PCMU"
	CodecOPUS Codec = "OPUS"
	CodecVP8  Codec = "VP8"
	CodecH264 Codec = "H264"
)

type Type string

const (
	TypeAudio Type = "audio"
	TypeVideo Type = "video"
)

type Operation string

const (
	OperationRecordingStart    Operation = "RecordingStart"
	OperationRecordingComplete Operation = "RecordingComplete"
	OperationRecordingUpload   Operation = "RecordingUpload"
)

// https://www.twilio.com/docs/video/api/status-callbacks#recordings-event-parameters
type RecordingCallbackParameters struct {
	// The AccountSid associated with this Room
	AccountSid string `form:"AccountSid"`

	// RecordingSid.
	RecordingSid string `form:"RecordingSid"`

	// Time of the event, conformant to UTC ISO 8601 Timestamp.
	Timestamp time.Time `form:"Timestamp"`

	// The Room event. For example, rooms-created. See Rooms Status Callback Events for the complete list.
	StatusCallbackEvent string `form:"StatusCallbackEvent"`

	// The Sid of the Room generating this event.
	RoomSid string `form:"RoomSid"`

	// The UniqueName of the Room generating this event.
	RoomName string `form:"RoomName"`

	// The Type of the Room generating this event.
	RoomType string `form:"RoomType"`

	// The Sid for the Participant generating this event.
	ParticipantSid string `form:"ParticipantSid"`

	// This recording's source TrackSID, MTxxx.
	SourceSid string `form:"SourceSid"`

	// The relative URL to retrieve this recording's metadata.
	RecordingUri string `form:"RecordingUri"`

	// URL to fetch the generated media.
	MediaUri string `form:"MediaUri"`

	// Duration of the recording. Unsure for data type.
	Duration interface{} `form:"Duration"`

	// Total number of bytes recorded. Unsure for data type.
	Size interface{} `form:"Size"`

	// URL to fetch the generated media if stored in external storage.
	MediaExternalLocation string `form:"MediaExternalLocation"`

	// Container of the recording. Container used are mka for audio recordings and mkv for video recordings.
	Container Container `form:"Container"`

	// Codec used for this recording. This could be PCMU or OPUS for audio recordings, and VP8 or H264
	// for video recordings.
	Codec Codec `form:"Codec"`

	// The Identity of the Participant generating this event. Participant identities are set via
	// the Participant's Access Token
	ParticipantIdentity string `form:"ParticipantIdentity"`

	// The name that was given to the source track of this recording. If no name is given, the SourceSid is used.
	TrackName string `form:"TrackName"`

	// The time in milliseconds elapsed between an arbitrary point in time, common to all Rooms, and the moment
	// when the source room of this track started. This information provides a synchronization mechanism for
	// recordings belonging to the same room.
	OffsetFromTwilioVideoEpoch int64 `form:"OffsetFromTwilioVideoEpoch"`

	// The Status of the Room generating this event.
	RoomStatus string `form:"RoomStatus"`

	// Only on participant-disconnected event.
	// The total duration the Participant remained connected to the Room.
	ParticipantDuration *uint64 `form:"ParticipantDuration"`

	// The type of track for this recording, audio or video.
	Type Type `form:"Type"`

	// Operation that failed: RecordingStart, RecordingComplete, RecordingUpload.
	FailedOperation Operation `form:"FailedOperation"`
}
