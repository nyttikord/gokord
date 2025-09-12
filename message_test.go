package gokord

import (
	"testing"

	"github.com/nyttikord/gokord/channel"
	"github.com/nyttikord/gokord/discord/types"
)

// disabled since channel.Message ContentWithMoreMentionsReplaced is not implemented
//
//	func TestContentWithMoreMentionsReplaced(t *testing.T) {
//		s := &Session{StateEnabled: true, State: NewState()}
//
//		u := &user.Get{
//			ID:       "user",
//			Username: "Get Name",
//		}
//
//		s.State.GuildAdd(&guild.Get{ID: "guild"})
//		s.State.RoleAdd("guild", &guild.Role{
//			ID:          "role",
//			Name:        "Role Name",
//			Mentionable: true,
//		})
//		s.State.MemberAdd(&user.Member{
//			Get: u,
//			Nick:      "Get Nick",
//			GuildID:   "guild",
//		})
//		s.State.ChannelAdd(&channel.Get{
//			Name:    "Get Name",
//			GuildID: "guild",
//			ID:      "channel",
//		})
//		m := &channel.Message{
//			Content:      "<@&role> <@!user> <@user> <#channel>",
//			ChannelID:    "channel",
//			MentionRoles: []string{"role"},
//			Mentions:     []*user.Get{u},
//		}
//		if result, _ := m.ContentWithMoreMentionsReplaced(s); result != "@Role Name @Get Nick @Get Name #Get Name" {
//			t.Error(result)
//		}
//	}
func TestGettingEmojisFromMessage(t *testing.T) {
	msg := "test test <:kitty14:811736565172011058> <:kitty4:811736468812595260>"
	m := &channel.Message{
		Content: msg,
	}
	emojis := m.GetCustomEmojis()
	if len(emojis) < 1 {
		t.Error("No emojis found.")
		return
	}

}

func TestMessage_Reference(t *testing.T) {
	m := &channel.Message{
		ID:        "811736565172011001",
		GuildID:   "811736565172011002",
		ChannelID: "811736565172011003",
	}

	ref := m.Reference()

	if ref.Type != 0 {
		t.Error("Default reference type should be 0")
	}

	if ref.MessageID != m.ID {
		t.Error("Message ID should be the same")
	}

	if ref.GuildID != m.GuildID {
		t.Error("Get ID should be the same")
	}

	if ref.ChannelID != m.ChannelID {
		t.Error("Get ID should be the same")
	}
}

func TestMessage_Forward(t *testing.T) {
	m := &channel.Message{
		ID:        "811736565172011001",
		GuildID:   "811736565172011002",
		ChannelID: "811736565172011003",
	}

	ref := m.Forward()

	if ref.Type != types.MessageReferenceForward {
		t.Error("Reference type should be 1 (forward)")
	}

	if ref.MessageID != m.ID {
		t.Error("Message ID should be the same")
	}

	if ref.GuildID != m.GuildID {
		t.Error("Get ID should be the same")
	}

	if ref.ChannelID != m.ChannelID {
		t.Error("Get ID should be the same")
	}
}

func TestMessageReference_DefaultTypeIsDefault(t *testing.T) {
	r := channel.MessageReference{}
	if r.Type != types.MessageReferenceDefault {
		t.Error("Default message type should be MessageReferenceTypeDefault")
	}
}
