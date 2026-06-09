package handler

import (
	"net/http"

	"taskmaster/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub *websocket.Hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Restrict in production
	},
}

func NewWSHandler(hub *websocket.Hub) *WSHandler {
	return &WSHandler{hub}
}

func (h *WSHandler) Handle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &websocket.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
