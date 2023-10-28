package ws

import (
	"chat-app/repository"
	"log"
	"time"
)

const TYPE_BOOTUP = "bootup"
const TYPE_CHAT = "chat"
const TYPE_CHANNEL_ADDED = "channel_added"

type Broadcastable interface {
	Setup(*Client, *WebsocketProvider) bool
	Broadcast(*WebsocketProvider)
}

type WebsocketMessage struct {
	Type string `json:"type"`
}

type BootupRequest struct {
	WebsocketMessage
	Username string `json:"username"`
}

type BootupResponse struct {
	WebsocketMessage
	Username string `json:"username"`
}

func (br *BootupRequest) Setup(client *Client, wp *WebsocketProvider) bool {
	client.Username = br.Username
	wp.UserClientMap[client.Username] = client
	return true
}

func (bm BootupRequest) Broadcast(wp *WebsocketProvider) {
	response := BootupResponse{}
	response.Type = TYPE_BOOTUP
	response.Username = bm.Username

	if user, exists := (*wp).UserClientMap[bm.Username]; exists {
		if err := user.Connection.WriteJSON(response); err != nil {
			log.Println("error sending message to recepient", user.Connection.RemoteAddr())
			wp.deleteClient(user)
		}
	}
}

type ChatRequest struct {
	WebsocketMessage
	Message repository.Message `json:"message"`
	Channel repository.Channel
}

type ChatResponse struct {
	WebsocketMessage
	Message repository.Message `json:"message"`
}

func (cr *ChatRequest) Setup(client *Client, wp *WebsocketProvider) bool {
	channel, err := wp.DB.GetChannel(cr.Message.ChannelID)
	if err != nil {
		return false
	}

	if len(cr.Message.Content) == 0 || len(cr.Message.From) == 0 || len(channel.Members) == 0 {
		return false
	}

	id, err := wp.DB.CreateMessage(cr.Message.From, cr.Message.ChannelID, cr.Message.Content)
	if err != nil {
		log.Println("error creating message db entry", client.Connection.RemoteAddr())
		return false
	}

	cr.Message.ID = id
	cr.Message.Timestamp = time.Now().UnixMilli()
	cr.Channel = channel
	return true
}

func (cr ChatRequest) Broadcast(wp *WebsocketProvider) {
	response := ChatResponse{}
	response.Type = TYPE_CHAT
	response.Message = cr.Message

	for _, recepient := range cr.Channel.Members {
		if to, exists := (*wp).UserClientMap[recepient]; exists {
			if err := to.Connection.WriteJSON(response); err != nil {
				log.Println("error sending message to recepient", to.Connection.RemoteAddr())
				wp.deleteClient(to)
			}
		}
	}
}

type ChannelAddedRequest struct {
	WebsocketMessage
	Username  string `json:"username"`
	ChannelID string `json:"channelID"`
	Channel   repository.Channel
}

type ChannelAddedResponse struct {
	WebsocketMessage
	ChannelID string `json:"channelID"`
}

func (car *ChannelAddedRequest) Setup(client *Client, wp *WebsocketProvider) bool {
	exists, err := wp.DB.UserExists(car.Username)
	if err != nil || !exists {
		return false
	}

	channel, err := wp.DB.GetChannel(car.ChannelID)
	if err != nil {
		return false
	}

	car.Channel = channel
	return true
}

func (car ChannelAddedRequest) Broadcast(wp *WebsocketProvider) {
	response := ChannelAddedResponse{}
	response.Type = TYPE_CHANNEL_ADDED
	response.ChannelID = car.ChannelID

	for _, recepient := range car.Channel.Members {
		if recepient == car.Username {
			continue
		}
		if to, exists := (*wp).UserClientMap[recepient]; exists {
			if err := to.Connection.WriteJSON(response); err != nil {
				log.Println("error sending message to recepient", to.Connection.RemoteAddr())
				wp.deleteClient(to)
			}
		}
	}
}
