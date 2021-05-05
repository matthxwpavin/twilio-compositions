package video

import "net/url"

type VideoUrl string

const (
	BaseUrl = VideoUrl("https://video.twilio.com")
)

func (url VideoUrl) WithCompositionHooksURI() string {
	return url.WithCompositionHooksURIAndPathParam("")
}

func (url VideoUrl) WithCompositionHooksURIAndPathParam(compositionHooksSid string) string {
	pathParam := ""
	if compositionHooksSid != "" {
		pathParam = "/" + compositionHooksSid
	}
	return string(url) + "/v1/CompositionHooks" + pathParam
}

func (url VideoUrl) WithCompositionURI() string {
	return string(url) + "/v1/Compositions"
}

func (url VideoUrl) WithCompositionURIMedia(compositionSid string) string {
	return url.WithCompositionURI() + "/" + compositionSid + "/Media"
}

func (url VideoUrl) WithCompositionURIAndQueryParameters(values url.Values) string {
	return url.WithCompositionURI() + "?" + values.Encode()
}

func (url VideoUrl) WithRoomsURI() string {
	return string(url) + "/v1/Rooms"
}

func (url VideoUrl) WithRoomsURIAndQueryParameters(values url.Values) string {
	return url.WithRoomsURI() + "?" + values.Encode()
}
