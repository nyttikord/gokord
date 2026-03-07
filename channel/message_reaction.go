package channel

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nyttikord/gokord/discord"
	. "github.com/nyttikord/gokord/discord/request"
	"github.com/nyttikord/gokord/emoji"
	"github.com/nyttikord/gokord/user"
)

// Reactions represents a reaction to a [Message].
type Reactions struct {
	Count int          `json:"count"`
	Me    bool         `json:"me"`
	Emoji *emoji.Emoji `json:"emoji"`
}

// AddReaction to a [Message].
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format
// (e.g. "hello:1234567654321").
func AddReaction(channelID, messageID uint64, emojiID string) Empty {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	req := NewSimple(http.MethodPut, discord.EndpointMessageReaction(channelID, messageID, emojiID, 0)).
		WithBucketID(discord.EndpointMessageReaction(channelID, 0, "", 0))
	return WrapAsEmpty(req)
}

// DeleteReaction to a [Message].
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format
// (e.g. "hello:1234567654321").
func DeleteReaction(channelID, messageID uint64, emojiID string, userID uint64) Empty {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	req := NewSimple(http.MethodDelete, discord.EndpointMessageReaction(channelID, messageID, emojiID, userID)).
		WithBucketID(discord.EndpointMessageReaction(channelID, 0, "", 0))
	return WrapAsEmpty(req)
}

// DeleteAllReactions from a [Message].
func DeleteAllReactions(channelID, messageID uint64) Empty {
	req := NewSimple(http.MethodDelete, discord.EndpointMessageReactionsAll(channelID, messageID)).
		WithBucketID(discord.EndpointMessageReactionsAll(channelID, 0))
	return WrapAsEmpty(req)
}

// DeleteEmojiReactions deletes all reactions of a certain [emoji.Emoji] from a [Message].
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format
// (e.g. "hello:1234567654321").
func DeleteEmojiReactions(channelID, messageID uint64, emojiID string) Empty {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	req := NewSimple(http.MethodDelete, discord.EndpointMessageReactions(channelID, messageID, emojiID)).
		WithBucketID(discord.EndpointMessageReactions(channelID, 0, ""))
	return WrapAsEmpty(req)
}

// ListReactions gets all the [Reactions] for a specific [emoji.Emoji].
//
// emojiID is either the Unicode emoji for the reaction, or a guild emoji identifier in name:id format
// (e.g. "hello:1234567654321").
// limit is the max number of users to return (max 100).
// If provided all reactions returned will be before beforeID.
// If provided all reactions returned will be after afterID.
func ListReactions(channelID, messageID uint64, emojiID string, limit int, beforeID, afterID string) Request[[]*user.User] {
	// emoji such as  #⃣ need to have # escaped
	emojiID = strings.ReplaceAll(emojiID, "#", "%23")
	uri := discord.EndpointMessageReactions(channelID, messageID, emojiID)

	v := url.Values{}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	return NewData[[]*user.User](http.MethodGet, uri).
		WithBucketID(discord.EndpointMessageReaction(channelID, 0, "", 0))
}
