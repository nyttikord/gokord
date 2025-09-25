package gokord

import (
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/guild"
)

// setGuildIds will set the GuildID on all the members of a guild.Guild.
// This is done as event data does not have it set.
func setGuildIds(g *guild.Guild) {
	for _, c := range g.Channels {
		c.GuildID = g.ID
	}

	for _, m := range g.Members {
		m.GuildID = g.ID
	}

	for _, vs := range g.VoiceStates {
		vs.GuildID = g.ID
	}
}

// onInterface handles all internal events and routes them to the appropriate internal handler.
func (s *Session) onInterface(i any) {
	switch t := i.(type) {
	case *event.Ready:
		for _, g := range t.Guilds {
			setGuildIds(g)
		}
		s.onReady(t)
	case *event.GuildCreate:
		setGuildIds(t.Guild)
	case *event.GuildUpdate:
		setGuildIds(t.Guild)
	case *event.VoiceServerUpdate:
		go s.onVoiceServerUpdate(t)
	case *event.VoiceStateUpdate:
		go s.onVoiceStateUpdate(t)
	}
	err := s.sessionState.onInterface(s, i)
	if err != nil {
		s.LogDebug("error dispatching internal event, %s", err)
	}
}

// onReady handles the ready event.
func (s *Session) onReady(r *event.Ready) {

	// Store the SessionID within the Session struct.
	s.sessionID = r.SessionID

	// Store the ResumeGatewayURL within the Session struct.
	s.resumeGatewayURL = r.ResumeGatewayURL
}
