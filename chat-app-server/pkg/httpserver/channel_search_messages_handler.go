package httpserver

import (
	"chat-app/repository"
	"encoding/json"
	"net/http"
)

type ChannelSearchMessagesRequest struct {
	ChannelID string `json:"channelID"`
	Start     int    `json:"start"`
	Stop      int    `json:"stop"`
}

type ChannelSearchMessagesResponse struct {
	Messages []repository.Message `json:"messages"`
}

func ChannelSearchMessagesHandler(hp HTTPProvider, w http.ResponseWriter, r *http.Request) {
	req := ChannelSearchMessagesRequest{}
	res := ChannelSearchMessagesResponse{}

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

	messages, err := hp.DB.SearchMessages(req.ChannelID, req.Start, req.Stop)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Messages = messages
	json.NewEncoder(w).Encode(res)
}
