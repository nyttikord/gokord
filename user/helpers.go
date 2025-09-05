package user

import "strings"

func avatarURL(avatarHash, defaultAvatarURL, staticAvatarURL, animatedAvatarURL, size string) string {
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

func bannerURL(bannerHash, staticBannerURL, animatedBannerURL, size string) string {
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
