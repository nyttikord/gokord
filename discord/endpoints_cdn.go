package discord

import "fmt"

var (
	EndpointCDN             = "https://cdn.discordapp.com"
	EndpointCDNAttachments  = EndpointCDN + "/attachments"
	EndpointCDNAvatars      = EndpointCDN + "/avatars"
	EndpointCDNIcons        = EndpointCDN + "/icons"
	EndpointCDNSplashes     = EndpointCDN + "/splashes"
	EndpointCDNChannelIcons = EndpointCDN + "/channel-icons"
	EndpointCDNBanners      = EndpointCDN + "/banners"
	EndpointCDNGuilds       = EndpointCDN + "/guilds"
	EndpointCDNRoleIcons    = EndpointCDN + "/role-icons"
)

func endpointCDN(base string, id uint64, hash string) string {
	return fmt.Sprintf("%s/%d/%s", base, id, hash)
}

func EndpointUserAvatar(uID uint64, aID string) string {
	return endpointCDN(EndpointCDNAvatars, uID, aID) + ".png"
}
func EndpointUserAvatarAnimated(uID uint64, aID string) string {
	return endpointCDN(EndpointCDNAvatars, uID, aID) + ".gif"
}
func EndpointDefaultUserAvatar(idx uint64) string {
	return fmt.Sprintf("%s/embed/avatars/%d.png", EndpointCDN, idx)
}
func EndpointUserBanner(uID uint64, cID string) string {
	return endpointCDN(EndpointCDNBanners, uID, cID) + ".png"
}
func EndpointUserBannerAnimated(uID uint64, cID string) string {
	return endpointCDN(EndpointCDNBanners, uID, cID) + ".gif"
}

func EndpointGuildIcon(gID uint64, hash string) string {
	return endpointCDN(EndpointCDNIcons, gID, hash) + ".png"
}
func EndpointGuildIconAnimated(gID uint64, hash string) string {
	return endpointCDN(EndpointCDNIcons, gID, hash) + ".gif"
}
func EndpointGuildSplash(gID uint64, hash string) string {
	return endpointCDN(EndpointCDNSplashes, gID, hash) + ".png"
}
func EndpointGuildBanner(gID uint64, hash string) string {
	return endpointCDN(EndpointCDNBanners, gID, hash) + ".png"
}
func EndpointGuildBannerAnimated(gID uint64, hash string) string {
	return endpointCDN(EndpointCDNBanners, gID, hash) + ".gif"
}

func endpointCDNGuildMember(gID, uID uint64, sub, hash string) string {
	return fmt.Sprintf("%s/%d/users/%d/%s/%s", EndpointCDNGuilds, gID, uID, sub, hash)
}
func EndpointGuildMemberAvatar(gId, uID uint64, aID string) string {
	return endpointCDNGuildMember(gId, uID, "avatars", aID) + ".png"
}
func EndpointGuildMemberAvatarAnimated(gId, uID uint64, aID string) string {
	return endpointCDNGuildMember(gId, uID, "avatars", aID) + ".gif"
}
func EndpointGuildMemberBanner(gId, uID uint64, hash string) string {
	return endpointCDNGuildMember(gId, uID, "banners", hash) + ".png"
}
func EndpointGuildMemberBannerAnimated(gId, uID uint64, hash string) string {
	return endpointCDNGuildMember(gId, uID, "banners", hash) + ".gif"
}

func EndpointRoleIcon(rID uint64, hash string) string {
	return endpointCDN(EndpointCDNRoleIcons, rID, hash) + ".png"
}
