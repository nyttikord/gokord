// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gokord

import (
	"context"
	"testing"
	"time"

	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/user"
)

func TestSession(t *testing.T) {
	if envBotToken == "" {
		t.Skip("Skipping session test, DG_TOKEN not set")
	}
	dgBot.EventManager().AddHandler(func(ctx context.Context, s bot.Session, r *event.Ready) {
		bot.Logger(ctx).Info("bot ready")
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	err := dgBot.OpenAndBlock(ctx)
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func TestMember_DisplayName(t *testing.T) {
	u := &user.User{
		GlobalName: "Global",
	}
	t.Run("no server nickname set", func(t *testing.T) {
		m := &user.Member{
			Nick: "",
			User: u,
		}
		want := u.DisplayName()
		if dn := m.DisplayName(); dn != want {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, want)
		}
	})
	t.Run("server nickname set", func(t *testing.T) {
		m := &user.Member{
			Nick: "Server",
			User: u,
		}
		if dn := m.DisplayName(); dn != m.Nick {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, m.Nick)
		}
	})
}
