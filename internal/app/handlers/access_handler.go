package handlers

import (
	"chatapp/internal/app"
	service "chatapp/internal/app/services"
	"chatapp/internal/lib/config"
	"chatapp/internal/lib/util"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
)

// Logs out an authenticated user
func LogoutHandler(rw http.ResponseWriter, req *http.Request) {
	gothic.Logout(rw, req)

	// TODO Need to get rid of this and only use Gorilla's session
	http.SetCookie(rw, &http.Cookie{
		Name:   config.GetInstance().Auth.Cookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1, // expire and remove cokkie
	})
	rw.Header()["Location"] = []string{"/chat"}
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

// LoginHandler handles login attempts by both auth providers and local app's auth.
// 	Provider's format: /auth/{provider}/{action}
func LoginHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		// ChatApp login attempt
		handleAppLogin(rw, req)
	} else {
		// Provider login attempt
		slice := strings.Split(req.URL.Path, "/")
		// provider := slice[2]
		action := slice[3]
		switch action {
		case "login":
			gothic.BeginAuthHandler(rw, req)
		case "callback":
			user, err := gothic.CompleteUserAuth(rw, req)
			if err != nil {
				fmt.Fprintln(rw, err)
				return
			}
			completeAuthentication(rw, req, &user)
		default:
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, "Auth action %s not supported", action)
		}
	}
}

// handleAppLogin handles application's login process
func handleAppLogin(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		rw.WriteHeader(400) // Return 400 Bad Request.
		return
	}
	// Read the form data
	req.ParseForm()
	username := req.Form.Get("uname")
	password := req.Form.Get("passwd")

	// TODO Sanatize data...
	// TODO Authenticate data

	creds := &service.CredentialsDTO{
		Username: username,
		Password: password,
	}
	gothUser := goth.User{}
	// LoginService will handle the authentication
	loginService := service.NewAuthenticationService(creds, &gothUser)
	err := loginService.Execute()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	completeAuthentication(rw, req, &gothUser)
}

// Process the Provider's Callback for the Login request
func completeAuthentication(rw http.ResponseWriter, req *http.Request, user *goth.User) {
	chtUsr := &app.ChatUser{User: *user}
	md5Hash := md5.New()
	io.WriteString(md5Hash, strings.ToLower(user.Email))
	chtUsr.UniqueId = fmt.Sprintf("%x", md5Hash.Sum(nil))
	avatarURL, err := app.StrategyList.GetAvatarURL(*chtUsr)
	if err != nil {
		log.Fatalln("Error when trying to access GetAvatarURL", "-", err)
	}
	// TODO Need to get rid of this and only use Gorilla's session
	// save some data
	authCookieValue := objx.New(map[string]interface{}{
		"userid":     chtUsr.UniqueId,
		"name":       user.Name,
		"avatar_url": avatarURL,
		"email":      user.Email,
	}).MustBase64()
	http.SetCookie(rw, &http.Cookie{
		Name:  config.GetInstance().Auth.Cookie,
		Value: authCookieValue,
		Path:  "/",
	})
	// Redirecting
	rw.Header()["Location"] = []string{"/chat"}
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

// Handler for welcoming new users
func WelcomeHandler(rw http.ResponseWriter, req *http.Request) {

	// TODO Do this check with every auth request??
	// Extract the access token from the cookie
	token, errCode, err := retrieveTokenFromCookie(req)
	// Check token
	if err != nil {
		rw.WriteHeader(errCode)
		fmt.Println(err)
		// Redirect to login ??
	}
	// Validate token
	claims, err := util.ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// Renew token only if it is not within the expired time
	if !util.TokenIsWithinExpiry(&claims) {
		util.RenewToken(claims.Username)
	}
	rw.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

// Extract access token from cookie
func retrieveTokenFromCookie(req *http.Request) (string, int, error) {
	// Get the request's cookies
	cookie, err := req.Cookie(config.GetInstance().Auth.Cookie)
	switch err {
	case http.ErrNoCookie:
		fmt.Println(http.StatusText(http.StatusUnauthorized) + "Cookie not found!")
		return "", http.StatusUnauthorized, err
	case nil:
		fmt.Println(http.StatusText(http.StatusBadRequest) + "Cookie is nil")
		return "", http.StatusBadRequest, err
	}
	// Get the token string from the cookie
	if token := cookie.Value; token == "" {
		fmt.Println(http.StatusText(http.StatusUnauthorized) + "Token not found!")
		return "", http.StatusBadRequest, err
	} else {
		return token, 0, nil
	}
}
