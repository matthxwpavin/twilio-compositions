package twilio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/matthxwpavin/twilio-compositions/video"
	"github.com/matthxwpavin/twilio-compositions/video/composition"
	"github.com/matthxwpavin/twilio-compositions/video/rooms"
	"net/http"
	"testing"
)

var twi = func() *Twilio {
	t, err := New(&Credential{
		AccountSid:   "ACbd053199ec191510366e3ff202715466",
		ApiKeySid:    "SK7269d82865b04ff37b71ca4a90bb8555",
		ApiKeySecret: "A9NWMeXdakDv6qDB4hT3sgs9ZVIaVszd",
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

func TestConfigRoomStatusCallback(t *testing.T) {
	if err := twi.ConfigRoomsStatusCallback(&rooms.ConfigCallbackParams{
		UniqueName:           "CNRoomsStatusCallback",
		StatusCallback:       "https://dev.clicknic.co/api/v1/videoCall//statusCallback/rooms",
		StatusCallbackMethod: http.MethodPost,
		Type:                 rooms.TypeGroup,
	}); err != nil {
		t.Errorf("error to config rooms status callback: %v", err)
	}
}

func TestListEnableCompositionHooks(t *testing.T) {
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
	trim := true
	AudSource := "*"
	res := composition.VGA

	_, err = twi.CreateComposition(&composition.ComposeParams{
		RoomSid:              "",
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
