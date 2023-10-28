package repository

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channelID"`
	From      string `json:"from"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type Channel struct {
	ID           string   `json:"id"`
	Members      []string `json:"members"`
	LastActivity Message  `json:"lastActivity"`
}
