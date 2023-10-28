package httpserver

import (
	"chat-app/repository"
	"encoding/json"
	"net/http"
)

type ChannelCreateRequest struct {
	Members []string `json:"members"`
}

type ChannelCreateResponse struct {
	Channel repository.Channel `json:"channel"`
}

func ChannelCreateHandler(hp HTTPProvider, w http.ResponseWriter, r *http.Request) {
	req := ChannelCreateRequest{}
	res := ChannelCreateResponse{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, repository.ErrorFailedToDecode.Error(), http.StatusBadRequest)
		return
	}

	for _, member := range req.Members {
		exists, err := hp.DB.UserExists(member)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, repository.ErrorUsernameDoesntExist.Error(), http.StatusBadRequest)
			return
		}
	}

	channel, err := hp.DB.CreateChannel(req.Members)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Channel = channel
	json.NewEncoder(w).Encode(res)
}
