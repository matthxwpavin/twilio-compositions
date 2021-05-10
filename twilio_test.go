package twilio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestListCompletedRooms(t *testing.T) {
	twi, err := New(&Credential{
		AccountSid:   "ACbd053199ec191510366e3ff202715466",
		ApiKeySid:    "SK7269d82865b04ff37b71ca4a90bb8555",
		ApiKeySecret: "A9NWMeXdakDv6qDB4hT3sgs9ZVIaVszd",
	})
	if err != nil {
		t.Errorf("error to initialize: %v", err)
	}

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
