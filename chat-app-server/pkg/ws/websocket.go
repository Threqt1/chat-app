package ws

import (
	"chat-app/pkg/httpserver"
	"chat-app/repository"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn
	Username   string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebsocketProvider struct {
	DB            repository.DatabaseService
	ClientUserMap map[*Client]string
	UserClientMap map[string]*Client
	Broadcast     chan Broadcastable
}

// reciever handles recieving messages from websocket clients.
func (wp WebsocketProvider) reciever(client *Client) {
	for {
		_, rawMessage, err := client.Connection.ReadMessage()
		if err != nil {
			wp.deleteClient(client)
			return
		}

		message := WebsocketMessage{}

		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println("error marshalling message json", client.Connection.RemoteAddr())
			continue
		}

		var parsed Broadcastable

		switch message.Type {
		case TYPE_BOOTUP:
			parsed = &BootupRequest{}
		case TYPE_CHAT:
			parsed = &ChatRequest{}
		case TYPE_CHANNEL_ADDED:
			parsed = &ChannelAddedRequest{}
		}

		err = json.Unmarshal(rawMessage, parsed)
		if err != nil {
			log.Println("error marshalling message json", client.Connection.RemoteAddr())
			continue
		}

		if success := parsed.Setup(client, &wp); !success {
			continue
		}

		wp.Broadcast <- parsed
	}
}

// broadcaster broadcasts recieved messages to all applicable clients.
func (wp WebsocketProvider) broadcaster() {
	for {
		message := <-wp.Broadcast

		message.Broadcast(&wp)
	}
}

// serverWebsocket upgrades the HTTP connection to a websocket.
func (wp WebsocketProvider) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error registered websocket client")
	}

	client := &Client{Connection: websocket}
	wp.ClientUserMap[client] = ""

	wp.reciever(client)

	wp.deleteClient(client)
}

func (wp WebsocketProvider) deleteClient(client *Client) {
	client.Connection.Close()
	user := wp.ClientUserMap[client]
	delete(wp.ClientUserMap, client)
	delete(wp.UserClientMap, user)
}

func (wp WebsocketProvider) Start() {
	go wp.broadcaster()
	http.HandleFunc("/", httpserver.StatusHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wp.serveWebsocket(w, r)
	})
	http.ListenAndServe(":8081", nil)
}

func NewWebsocketProvider(db repository.DatabaseService) WebsocketProvider {
	wp := WebsocketProvider{}

	wp.DB = db
	wp.ClientUserMap = make(map[*Client]string)
	wp.UserClientMap = make(map[string]*Client)
	wp.Broadcast = make(chan Broadcastable)

	return wp
}
