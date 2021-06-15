package app

import (
	"chatapp/internal/lib/config"
	"net/http"

	"github.com/markbates/goth"
)

type ChatUserInterface interface {
	UniqueID() string
	AvatarURL() string
}

type ChatUser struct {
	goth.User
	UniqueId string
}

func (user ChatUser) UniqueID() string {
	return user.UniqueId
}

func (user ChatUser) AvatarURL() string {
	return user.User.AvatarURL
}

type authHandler struct {
	next http.Handler
}

//
func (h *authHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// check if user is logged in
	_, err := req.Cookie(config.GetInstance().Auth.Cookie)
	if err == http.ErrNoCookie {
		// not authenticated
		rw.Header().Set("Location", "/login")
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// success - call the next handler
	h.next.ServeHTTP(rw, req)
}

// Wraps other `http.Handler`s within `authHandler`
func ForceAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
