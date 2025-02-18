package twilio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/matthxwpavin/twilio-compositions/video"
	"github.com/matthxwpavin/twilio-compositions/video/composition"
	"github.com/matthxwpavin/twilio-compositions/video/rooms"
	"github.com/pelletier/go-toml"
)

var twi = func() *Twilio {
	tree, err := toml.LoadFile("./credentials.toml")
	if err != nil {
		panic(err)
	}

	_map := tree.ToMap()["twilio"].(map[string]interface{})
	t := NewWithHttpClient(&Credential{
		AccountSid:   _map["account_sid"].(string),
		ApiKeySid:    _map["api_key_sid"].(string),
		ApiKeySecret: _map["api_key_secret"].(string),
	}, &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	})
	if err != nil {
		panic(err)
	}
	return t
}()

func TestListCompletedRooms(t *testing.T) {
	rooms, err := twi.ListCompletedRooms(1)
	if err != nil {
		t.Errorf("error to list completed rooms: %v", err)
	}

	bb, err := json.Marshal(rooms)
	if err != nil {
		t.Errorf("error to marshal: %v", err)
	}

	dst := &bytes.Buffer{}
	if err := json.Indent(dst, bb, "", "\t"); err != nil {
		t.Errorf("error to indent: %v", err)
	}
	fmt.Println(dst.String())
}

func TestListEnabledCompositionHooks(t *testing.T) {
	hooks, err := twi.ListEnabledCompositionHooks()
	if err != nil {
		t.Errorf("error to list composition hooks: %v", err)
	}

	bb, err := json.Marshal(hooks)
	if err != nil {
		t.Errorf("error to marshal composition hooks: %v", err)
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, bb, "", "\t"); err != nil {
		t.Errorf("error to indent: %v", err)
	}
	fmt.Println(dst.String())
}

func TestCreateComposition(t *testing.T) {
	v, err := video.NewVideoLayout(composition.VGA)
	if err != nil {
		t.Errorf("error to new video composition: %v", err)
	}

	reuse := video.ReuseShowOldest
	reg := &video.Region{
		Name: "grid",
		Prop: &video.RegionProp{
			Reuse:                &reuse,
			VideoSources:         []string{"*"},
			VideoSourcesExcluded: nil,
		},
	}

	if err := v.AddRegion(reg); err != nil {
		t.Errorf("error to add region: %v", err)
	}

	var (
		trim      = true
		AudSource = "*"
		res       = composition.VGA
	)
	comp, err := twi.CreateComposition(&composition.ComposeParams{
		RoomSid:              "RM25d7091d712e6f2ef1a589be78976596",
		VideoLayout:          v,
		AudioSources:         &AudSource,
		AudioSourcesExcluded: nil,
		Resolution:           &res,
		Format:               composition.MP4,
		Trim:                 &trim,
	})
	if err != nil {
		t.Errorf("error to create composition: %v", err)
	}
	jsonPrint(comp)
}

func TestDeleteCompositionHooks(t *testing.T) {
	if err := twi.DeleteCompositionHooks("HK9ef12a9c3d22c3c3b05b5f1420125dfc"); err != nil {
		t.Errorf("error to delete composition hooks: %v", err)
	}
}

func TestCreateCompositionHooks(t *testing.T) {
	v, err := video.NewVideoLayout(composition.VGA)
	if err != nil {
		t.Errorf("error to new video composition: %v", err)
	}

	reuse := video.ReuseShowOldest
	reg := &video.Region{
		Name: "grid",
		Prop: &video.RegionProp{
			Reuse:                &reuse,
			VideoSources:         []string{"*"},
			VideoSourcesExcluded: nil,
		},
	}

	if err := v.AddRegion(reg); err != nil {
		t.Errorf("error to add region: %v", err)
	}

	var (
		trim      = true
		AudSource = "*"
		res       = composition.VGA
		enabled   = true
	)
	_, err = twi.CreateCompositionHooks(&composition.HooksParams{
		FriendlyName:         "ClicknicCompositionHooks",
		Enabled:              &enabled,
		VideoLayout:          v,
		AudioSources:         &AudSource,
		AudioSourcesExcluded: nil,
		Resolution:           &res,
		Format:               composition.MP4,
		Trim:                 &trim,
	})
	if err != nil {
		t.Errorf("error to create composition: %v", err)
	}
}

func TestUpdateCompositionHooks(t *testing.T) {
	v, err := video.NewVideoLayout(composition.VGA)
	if err != nil {
		t.Errorf("error to new video composition: %v", err)
	}

	reuse := video.ReuseShowOldest
	reg := &video.Region{
		Name: "grid",
		Prop: &video.RegionProp{
			Reuse:                &reuse,
			VideoSources:         []string{"*"},
			VideoSourcesExcluded: nil,
		},
	}

	if err := v.AddRegion(reg); err != nil {
		t.Errorf("error to add region: %v", err)
	}

	var (
		trim      = true
		AudSource = "*"
		res       = composition.VGA
		enabled   = false
	)
	_, err = twi.UpdateCompositionHooks(
		"HK9ef12a9c3d22c3c3b05b5f1420125dfc",
		&composition.HooksParams{
			FriendlyName:         "ClicknicCompositionHooks",
			Enabled:              &enabled,
			VideoLayout:          v,
			AudioSources:         &AudSource,
			AudioSourcesExcluded: nil,
			Resolution:           &res,
			Format:               composition.MP4,
			Trim:                 &trim,
		})
	if err != nil {
		t.Errorf("error to create composition: %v", err)
	}
}

func TestCreateRoom(t *testing.T) {
	_type := rooms.RoomType("group-small")
	uniqueName := "TestRoom2"
	callbackUrl := "https://xxxxxxx/api/v1/rooms/statusCallback"
	callbackMethod := "POST"
	param := &rooms.RoomPostParams{
		Type:                 &_type,
		UniqueName:           &uniqueName,
		StatusCallback:       &callbackUrl,
		StatusCallbackMethod: &callbackMethod,
	}

	room, err := twi.CreateRoom(param)
	if err != nil {
		t.Errorf("error to create room: %v", err)
	}

	fmt.Println(room)
}

func TestListCompositionsByRoomSid(t *testing.T) {
	status := composition.StatusCompleted
	roomSid := "RMac8c929571dfa4c7262d48dca8c5c355"
	param := composition.GetParams{
		Status:  &status,
		RoomSid: &roomSid,
	}

	ret, err := twi.ListCompositions(&param)
	if err != nil {
		t.Errorf("error to list compositions: %v", err)
	}

	jsonPrint(ret)
}

func TestListCompositions(t *testing.T) {
	status := composition.StatusCompleted
	_, err := time.Parse("2006-01-02 15:04:05Z07:00", "2021-05-18 00:00:00+00:00")
	if err != nil {
		t.Errorf("error to parse time: %v", err)
	}

	// after := afterDate.Format(time.RFC3339)
	param := composition.GetParams{
		Status: &status,
		// DateCreatedAfter:  &after,
		DateCreatedBefore: nil,
		RoomSid:           nil,
	}

	ret, err := twi.ListCompositions(&param)
	if err != nil {
		t.Errorf("error to list compositions: %v", err)
	}

	jsonPrint(ret)
}

func TestGetRoomBySid(t *testing.T) {
	room, err := twi.GetRoomInstance("RM25d7091d712e6f2ef1a589be78976596")
	if err != nil {
		t.Errorf("error to get a room: %v", err)
	}

	jsonPrint(room)
}

func TestListRecordings(t *testing.T) {
	recs, err := twi.ListRecordings(
		RecordingFilter{
			MediaType: MediaTypeAudio,
			RoomSid:   "RM06bc1c5ac394effdb919741e792776b6",
		},
	)
	if err != nil {
		t.Error("could not get recordings", err)
	}
	jsonPrint(recs)
}

func TestGetRecordingMedia(t *testing.T) {
	media, err := twi.GetRecordingMedia("RT99545ec1d5c10b9bed40195372544a9d")
	if err != nil {
		t.Error("could not get recording media", err)
	}
	print(media.RedirectTo)
}

func TestAuthenticateMediaLink(t *testing.T) {
	url, err := twi.AuthenticateMediaLink(
		context.Background(),
		"https://video.twilio.com/v1/Recordings/RTc77b96e6589c7b0dc1b8689f153fb569/Media",
		nil,
	)
	if err != nil {
		t.Errorf("Could not authenticate media link: %v", err)
	}
	print(url)
}

func jsonPrint(scheme interface{}) {
	b, err := json.Marshal(scheme)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	if err := json.Indent(buf, b, "", "\t"); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}

func TestGetRoomParticipants(t *testing.T) {
	resp, err := twi.GetParticipantsByRoomSid("RMfcb5d69c724ce93fdf8f14f80134e853")
	if err != nil {
		t.Fatalf("Could not get room participants: %v", err)
	}
	jsonPrint(resp)
}
