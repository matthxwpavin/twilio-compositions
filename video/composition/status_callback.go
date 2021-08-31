package composition

import "time"

// https://www.twilio.com/docs/video/api/status-callbacks#compositions-events
const (
	StatusCallbackEnqueued   = "composition-enqueued"
	StatusCallbackHookFailed = "composition-hook-failed"
	StatusCallbackStarted    = "composition-started"
	StatusCallbackAvailable  = "composition-available"
	StatusCallbackProgress   = "composition-progress"
	StatusCallbackFailed     = "composition-failed"
)

// https://www.twilio.com/docs/video/api/status-callbacks#compositions-event-parameters
type CallbackParam struct {
	AccountSid          string    `form:"AccountSid"`
	RoomSid             string    `form:"RoomSid"`
	HookSid             string    `form:"HookSid"`
	HookUri             string    `form:"HookUri"`
	HookFriendlyName    string    `form:"HookFriendlyName"`
	CompositionSid      string    `form:"CompositionSid"`
	CompositionUri      string    `form:"CompositionUri"`
	MediaUri            string    `form:"MediaUri"`
	Duration            uint      `form:"Duration"`
	Size                uint      `form:"Size"`
	PercentageDone      float64   `form:"PercentageDone"`
	SecondsRemaining    float64   `form:"SecondsRemaining"`
	FailedOperation     string    `form:"FailedOperation"`
	ErrorMessage        string    `form:"ErrorMessage"`
	StatusCallbackEvent string    `form:"StatusCallbackEvent"`
	Timestamp           time.Time `form:"Timestamp"`
}
