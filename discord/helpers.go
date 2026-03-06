package discord

import (
	"strings"
	"time"
)

// AvatarURL returns the URL to get the avatar from a hash.
func AvatarURL(avatarHash, defaultAvatarURL, staticAvatarURL, animatedAvatarURL, size string) string {
	var URL string
	if avatarHash == "" {
		URL = defaultAvatarURL
	} else if strings.HasPrefix(avatarHash, "a_") {
		URL = animatedAvatarURL
	} else {
		URL = staticAvatarURL
	}

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

// BannerURL returns the URL to get the avatar from a hash.
func BannerURL(bannerHash, staticBannerURL, animatedBannerURL, size string) string {
	var URL string
	if bannerHash == "" {
		return ""
	} else if strings.HasPrefix(bannerHash, "a_") {
		URL = animatedBannerURL
	} else {
		URL = staticBannerURL
	}

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

// IconURL returns the URL to get the avatar from a hash.
func IconURL(iconHash, staticIconURL, animatedIconURL, size string) string {
	var URL string
	if iconHash == "" {
		return ""
	} else if strings.HasPrefix(iconHash, "a_") {
		URL = animatedIconURL
	} else {
		URL = staticIconURL
	}

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

var QuoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

const DiscordEpoch int64 = 1420070400000

// SnowflakeTimestamp returns the creation time of a Snowflake ID relative to the creation of Discord.
func SnowflakeTimestamp(id uint64) time.Time {
	timestamp := int64(id>>22) + DiscordEpoch
	return time.Unix(0, timestamp*int64(time.Millisecond))
}
