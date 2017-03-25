package types

// MessageType is the type of message received from the bot
type MessageType string

const (
	// MessageTypeUnknown is an unknown message
	MessageTypeUnknown MessageType = "unknown"
	// MessageTypeAdd is for adding apps
	MessageTypeAdd MessageType = "add"
	// MessageTypeDelete is for deleting apps
	MessageTypeDelete MessageType = "delete"
	// MessageTypeList is for listing apps
	MessageTypeList MessageType = "list"
	// MessageTypeSearch is for searching apps
	MessageTypeSearch MessageType = "search"
	// MessageTypeSync is for syncing
	MessageTypeSync MessageType = "sync"
	// MessageTypeAnnounce is for announcing app news
	MessageTypeAnnounce MessageType = "announce"
	// MessageTypeHelp is for help
	MessageTypeHelp MessageType = "help"
)
