// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gokord

import (
	"testing"

	"github.com/nyttikord/gokord/user"
)

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
