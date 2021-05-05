package room

import "time"

const (
	// https://www.twilio.com/docs/video/api/status-callbacks#rooms-callback-events

	// Room created.
	StatusCallbackCreated = "room-created"

	// Room completed.
	// (Note: Rooms created by the REST API will
	// fire room-ended event when room is empty for 5 minutes.)
	StatusCallbackEnded = "room-ended"

	// Participant joined the Room.
	StatusPartCon = "participant-connected"

	// Participant left the Room.
	StatusPartDisCon = "participant-disconnected"

	// Participant added a Track.
	StatusTrackAdded = "track-added"

	// Participant removed a Track.
	StatusTrackRemoved = "track-removed"

	// Participant unpaused a Track.
	StatusTrackEnabled = "track-enabled"

	// Participant paused a Track.
	StatusTrackDisabled = "track-disabled"

	// Recording for a Track began
	StatusRecStarted = "recording-started"

	// Recording for a Track completed
	StatusRecCompleted = "recording-completed"

	// Failure during a recording operation request
	StatusRecFailed = "recording-failed"
)

type RoomInstanceList struct {
	Rooms []RoomInstance `json:"rooms"`
	Meta  struct {
		Page            int    `json:"page"`
		PageSize        int    `json:"page_size"`
		FirstPageURL    string `json:"first_page_url"`
		PreviousPageURL string `json:"previous_page_url"`
		URL             string `json:"url"`
		NextPageURL     string `json:"next_page_url"`
		Key             string `json:"key"`
	} `json:"meta"`
}

type RoomInstance struct {
	// The SID of the Account that created the Room resource.
	AccountSid string `json:"account_sid"`

	// The date and time in GMT when the resource was created specified in ISO 8601 format.
	DateCreated time.Time `json:"date_created"`

	// The date and time in GMT when the resource was last updated specified in ISO 8601 format.
	DateUpdated time.Time `json:"date_updated"`

	// The status of the room. Can be: in-progress, failed, or completed.
	Status string `json:"status"`

	// The type of room. Can be: go, peer-to-peer, group-small, or group. The default value is group.
	Type string `json:"type"`

	// The unique string that we created to identify the Room resource.
	Sid string `json:"sid"`

	// Deprecated, now always considered to be true.
	EnableTurn bool `json:"enable_turn"`

	// An application-defined string that uniquely identifies the resource. It can be used
	// as a room_sid in place of the resource's sid in the URL to address the resource.
	// This value is unique for in-progress rooms. SDK clients can use this name to connect to the room.
	// REST API clients can use this name in place of the Room SID
	// to interact with the room as long as the room is in-progress.
	UniqueName string `json:"unique_name"`

	// The maximum number of concurrent Participants allowed in the room.
	MaxParticipants int `json:"max_participants"`

	// The maximum number of published audio, video, and data tracks
	// all participants combined are allowed to publish in the room at the same time.
	MaxConcurrentPublishedTracks int `json:"max_concurrent_published_tracks"`

	// The duration of the room in seconds.
	Duration int `json:"duration"`

	// The HTTP method we use to call status_callback. Can be POST or GET and defaults to POST.
	StatusCallbackMethod string `json:"status_callback_method"`

	// The URL we call using the status_callback_method to send status information
	// to your application on every room event. See Status Callbacks for more info.
	StatusCallback string `json:"status_callback"`

	// Whether to start recording when Participants connect. This feature is not available in peer-to-peer rooms.
	RecordParticipantsOnConnect bool `json:"record_participants_on_connect"`

	// An array of the video codecs that are supported when publishing a track in the room.
	// Can be: VP8 and H264. This feature is not available in peer-to-peer rooms
	VideoCodecs []string `json:"video_codecs"`

	// The region for the media server in Group Rooms.
	// Can be: one of the available Media Regions. This feature is not available in peer-to-peer rooms.
	MediaRegion string `json:"media_region"`

	// The UTC end time of the room in ISO 8601 format.
	EndTime time.Time `json:"end_time"`

	// The absolute URL of the resource.
	URL string `json:"url"`

	// The URLs of related resources.
	Links struct {
		Participants   string `json:"participants"`
		Recordings     string `json:"recordings"`
		RecordingRules string `json:"recording_rules"`
	} `json:"links"`
}

// https://www.twilio.com/docs/video/api/status-callbacks#rooms-event-parameters
type RoomCallBack struct {
	// The AccountSid associated with this Room
	AccountSid string `form:"AccountSid"`

	// The UniqueName of the Room generating this event.
	RoomName string `form:"RoomName"`

	// The Sid of the Room generating this event.
	RoomSid string `form:"RoomSid"`

	// The Status of the Room generating this event.
	RoomStatus string `form:"RoomStatus"`

	// The Type of the Room generating this event.
	RoomType string `form:"RoomType"`

	// The Room event. For example, room-created. See Rooms Status Callback Events for the complete list.
	StatusCallbackEvent string `form:"StatusCallbackEvent"`

	// Time of the event, conformant to UTC ISO 8601 Timestamp.
	Timestamp time.Time `form:"Timestamp"`

	// The Sid for the Participant generating this event.
	ParticipantSid string `form:"ParticipantSid"`

	// Only on participant-disconnected event.
	// The total duration the Participant remained connected to the Room.
	ParticipantDuration *uint64 `form:"ParticipantDuration"`

	// The Identity of the Participant generating this event. Participant identities are set via the Participant's Access Token
	ParticipantIdentity string `form:"ParticipantIdentity"`

	// Only on room-ended
	// The total duration of the Room, in seconds.
	RoomDuration *uint64 `form:"RoomDuration"`

	// An incrementing integer representing the order of the events in the Room.
	SequenceNumber uint64 `form:"SequenceNumber"`

	// The Sid of the Track.
	TrackSid string `form:"TrackSid"`

	// The Kind of the Track (data, audio or video).
	TrackKind string `form:"TrackKind"`
}
