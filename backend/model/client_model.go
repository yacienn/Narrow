package model

import "github.com/gorilla/websocket"

type Client struct {
	Id   string
	Conn *websocket.Conn
}
