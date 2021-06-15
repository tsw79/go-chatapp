package app

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"github.com/tsw79/debug/trace"
)

type Room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan *message
	// Use two channels to avoid multiple clients from accessing the same data at the same time!
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// debgug and trace activity in the room
	Tracer trace.DebugTracer
}

/* Makes a new room */
func NewRoom() *Room {
	return &Room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		// initiate newRoom with a nill tracer
		Tracer: trace.Off(),
	}
}

func (rm *Room) Run() {
	for {
		select {
		case client := <-rm.join:
			// joining
			rm.clients[client] = true
			rm.Tracer.Trace("New client joined.")
		case client := <-rm.leave:
			// leaving
			delete(rm.clients, client)
			close(client.send)
			rm.Tracer.Trace("Client left.")
		case msg := <-rm.forward:
			rm.Tracer.Trace("Message received: ", msg.Message)
			// forward message to all clients
			for client := range rm.clients {
				client.send <- msg
				rm.Tracer.Trace(" --- sent to client.")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

/* Turning room into an HTTP handler */
func (rm *Room) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     rm,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	rm.join <- client
	defer func() { rm.leave <- client }()
	go client.write()
	client.read()
}
