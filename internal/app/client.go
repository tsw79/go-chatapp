package app

import (
	"time"

	"github.com/gorilla/websocket"
)

/* Client represents a single chatting user. */
type client struct {
	socket   *websocket.Conn        // socket is the web socket for this client.
	send     chan *message          // send is a channel on which messages are sent.
	room     *Room                  // room is the room this client is chatting in.
	userData map[string]interface{} // userData holds information about the user
}

/* Read from socket, whilst looping, sending any received messages to the forward channel */
func (clt *client) read() {
	defer clt.socket.Close()
	for {
		var msg *message
		if err := clt.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = clt.userData["name"].(string)
		if avatarURL, ok := clt.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		clt.room.forward <- msg
	}
}

/* Write from socket, whilst looping, continually accepts messages and writing it out of the coket */
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
