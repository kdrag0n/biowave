// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains variables for all known Discord end points.  All functions
// throughout the Discordgo package use these variables for all connections
// to Discord.  These are all exported and you may modify them if needed.

package discordgo

// APIVersion is the Discord API version used for the REST and Websocket API.
var APIVersion = "6"

// Known Discord API Endpoints.
var (
	EndpointStatus     = "https://status.discordapp.com/api/v2/"
	EndpointSm         = EndpointStatus + "scheduled-maintenances/"
	EndpointSmActive   = EndpointSm + "active.json"
	EndpointSmUpcoming = EndpointSm + "upcoming.json"

	EndpointDiscord    = "https://discordapp.com/"
	EndpointAPI        = EndpointDiscord + "api/v" + APIVersion + "/"
	EndpointGuilds     = EndpointAPI + "guilds/"
	EndpointChannels   = EndpointAPI + "channels/"
	EndpointUsers      = EndpointAPI + "users/"
	EndpointGateway    = EndpointAPI + "gateway"
	EndpointGatewayBot = EndpointGateway + "/bot"
	EndpointWebhooks   = EndpointAPI + "webhooks/"

	EndpointCDN             = "https://cdn.discordapp.com/"
	EndpointCDNAttachments  = EndpointCDN + "attachments/"
	EndpointCDNAvatars      = EndpointCDN + "avatars/"
	EndpointCDNIcons        = EndpointCDN + "icons/"
	EndpointCDNSplashes     = EndpointCDN + "splashes/"
	EndpointCDNChannelIcons = EndpointCDN + "channel-icons/"

	EndpointAuth           = EndpointAPI + "auth/"
	EndpointLogin          = EndpointAuth + "login"
	EndpointLogout         = EndpointAuth + "logout"
	EndpointVerify         = EndpointAuth + "verify"
	EndpointVerifyResend   = EndpointAuth + "verify/resend"
	EndpointForgotPassword = EndpointAuth + "forgot"
	EndpointResetPassword  = EndpointAuth + "reset"
	EndpointRegister       = EndpointAuth + "register"

	EndpointVoice        = EndpointAPI + "/voice/"
	EndpointVoiceRegions = EndpointVoice + "regions"
	EndpointVoiceIce     = EndpointVoice + "ice"

	EndpointTutorial           = EndpointAPI + "tutorial/"
	EndpointTutorialIndicators = EndpointTutorial + "indicators"

	EndpointTrack        = EndpointAPI + "track"
	EndpointSso          = EndpointAPI + "sso"
	EndpointReport       = EndpointAPI + "report"
	EndpointIntegrations = EndpointAPI + "integrations"

	EndpointUser               = func(uID uint64) string { return EndpointUsers + idToStr(uID) }
	EndpointUserAvatar         = func(uID uint64, aID string) string { return EndpointCDNAvatars + idToStr(uID) + "/" + aID + ".png" }
	EndpointUserAvatarAnimated = func(uID uint64, aID string) string { return EndpointCDNAvatars + idToStr(uID) + "/" + aID + ".gif" }
	EndpointUserSettings       = func(uID uint64) string { return EndpointUsers + idToStr(uID) + "/settings" }
	EndpointUserGuilds         = func(uID uint64) string { return EndpointUsers + idToStr(uID) + "/guilds" }
	EndpointUserGuild          = func(uID, gID uint64) string { return EndpointUsers + idToStr(uID) + "/guilds/" + idToStr(gID) }
	EndpointUserGuildSettings  = func(uID, gID uint64) string {
		return EndpointUsers + idToStr(uID) + "/guilds/" + idToStr(gID) + "/settings"
	}
	EndpointUserChannels    = func(uID uint64) string { return EndpointUsers + idToStr(uID) + "/channels" }
	EndpointUserDevices     = func(uID uint64) string { return EndpointUsers + idToStr(uID) + "/devices" }
	EndpointUserConnections = func(uID uint64) string { return EndpointUsers + idToStr(uID) + "/connections" }
	EndpointUserNotes       = func(uID uint64) string { return EndpointUsers + "@me/notes/" + idToStr(uID) }

	EndpointGuild           = func(gID uint64) string { return EndpointGuilds + idToStr(gID) }
	EndpointGuildChannels   = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/channels" }
	EndpointGuildMembers    = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/members" }
	EndpointGuildMember     = func(gID, uID uint64) string { return EndpointGuilds + idToStr(gID) + "/members/" + idToStr(uID) }
	EndpointGuildMemberRole = func(gID, uID, rID uint64) string {
		return EndpointGuilds + idToStr(gID) + "/members/" + idToStr(uID) + "/roles/" + idToStr(rID)
	}
	EndpointGuildBans            = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/bans" }
	EndpointGuildBan             = func(gID, uID uint64) string { return EndpointGuilds + idToStr(gID) + "/bans/" + idToStr(uID) }
	EndpointGuildIntegrations    = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/integrations" }
	EndpointGuildIntegration     = func(gID, iID uint64) string { return EndpointGuilds + idToStr(gID) + "/integrations/" + idToStr(iID) }
	EndpointGuildIntegrationSync = func(gID, iID uint64) string {
		return EndpointGuilds + idToStr(gID) + "/integrations/" + idToStr(iID) + "/sync"
	}
	EndpointGuildRoles    = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/roles" }
	EndpointGuildRole     = func(gID, rID uint64) string { return EndpointGuilds + idToStr(gID) + "/roles/" + idToStr(rID) }
	EndpointGuildInvites  = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/invites" }
	EndpointGuildEmbed    = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/embed" }
	EndpointGuildPrune    = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/prune" }
	EndpointGuildIcon     = func(gID uint64, hash string) string { return EndpointCDNIcons + idToStr(gID) + "/" + hash + ".png" }
	EndpointGuildSplash   = func(gID uint64, hash string) string { return EndpointCDNSplashes + idToStr(gID) + "/" + hash + ".png" }
	EndpointGuildWebhooks = func(gID uint64) string { return EndpointGuilds + idToStr(gID) + "/webhooks" }

	EndpointChannel            = func(cID uint64) string { return EndpointChannels + idToStr(cID) }
	EndpointChannelPermissions = func(cID uint64) string { return EndpointChannels + idToStr(cID) + "/permissions" }
	EndpointChannelPermission  = func(cID, tID uint64) string { return EndpointChannels + idToStr(cID) + "/permissions/" + idToStr(tID) }
	EndpointChannelInvites     = func(cID uint64) string { return EndpointChannels + idToStr(cID) + "/invites" }
	EndpointChannelTyping      = func(cID uint64) string { return EndpointChannels + idToStr(cID) + "/typing" }
	EndpointChannelMessages    = func(cID uint64) string { return EndpointChannels + idToStr(cID) + "/messages" }
	EndpointChannelMessage     = func(cID, mID uint64) string { return EndpointChannels + idToStr(cID) + "/messages/" + idToStr(mID) }
	EndpointChannelMessageAck  = func(cID, mID uint64) string {
		return EndpointChannels + idToStr(cID) + "/messages/" + idToStr(mID) + "/ack"
	}
	EndpointChannelMessagesBulkDelete = func(cID uint64) string { return EndpointChannel(cID) + "/messages/bulk-delete" }
	EndpointChannelMessagesPins       = func(cID uint64) string { return EndpointChannel(cID) + "/pins" }
	EndpointChannelMessagePin         = func(cID, mID uint64) string { return EndpointChannel(cID) + "/pins/" + idToStr(mID) }

	EndpointGroupIcon = func(cID uint64, hash string) string { return EndpointCDNChannelIcons + idToStr(cID) + "/" + hash + ".png" }

	EndpointChannelWebhooks = func(cID uint64) string { return EndpointChannel(cID) + "/webhooks" }
	EndpointWebhook         = func(wID uint64) string { return EndpointWebhooks + idToStr(wID) }
	EndpointWebhookToken    = func(wID uint64, token string) string { return EndpointWebhooks + idToStr(wID) + "/" + token }

	EndpointMessageReactionsAll = func(cID, mID uint64) string {
		return EndpointChannelMessage(cID, mID) + "/reactions"
	}
	EndpointMessageReactions = func(cID, mID uint64, eID string) string {
		return EndpointChannelMessage(cID, mID) + "/reactions/" + eID
	}
	EndpointMessageReaction = func(cID, mID uint64, eID string, uID uint64) string {
		return EndpointMessageReactions(cID, mID, eID) + "/" + idToStr(uID)
	}

	EndpointRelationships       = func() string { return EndpointUsers + "@me" + "/relationships" }
	EndpointRelationship        = func(uID uint64) string { return EndpointRelationships() + "/" + idToStr(uID) }
	EndpointRelationshipsMutual = func(uID uint64) string { return EndpointUsers + idToStr(uID) + "/relationships" }

	EndpointGuildCreate = EndpointAPI + "guilds"

	EndpointInvite = func(iID uint64) string { return EndpointAPI + "invite/" + idToStr(iID) }

	EndpointIntegrationsJoin = func(iID uint64) string { return EndpointAPI + "integrations/" + idToStr(iID) + "/join" }

	EndpointEmoji = func(eID uint64) string { return EndpointAPI + "emojis/" + idToStr(eID) + ".png" }

	EndpointOauth2          = EndpointAPI + "oauth2/"
	EndpointApplications    = EndpointOauth2 + "applications"
	EndpointApplication     = func(aID uint64) string { return EndpointApplications + "/" + idToStr(aID) }
	EndpointApplicationsBot = func(aID uint64) string { return EndpointApplications + "/" + idToStr(aID) + "/bot" }
)
