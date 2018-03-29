package core

// A Permission is an internal or chat permission that can be granted.
type Permission uint8

// Permissions
const (
	// Bot: Administrative, meta
	PermDeveloper Permission = iota
	PermUnknown
	PermFull

	// User: General
	PermInvite
	PermView

	// User: Administrative
	PermAdmin
	PermKick
	PermBan
	PermChannel
	PermGroup
	PermAuditLogs
	PermRoles
	PermPermissions
	PermWebhooks
	PermEmoji

	// Message
	PermSend
	PermRead
	PermTTS
	PermDelete
	PermEmbed
	PermUpload
	PermHistory
	PermMassMention
	PermForeignEmoji
	PermReact

	// User: Voice
	PermVoiceConnect
	PermVoiceTalk
	PermVoiceMute
	PermVoiceDeafen
	PermVoiceActivity

	// User: Personal
	PermNickname
	PermManageNicknames
)
