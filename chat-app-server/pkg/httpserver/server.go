package httpserver

import (
	"chat-app/repository"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type HTTPProvider struct {
	DB     repository.DatabaseService
	Router *mux.Router
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func hpProvider(hp HTTPProvider, handler func(HTTPProvider, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(hp, w, r)
	}
}

func (hp HTTPProvider) Start() {
	hp.Router = mux.NewRouter()

	hp.Router.HandleFunc("/", StatusHandler).Methods(http.MethodGet)
	hp.Router.HandleFunc("/user/register", hpProvider(hp, UserRegisterHandler)).Methods(http.MethodPost)
	hp.Router.HandleFunc("/user/login", hpProvider(hp, UserLoginHandler)).Methods(http.MethodPost)
	hp.Router.HandleFunc("/user/channels", hpProvider(hp, UserChannelsHandler)).Methods(http.MethodPost)
	hp.Router.HandleFunc("/channels/get", hpProvider(hp, ChannelGetHandler)).Methods(http.MethodPost)
	hp.Router.HandleFunc("/channels/create", hpProvider(hp, ChannelCreateHandler)).Methods(http.MethodPost)
	hp.Router.HandleFunc("/channels/search", hpProvider(hp, ChannelSearchMessagesHandler)).Methods(http.MethodPost)

	hp.Router.Use(headerMiddleware)
	handler := cors.Default().Handler(hp.Router)
	http.ListenAndServe(":8080", handler)
}

func NewHTTPProvider(db repository.DatabaseService) (HTTPProvider, error) {
	return HTTPProvider{DB: db}, nil
}
