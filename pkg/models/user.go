package models

// User is a user of the gaia-bot
type User struct {
	// Github handle of the user
	Handle string `json:"handle"`
	// List of commands this user has access to
	Commands []string `json:"commands"`
}
