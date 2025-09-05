// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gokord

import (
	user2 "github.com/nyttikord/gokord/user"
	"testing"
)

func TestMember_DisplayName(t *testing.T) {
	user := &user2.User{
		GlobalName: "Global",
	}
	t.Run("no server nickname set", func(t *testing.T) {
		m := &Member{
			Nick:       "",
			user2.User: user,
		}
		want := user.DisplayName()
		if dn := m.DisplayName(); dn != want {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, want)
		}
	})
	t.Run("server nickname set", func(t *testing.T) {
		m := &Member{
			Nick:       "Server",
			user2.User: user,
		}
		if dn := m.DisplayName(); dn != m.Nick {
			t.Errorf("Member.DisplayName() = %v, want %v", dn, m.Nick)
		}
	})
}
