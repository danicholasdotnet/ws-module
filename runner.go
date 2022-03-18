package wxws

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = NewHub()

func Start() {
	go hub.Run()
	Listen()
}

func Listen() {
	http.HandleFunc("/", httpHandler)
	address := fmt.Sprint(Host, ":", Port)
	fmt.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
func httpHandler(w http.ResponseWriter, r *http.Request) {
	from := r.RemoteAddr
	page := r.URL

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("request from", from, "on", page, "=> error:", err)
		return
	}
	log.Println("request from", from, "on", page, "=> succesfully upgraded")

	c := &Connection{
		// value will be user id in future
		ID:   conn.RemoteAddr().(*net.TCPAddr).Port,
		Send: make(chan []byte, 256),
		WS:   conn,
	}

	go c.WriteLoop(hub)
	go c.ReadLoop(hub)
}
