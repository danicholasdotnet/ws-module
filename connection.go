package wxws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID   int
	WS   *websocket.Conn
	Send chan []byte
}

func (c *Connection) Write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(WriteWait))
	return c.WS.WriteMessage(mt, payload)
}

func (c *Connection) ReadLoop(h *Hub) {
	defer func() {
		c.UnregisterAll(h)
		c.WS.Close()
		delete(Population, c.ID)
		h.Broadcast <- Message{
			Channel: "population",
			Data:    Population,
		}
	}()

	// set limits
	c.WS.SetReadLimit(MaxMessageSize)
	c.WS.SetReadDeadline(time.Now().Add(PongWait))
	c.WS.SetPongHandler(func(string) error {
		c.WS.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		// get slice of bytes from ws connection
		msgType, msgData, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("WS READ ERR:", err)
			}
			break
		}
		if msgType != 1 {
			log.Println("Incoming message was binary, not slice of bytes.")
			break
		}

		// get packet struct from slice of bytes
		log.Println("receiving data", string(msgData), "from connection", c.WS.RemoteAddr())
		p := &Packet{}
		json.Unmarshal(msgData, p)

		// update subscriptions to match interests
		c.UpdateSubs(h, p.Interests)

		// use the relevant handler for the event
		handler := HandlersArray[p.Event.Name]
		if handler == nil {
			log.Println("EVENT SENT DID NOT HAVE A HANDLER ASSOCIATED")
			log.Println("event name in sent packet:", p.Event.Name)
			log.Println("events that are registered:", HandlersArray)
		} else {
			// construct outgoing message based on incoming packet
			m, err := handler(p, c)
			if err != nil {
				log.Println("HANDLER ERROR:", err)
			} else {
				// broadcast message to anyone subscribed
				h.Broadcast <- m
			}
		}
	}
}

func (c *Connection) WriteLoop(h *Hub) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.UnregisterAll(h)
		c.WS.Close()
		close(c.Send)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Write(websocket.TextMessage, message)
			if err != nil {
				log.Println("WS WRITE ERR:", err)
				return
			}
		case <-ticker.C:
			err := c.Write(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("WS PING ERR:", err)
				return
			}
		}
	}
}

// unregister from every channel
func (c *Connection) UnregisterAll(h *Hub) {
	for ch := range h.Channels {
		s := &Subscription{
			Conn:    c,
			Channel: ch,
		}
		h.Unregister <- s
	}
}

func (c *Connection) UpdateSubs(h *Hub, interests []string) {
	// for each channel that exists currently
	for ch := range h.Channels {
		s := &Subscription{
			Conn:    c,
			Channel: ch,
		}
		present, index := StrSearch(ch, interests)
		// if the channel is one we want, register, else unregister
		if present {
			h.Register <- s
			interests = SliceRemove(interests, index)
		} else {
			h.Unregister <- s
		}
	}
	// register for all remaining channels, they don't exist yet
	for _, ch := range interests {
		h.Register <- &Subscription{
			Conn:    c,
			Channel: ch,
		}
	}
}
