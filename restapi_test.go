package gokord

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// TestChannelMessageSend tests the ChannelMessageSend() function. This should not return an error.
func TestChannelMessageSend(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	ctx := context.Background()

	_, err := dg.ChannelAPI().MessageSend(ctx, envChannel, "Running REST API Tests!")
	if err != nil {
		t.Errorf("ChannelMessageSend returned error: %+v", err)
	}
}

/*
// removed for now, only works on BOT accounts now
func TestUserAvatar(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot TestUserAvatar, dg not set.")
	}

	u, err := dg.Application("@me")
	if err != nil {
		t.Error("error fetching @me user,", err)
	}

	a, err := dg.UserAvatar(u.ID)
	if err != nil {
		if err.Error() == `HTTP 404 NOT FOUND, {"code": 0, "message": "404: Not Found"}` {
			t.Skip("Skipped, @me doesn't have an Avatar")
		}
		t.Errorf(err.Error())
	}

	if a == nil {
		t.Errorf("a == nil, should be image.Image")
	}
}
*/

/* Running this causes an error due to 2/hour rate limit on username changes
func TestUserUpdate(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot test logout, dg not set.")
	}

	u, err := dg.Application("@me")
	if err != nil {
		t.Errorf(err.Error())
	}

	s, err := dg.UserUpdate(envEmail, envPassword, "testname", u.Avatar, "")
	if err != nil {
		t.Error(err.Error())
	}
	if s.Username != "testname" {
		t.Error("Username != testname")
	}
	s, err = dg.UserUpdate(envEmail, envPassword, u.Username, u.Avatar, "")
	if err != nil {
		t.Error(err.Error())
	}
	if s.Username != u.Username {
		t.Error("Username != " + u.Username)
	}
}
*/

//func (s *Session) UserChannelCreate(recipientID string) (st *Application, err error) {

func TestUserChannelCreate(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot TestUserChannelCreate, dg not set.")
	}

	if envAdmin == "" {
		t.Skip("Skipped, DG_ADMIN not set.")
	}

	ctx := context.Background()

	_, err := dg.UserAPI().ChannelCreate(ctx, envAdmin)
	if err != nil {
		t.Error(err.Error())
	}

	// TODO make sure the channel was added
}

func TestUserGuilds(t *testing.T) {
	if dg == nil {
		t.Skip("Cannot TestUserGuilds, dg not set.")
	}

	ctx := context.Background()

	_, err := dg.GuildAPI().UserGuilds(ctx, 10, "", "", false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGateway(t *testing.T) {
	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}
	_, err := dg.Gateway(context.Background())
	if err != nil {
		t.Errorf("Gateway() returned error: %+v", err)
	}
}

func TestGatewayBot(t *testing.T) {
	if dgBot == nil {
		t.Skip("Skipping, dgBot not set.")
	}
	_, err := dgBot.GatewayBot(context.Background())
	if err != nil {
		t.Errorf("GatewayBot() returned error: %+v", err)
	}
}

func TestVoiceRegions(t *testing.T) {
	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.GuildAPI().VoiceRegions(context.Background())
	if err != nil {
		t.Errorf("VoiceRegions() returned error: %+v", err)
	}
}
func TestGuildRoles(t *testing.T) {
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.GuildAPI().Roles(context.Background(), envGuild)
	if err != nil {
		t.Errorf("GuildRoles(envGuild) returned error: %+v", err)
	}

}

func TestGuildMemberNickname(t *testing.T) {
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	err := dg.GuildAPI().MemberNickname(context.Background(), envGuild, "@me/nick", "B1nzyRocks")
	if err != nil {
		t.Errorf("GuildNickname returned error: %+v", err)
	}
}

// TestChannelMessageSend2 tests the ChannelMessageSend() function. This should not return an error.
func TestChannelMessageSend2(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.ChannelAPI().MessageSend(context.Background(), envChannel, "All done running REST API Tests!")
	if err != nil {
		t.Errorf("ChannelMessageSend returned error: %+v", err)
	}
}

// TestGuildPruneCount tests PruneCount() function. This should not return an error.
func TestGuildPruneCount(t *testing.T) {
	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.GuildAPI().PruneCount(context.Background(), envGuild, 1)
	if err != nil {
		t.Errorf("PruneCount returned error: %+v", err)
	}
}

/*
// TestGuildPrune tests Prune() function. This should not return an error.
func TestGuildPrune(t *testing.T) {

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	_, err := dg.Prune(envGuild, 1)
	if err != nil {
		t.Errorf("Prune returned error: %+v", err)
	}
}
*/

func Test_unmarshal(t *testing.T) {
	err := unmarshal([]byte{}, &struct{}{})
	if !errors.Is(err, ErrJSONUnmarshal) {
		t.Errorf("Unexpected error type: %T", err)
	}
}

// roundTripperFunc implements http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
