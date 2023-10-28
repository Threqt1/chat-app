package httpserver

import (
	"chat-app/repository"
	"encoding/json"
	"net/http"
)

type ChannelGetRequest struct {
	ChannelID string `json:"channelID"`
}

type ChannelGetResponse struct {
	Channel repository.Channel `json:"channel"`
}

func ChannelGetHandler(hp HTTPProvider, w http.ResponseWriter, r *http.Request) {
	req := ChannelGetRequest{}
	res := ChannelGetResponse{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, repository.ErrorFailedToDecode.Error(), http.StatusBadRequest)
		return
	}

	exists, err := hp.DB.ChannelExists(req.ChannelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, repository.ErrorChannelDoesntExist.Error(), http.StatusBadRequest)
		return
	}

	channel, err := hp.DB.GetChannel(req.ChannelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Channel = channel
	json.NewEncoder(w).Encode(res)
}
