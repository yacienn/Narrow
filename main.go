package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket upgrade error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Player connected!")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			return
		}

		fmt.Println("received:", string(message))

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/game", gameHandler)
	fmt.Println("Server listening on :8080")

	// Listen on all network interfaces, not just localhost,
	// so other PCs on the LAN can reach it.
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
