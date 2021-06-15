package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"chatapp/internal/app"
	"chatapp/internal/app/handlers"
	handler "chatapp/internal/app/handlers"
	"chatapp/internal/lib"
	"chatapp/internal/lib/config"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/tsw79/debug/trace"
)

// Bootstrap the application and start the server
func Run() {
	var port = flag.String("port", "8085", "The address the server communicates.")
	flag.Parse()

	// load app defaults
	load()
	// replace the specifier with the correct addr
	tcpAddress := fmt.Sprintf("%s:%s", config.GetInstance().Server.Host, *port)
	serverAddr := fmt.Sprintf("http://%s", tcpAddress)
	// setup app
	setup(serverAddr)
	// create a new room
	room := app.NewRoom()
	// assign new Trace object to room
	room.Tracer = trace.New(os.Stdout)
	// Register URL paths
	registerHandlers(room)
	// run room inside a Go routine
	go room.Run()
	// start the web server
	startServer(tcpAddress)
}

// Initialize app
func load() {
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading `.env` file")
	}
	// Load config
	// config, err = lib.LoadConfig("./app.config.json")
	if err = config.LoadConfig("./app.config.json"); err != nil {
		log.Fatal("Error loading config file:", err)
	}
}

// Setup access to various OAuth2 service providers, and init db
func setup(serverAddr string) {
	// init db
	// db.DbConnect()

	maxAge := 86400 * 30 // 30 days
	// create a cookie
	store := sessions.NewCookieStore(
		[]byte(os.Getenv("AUTH_SECRET_HASH")),
	)
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = false  // Set to true when serving over https
	gothic.Store = store
	// Set up auth providers
	goth.UseProviders(
		facebook.New(
			os.Getenv("AUTH_FACEBOOK_CLIENT_ID"),
			os.Getenv("AUTH_FACEBOOK_CLIENT_SECRET"),
			serverAddr+os.Getenv("AUTH_FACEBOOK_REDIRECT_URL"),
		),
		github.New(
			os.Getenv("AUTH_GITHUB_CLIENT_ID"),
			os.Getenv("AUTH_GITHUB_CLIENT_SECRET"),
			serverAddr+os.Getenv("AUTH_GITHUB_REDIRECT_URL"),
		),
		google.New(
			os.Getenv("AUTH_GOOGLE_CLIENT_ID"),
			os.Getenv("AUTH_GOOGLE_CLIENT_SECRET"),
			serverAddr+os.Getenv("AUTH_GOOGLE_REDIRECT_URL"),
			"email",
			"profile",
		),
	)
}

// Register URL paths
func registerHandlers(room *app.Room) {
	http.Handle("/room", room)
	http.Handle("/login", &lib.TemplateHandler{
		Files: []string{"layout/app.auth.html", "user/login.html"},
		Data:  handlers.Page{Title: "Login"},
	})
	http.Handle("/register", &lib.TemplateHandler{
		Files: []string{"layout/app.auth.html", "user/register.html"},
		Data:  handlers.Page{Title: "Register"},
	})
	http.Handle("/verify", &lib.TemplateHandler{
		Files: []string{"layout/app.auth.html", "user/verify_email.html"},
		Data:  handlers.Page{Title: "Verify email"},
	})
	http.Handle("/chat", app.ForceAuth(&lib.TemplateHandler{
		Files: []string{"layout/app.html", "block/header.html", "chat.html"},
		Data:  handlers.Page{Title: "Chat Room"},
	}))
	http.Handle("/upload", app.ForceAuth(&lib.TemplateHandler{
		Files: []string{"layout/app.html", "block/header.html", "upload.html"},
		Data:  handlers.Page{Title: "Upload"},
	}))
	http.HandleFunc("/register_s", handler.RegistrationHandler)
	http.HandleFunc("/uploader", app.UploaderHandler)
	http.HandleFunc("/auth/", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/welcome_s", handler.WelcomeHandler)
	// Static handlers
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./web/uploads/avatars"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/static/css"))))
}

// Start HTTP server
func startServer(addr string) {
	log.Print("Starting HTTP server on ", addr)
	// start the web server
	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println("HTTP server running on ", addr)
		log.Fatal("ListenAndServe:", err)
		panic(err)
	}
	println("Running code after ListenAndServe (only happens when server shuts down..)")
}
