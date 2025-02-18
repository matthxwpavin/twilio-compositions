package participants

import "time"

type ParticipantInstance struct {
	AccountSid  string      `json:"account_sid"`
	RoomSid     string      `json:"room_sid"`
	DateCreated time.Time   `json:"date_created"`
	DateUpdated time.Time   `json:"date_updated"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     interface{} `json:"end_time"`
	Sid         string      `json:"sid"`
	Identity    string      `json:"identity"`
	Status      string      `json:"status"`
	URL         string      `json:"url"`
	Duration    interface{} `json:"duration"`
	Links       struct {
		PublishedTracks  string `json:"published_tracks"`
		SubscribedTracks string `json:"subscribed_tracks"`
		SubscribeRules   string `json:"subscribe_rules"`
		Anonymize        string `json:"anonymize"`
	} `json:"links"`
}
