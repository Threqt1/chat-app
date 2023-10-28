package repository

type DatabaseService interface {
	CreateUser(string, string) (User, error)
	GetUser(string) (User, error)
	UserExists(string) (bool, error)
	CheckPassword(string, string) (bool, error)
	CreateChannel([]string) (Channel, error)
	GetChannel(string) (Channel, error)
	ChannelExists(string) (bool, error)
	GetChannels(string, int, int) ([]Channel, error)
	UpdateChannels(string, string, float64) error
	CreateMessage(string, string, string) (string, error)
	GetMessage(string) (Message, error)
	SearchMessages(string, int, int) ([]Message, error)
}

type SnowflakeService interface {
	// GenerateSnowflake generates a unique Twitter snowflake ID.
	GenerateSnowflake() string
}

type HTTPService interface {
	// Start the HTTP server.
	Start()
}

type WebsocketService interface {
	// Start the Websocket server.
	Start()
}
