package httpserver

import (
	"chat-app/repository"
	"encoding/json"
	"net/http"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	User repository.User `json:"user"`
}

func UserLoginHandler(hp HTTPProvider, w http.ResponseWriter, r *http.Request) {
	req := UserLoginRequest{}
	res := UserLoginResponse{}

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

	valid, err := hp.DB.CheckPassword(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, repository.ErrorInvalidCredentials.Error(), http.StatusBadRequest)
		return
	}

	user, err := hp.DB.GetUser(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.User = user
	res.User.Password = ""
	json.NewEncoder(w).Encode(res)
}
