package ws

import (
	"log"
	"net/http"
	"pulse/repository"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	userID int64
	conn   *websocket.Conn
}

type Message struct {
	Type        string `json:"type"`
	FollowerID  string `json:"followerId,omitempty"`
	FollowingID string `json:"followingId,omitempty"`
}

type Handler struct {
	hub  *Hub
	user *repository.UserRepository
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewClient(conn *websocket.Conn, userID int64) *Client {
	return &Client{
		conn:   conn,
		userID: userID,
	}
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	userID := c.GetInt64("userID")
	// user, userErr := h.user.Me(userID)
	// if userErr != nil {
	// 	return
	// }
	client := NewClient(conn, userID)
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("ws disconnect:", err)
			return
		}
	}
}
