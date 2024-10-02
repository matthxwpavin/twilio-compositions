package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ajg/form"
	"github.com/matthxwpavin/twilio-compositions/video"
	"github.com/matthxwpavin/twilio-compositions/video/composition"
	"github.com/matthxwpavin/twilio-compositions/video/recording"
	"github.com/matthxwpavin/twilio-compositions/video/rooms"
	"github.com/spf13/viper"
)

type Credential struct {
	AccountSid,
	ApiKeySid,
	ApiKeySecret string
}

func LoadCredentials() *Credential {
	cdtMap := viper.GetStringMapString("twilio")
	return &Credential{
		AccountSid:   cdtMap["account_sid"],
		ApiKeySid:    cdtMap["api_key_sid"],
		ApiKeySecret: cdtMap["api_key_secret"],
	}
}

type Twilio struct {
	cred    *Credential
	baseUrl video.VideoUrl
	client  *http.Client
}

func New(credential *Credential) *Twilio {
	return NewWithHttpClient(credential, http.DefaultClient)
}

func NewWithHttpClient(credential *Credential, httpClient *http.Client) *Twilio {
	if credential == nil {
		panic("credential must not be nil")
	}
	if httpClient == nil {
		panic("http client must not be nil")
	}

	return &Twilio{
		cred:    credential,
		baseUrl: video.BaseUrl,
		client:  httpClient,
	}
}

func (t *Twilio) CreateComposition(
	param *composition.ComposeParams,
) (*composition.Composition, error) {
	if err := t.validateResolution(param); err != nil {
		return nil, err
	}

	form, err := t.formValues(param)
	if err != nil {
		return nil, err
	}

	ret := &composition.Composition{}
	if err := t.request(
		http.MethodPost,
		t.baseUrl.WithCompositionURI(),
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
		nil,
		ret,
	); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *Twilio) CreateCompositionHooks(
	param *composition.HooksParams,
) (*composition.CompositionHooks, error) {
	return t.requestCompositionHooks(http.MethodPost, t.baseUrl.WithCompositionHooksURI(), param)
}

func (t *Twilio) UpdateCompositionHooks(
	hooksSid string,
	param *composition.HooksParams,
) (*composition.CompositionHooks, error) {
	if hooksSid == "" {
		return nil, errors.New("Hooks SID must not be empty")
	}
	return t.requestCompositionHooks(
		http.MethodPost,
		t.baseUrl.WithCompositionHooksURIAndPathParam(hooksSid),
		param,
	)
}

func (t *Twilio) requestCompositionHooks(
	method, url string,
	param *composition.HooksParams,
) (*composition.CompositionHooks, error) {
	if param.FriendlyName == "" {
		return nil, errors.New("Error, Friendly Name must not be nil.")
	}
	if err := t.validateResolution(param); err != nil {
		return nil, err
	}

	form, err := t.formValues(param)
	if err != nil {
		return nil, err
	}

	ret := &composition.CompositionHooks{}
	if err := t.request(
		method,
		url,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
		nil,
		ret,
	); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *Twilio) ListEnabledCompositionHooks() (*composition.CompositionHooksList, error) {
	ret := &composition.CompositionHooksList{}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithCompositionHooksURI(),
		"",
		nil,
		nil,
		ret,
	); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Twilio) DeleteCompositionHooks(hooksSid string) error {
	return t.request(
		http.MethodDelete,
		t.baseUrl.WithCompositionHooksURIAndPathParam(hooksSid),
		"",
		nil,
		nil,
		nil,
	)
}

func (t *Twilio) ListCompositions(
	param *composition.GetParams,
) (*composition.CompositionList, error) {
	ret := &composition.CompositionList{}
	values, err := form.EncodeToValues(param)
	if err != nil {
		return nil, err
	}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithCompositionURIAndQueryParameters(values),
		"",
		nil,
		nil,
		ret,
	); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Twilio) ListRoomCompletedCompositions(
	roomSid string,
) (*composition.CompositionList, error) {
	ret := &composition.CompositionList{}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithCompositionURIAndQueryParameters(url.Values{
			"RoomSid": []string{roomSid},
			"Status":  []string{string(composition.StatusCompleted)},
		}),
		"",
		nil,
		nil,
		ret,
	); err != nil {
		return nil, err
	}

	return ret, nil
}

type RecordingFilter struct {
	MediaType string
}

func (t *Twilio) ListRecordingsByRoomSid(
	roomSid string,
	filter *RecordingFilter,
) (*recording.RecordingList, error) {
	params := url.Values{"GroupingSid": []string{roomSid}}
	if filter != nil {
		if filter.MediaType != "" {
			params.Set("MediaType", filter.MediaType)
		}
	}

	dst := &recording.RecordingList{}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithRecordingsURIAndQueryParam(params),
		"",
		nil,
		nil,
		dst,
	); err != nil {
		return nil, err
	}
	return dst, nil
}

func (t *Twilio) GetRecordingMedia(recordingSid string) (*recording.Media, error) {
	dst := &recording.Media{}
	return dst, t.request(
		http.MethodGet,
		t.baseUrl.WithRecordingsURI()+fmt.Sprintf("/%s/Media", recordingSid),
		"",
		nil,
		func(status int) bool { return status == http.StatusFound },
		dst,
	)
}

func (t *Twilio) GetCompositionMedia(comSid string) (*composition.Composition, error) {
	ret := &composition.Composition{}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithCompositionURIMedia(comSid),
		"",
		nil,
		nil,
		ret,
	); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *Twilio) fireWithAuth(req *http.Request, checkStatus func(int) bool) ([]byte, error) {
	req.SetBasicAuth(t.cred.ApiKeySid, t.cred.ApiKeySecret)
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	if checkStatus == nil {
		checkStatus = func(i int) bool {
			return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
		}
	}

	if !checkStatus(resp.StatusCode) {
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"unexpected status code, status: %v, message: %s",
			resp.StatusCode,
			msg,
		)
	}
	return io.ReadAll(resp.Body)
}

func (t *Twilio) validateResolution(param video.VideoLayouter) error {
	resolution := param.GetVideoLayout().GetResolution()
	if param.GetResolution() == nil {
		resolution = composition.VGA
		if resolution == param.GetVideoLayout().GetResolution() {
			return nil
		}
	} else {
		if *param.GetResolution() == resolution {
			return nil
		}
	}

	return errors.New("Error, miss match resolution that given video layout.")
}

func (t *Twilio) GetRoomInstance(roomSid string) (*rooms.RoomInstance, error) {
	dst := &rooms.RoomInstance{}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithRoomsURI()+"/"+roomSid,
		"",
		nil,
		nil,
		dst,
	); err != nil {
		return nil, err
	}
	return dst, nil
}

func (t *Twilio) ListCompletedRooms(size uint) (*rooms.RoomInstanceList, error) {
	return t.ListRooms(url.Values{
		"Status":   []string{"completed"},
		"PageSize": []string{fmt.Sprintf("%d", size)},
	})
}

func (t *Twilio) ListRooms(params url.Values) (*rooms.RoomInstanceList, error) {
	dst := &rooms.RoomInstanceList{}
	if err := t.request(
		http.MethodGet,
		t.baseUrl.WithRoomsURIAndQueryParameters(params),
		"",
		nil,
		nil,
		dst,
	); err != nil {
		return nil, err
	}
	return dst, nil
}

func (t *Twilio) request(
	method, url, contentType string,
	body io.Reader,
	checkStatus func(int) bool,
	dst interface{},
) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	respBody, err := t.fireWithAuth(req, checkStatus)
	if err != nil {
		return err
	}

	if dst == nil {
		return nil
	}
	return json.Unmarshal(respBody, dst)
}

func (t *Twilio) formValues(p video.VideoLayouter) (url.Values, error) {

	var regionBytes []byte
	layout := p.GetVideoLayout()
	hasVideolayout := layout != nil
	if hasVideolayout {
		regionMap := make(map[string]interface{})
		for _, r := range layout.GetRegions() {
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

func (t *Twilio) CreateRoom(param *rooms.RoomPostParams) (*rooms.RoomInstance, error) {
	body, err := form.EncodeToValues(param)
	if err != nil {
		return nil, err
	}

	resp := &rooms.RoomInstance{}
	if err := t.request(
		http.MethodPost,
		t.baseUrl.WithRoomsURI(),
		"application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()),
		nil,
		resp,
	); err != nil {
		return nil, err
	}
	return resp, nil
}
