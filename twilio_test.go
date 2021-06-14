package twilio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/matthxwpavin/twilio-compositions/video"
	"github.com/matthxwpavin/twilio-compositions/video/composition"
	"github.com/matthxwpavin/twilio-compositions/video/rooms"
	"github.com/pelletier/go-toml"
)

var twi = func() *Twilio {
	path, err := filepath.Abs("../cn-std-api-server/conf/server_conf_dev.toml")
	if err != nil {
		panic(err)
	}

	tree, err := toml.LoadFile(filepath.Clean(path))
	if err != nil {
		panic(err)
	}

	_map := tree.ToMap()["twilio"].(map[string]interface{})
	t, err := New(&Credential{
		AccountSid:   _map["account_sid"].(string),
		ApiKeySid:    _map["api_key_sid"].(string),
		ApiKeySecret: _map["api_key_secret"].(string),
	})
	if err != nil {
		panic(err)
	}
	return t
}()

func TestListCompletedRooms(t *testing.T) {
	rooms, err := twi.ListCompletedRooms(30)
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
	_, err = twi.CreateComposition(&composition.ComposeParams{
		RoomSid:              "x",
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

func TestDeleteCompositionHooks(t *testing.T) {
	if err := twi.DeleteCompositionHooks("HKf4f715876f9c46d8767542315c7b87ee"); err != nil {
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

func TestListRoomsByUniqueName(t *testing.T) {
	rooms, err := twi.ListRooms(url.Values{
		"UniqueName": []string{"x"},
		"Status":     []string{"completed"},
	})
	if err != nil {
		t.Errorf("error to list rooms: %v", err)
	}

	jsonPrint(rooms)
}

func TestListCompositions(t *testing.T) {
	status := composition.StatusCompleted
	_, err := time.Parse("2006-01-02 15:04:05Z07:00", "2021-05-18 00:00:00+00:00")
	if err != nil {
		t.Errorf("error to parse time: %v", err)
	}

	roomSid := ""
	//after := afterDate.Format(time.RFC3339)
	param := composition.GetParams{
		Status: &status,
		//DateCreatedAfter:  &after,
		//DateCreatedBefore: nil,
		RoomSid: &roomSid,
	}

	ret, err := twi.ListCompositions(&param)
	if err != nil {
		t.Errorf("error to list compositions: %v", err)
	}

	jsonPrint(ret)
}

func TestGetRoomBySid(t *testing.T) {
	room, err := twi.GetRoomInstance("RMfa70a2c1699bf4d5045d2525ea5af946")
	if err != nil {
		t.Errorf("error to get a room: %v", err)
	}

	jsonPrint(room)
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

	fmt.Printf(buf.String())
}
