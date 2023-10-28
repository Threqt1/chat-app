package database

// generateChannelsKey returns the key for the channels list sorted set in the database.
func generateChannelsKey(username string) string {
	return "channels#" + username
}

// generateChannelKey returns the key for the channel sorted set in the database.
func generateChannelKey(id string) string {
	return "channel#" + id
}

func generateChannelListKey(id string) string {
	return "channelList#" + id
}

// generateUsernameKey generates the key for setting a username in the database.
func generateUsernameKey(username string) string {
	return "user#" + username
}

// generateMessageKey generates the key for setting a message in the database.
func generateMessageKey(snowflake string) string {
	return "message#" + snowflake
}
