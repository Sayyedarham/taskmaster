package handler

import (
	"net/http"

	wsinternal "taskmaster/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub *wsinternal.Hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Restrict in production
	},
}

func NewWSHandler(hub *wsinternal.Hub) *WSHandler {
	return &WSHandler{hub}
}

func (h *WSHandler) Handle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &wsinternal.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register(client)
	go client.WritePump()
	go client.ReadPump()
}
