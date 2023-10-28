package httpserver

import (
	"chat-app/repository"
	"encoding/json"
	"net/http"
)

type UserChannelsRequest struct {
	Username string `json:"username"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
}

type UserChannelsResponse struct {
	Channels []repository.Channel `json:"channels"`
}

func UserChannelsHandler(hp HTTPProvider, w http.ResponseWriter, r *http.Request) {
	req := UserChannelsRequest{}
	res := UserChannelsResponse{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, repository.ErrorFailedToDecode.Error(), http.StatusBadRequest)
		return
	}

	exists, err := hp.DB.UserExists(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, repository.ErrorUsernameDoesntExist.Error(), http.StatusBadRequest)
		return
	}

	channels, err := hp.DB.GetChannels(req.Username, req.Start, req.End)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Channels = channels
	json.NewEncoder(w).Encode(res)
}
