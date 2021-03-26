package twilio

import (
	"github.com/matthxwpavin/twilio-compositions/video"
	"github.com/matthxwpavin/twilio-compositions/video/composition"
	"github.com/matthxwpavin/twilio-compositions/video/compositionhooks"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func New(credential *Credential) (*Twilio, error) {
	if credential.AccountSid == "" || credential.ApiKeySecret == "" || credential.ApiKeySid == "" {
		return nil, errors.New("Error, each credential must not be empty")
	}
	return &Twilio{
		cred:    credential,
		baseUrl: video.BaseUrl,
		client:  http.DefaultClient,
	}, nil
}

func (t *Twilio) CreateCompositionHooks(param *compositionhooks.CreateParams) (*compositionhooks.CompositionHooks, error) {
	if param.FriendlyName == "" {
		return nil, errors.New("Error, Friendly Name must not be nil.")
	}
	if err := t.validateResolution(param); err != nil {
		return nil, err
	}

	form, err := param.FormValues()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		t.baseUrl.WithCompositionHooksURI(),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	respBody, err := t.fireWithAuth(req)
	if err != nil {
		return nil, err
	}

	ret := &compositionhooks.CompositionHooks{}
	if err := json.Unmarshal(respBody, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *Twilio) ListEnabledCompositionHooks() (*compositionhooks.CompositionHooksList, error) {
	req, err := http.NewRequest(http.MethodGet, t.baseUrl.WithCompositionHooksURI(), nil)
	if err != nil {
		return nil, err
	}

	respBody, err := t.fireWithAuth(req)
	if err != nil {
		return nil, err
	}

	ret := &compositionhooks.CompositionHooksList{}
	if err := json.Unmarshal(respBody, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Twilio) DeleteCompositionHooks(hooksSid string) error {
	req, err := http.NewRequest(http.MethodDelete, t.baseUrl.WithCompositionHooksURIAndPathParam(hooksSid), nil)
	if err != nil {
		return err
	}

	if _, err := t.fireWithAuth(req); err != nil {
		return err
	}

	return nil
}

func (t *Twilio) ListRoomCompletedCompositions(roomSid string) (*composition.CompositionList, error) {
	req, err := http.NewRequest(http.MethodGet, t.baseUrl.WithCompositionURIAndQueryParameters(url.Values{
		"RoomSid": []string{roomSid},
		"Status":  []string{composition.StatusCompleted},
	}), nil)
	if err != nil {
		return nil, err
	}

	respBody, err := t.fireWithAuth(req)
	if err != nil {
		return nil, err
	}

	ret := &composition.CompositionList{}
	if err := json.Unmarshal(respBody, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Twilio) GetCompositionMedia(comSid string) (*composition.Composition, error) {
	req, err := http.NewRequest(http.MethodGet, t.baseUrl.WithCompositionURIMedia(comSid), nil)
	if err != nil {
		return nil, err
	}

	respBody, err := t.fireWithAuth(req)
	if err != nil {
		return nil, err
	}

	ret := &composition.Composition{}
	if err := json.Unmarshal(respBody, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *Twilio) fireWithAuth(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(t.cred.ApiKeySid, t.cred.ApiKeySecret)
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(body))
	}

	return io.ReadAll(resp.Body)
}

func (t *Twilio) validateResolution(param *compositionhooks.CreateParams) error {
	resolution := param.VideoLayout.GetResolution()
	if param.Resolution == nil {
		resolution = string(compositionhooks.VGA)
		if resolution == param.VideoLayout.GetResolution() {
			return nil
		}
	} else {
		if *param.Resolution == resolution {
			return nil
		}
	}

	return errors.New("Error, miss match resolution that given video layout.")
}
