package event

import (
	"encoding/json"

	"github.com/nyttikord/gokord/internal/structs"
)

func (t *ThreadListSync) UnmarshalJSON(data []byte) error {
	type tt ThreadListSync
	var v struct {
		tt
		ChannelIDs []string `json:"channel_ids"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*t = ThreadListSync(v.tt)
	t.ChannelIDs = structs.SnowflakesToUints(v.ChannelIDs)
	return nil
}

func (t *ThreadMembersUpdate) UnmarshalJSON(data []byte) error {
	type tt ThreadMembersUpdate
	var v struct {
		tt
		RemovedMembers []string `json:"removed_member_ids"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*t = ThreadMembersUpdate(v.tt)
	t.RemovedMembers = structs.SnowflakesToUints(v.RemovedMembers)
	return nil
}

func (m *MessageDeleteBulk) UnmarshalJSON(data []byte) error {
	type t MessageDeleteBulk
	var v struct {
		t
		Messages []string `json:"ids"`
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*m = MessageDeleteBulk(v.t)
	m.Messages = structs.SnowflakesToUints(v.Messages)
	return nil
}
