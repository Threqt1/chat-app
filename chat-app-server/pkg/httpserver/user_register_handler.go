package httpserver

import (
	"chat-app/repository"
	"encoding/json"
	"net/http"
)

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	User repository.User `json:"user"`
}

func UserRegisterHandler(hp HTTPProvider, w http.ResponseWriter, r *http.Request) {
	req := UserRegisterRequest{}
	res := UserRegisterResponse{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, repository.ErrorFailedToDecode.Error(), http.StatusBadRequest)
		return
	}

	unique, err := hp.DB.UserExists(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if unique {
		http.Error(w, repository.ErrorUsernameNotUnique.Error(), http.StatusBadRequest)
		return
	}

	user, err := hp.DB.CreateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.User = user
	res.User.Password = ""
	json.NewEncoder(w).Encode(res)
}
