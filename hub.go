package wxws

import (
	"encoding/json"
	"log"
)

type Hub struct {
	Broadcast  chan Message
	Register   chan *Subscription
	Unregister chan *Subscription
	Channels   map[string]map[*Connection]bool
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Subscription),
		Unregister: make(chan *Subscription),
		Channels:   make(map[string]map[*Connection]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connectionMap := h.Channels[s.Channel]
			if connectionMap == nil {
				connectionMap = make(map[*Connection]bool)
				h.Channels[s.Channel] = connectionMap
			}
			h.Channels[s.Channel][s.Conn] = true

		case s := <-h.Unregister:
			connectionMap := h.Channels[s.Channel]
			if connectionMap != nil {
				if connectionMap[s.Conn] {
					delete(connectionMap, s.Conn)
					if len(connectionMap) == 0 {
						delete(h.Channels, s.Channel)
					}
				}
			}

		case m := <-h.Broadcast:
			// generate string of bytes from message data interface
			json, err := json.Marshal(map[string]interface{}{m.Channel: m.Data})
			log.Println("broadcasting data", string(json), "on channel", m.Channel)
			if err != nil {
				log.Println("JSON MARSHAL ERROR:", err)
			} else {
				// send it to each connection
				connectionMap := h.Channels[m.Channel]
				for c := range connectionMap {
					c.Send <- json
				}
			}
		}
	}
}
