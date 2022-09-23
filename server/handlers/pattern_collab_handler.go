package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	patternId string
	send      chan []byte
}

type Hub struct {
	patterns                map[string]map[*Client]bool
	broadcast               chan *SocketMessage
	register                chan *Client
	unregister              chan *Client
	updateSubscribedPattern chan *PatternUpdateRequest
}

type SocketMessage struct {
	Type   string      `json:"type"`
	Topics []string    `json:"topics,omitempty"`
	Topic  string      `json:"topic,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type PatternUpdateRequest struct {
	oldPatternId string
	newPatternId string
	client       *Client
}

func (c *Client) processMessage(message []byte) {
	sockMessage := &SocketMessage{}
	err := json.Unmarshal(message, sockMessage)
	if err != nil {
		return
	}
	switch sockMessage.Type {
	case "ping":
		pongMessage := &SocketMessage{Type: "pong"}
		msg, err := json.Marshal(pongMessage)
		if err != nil {
			return
		}
		c.send <- msg
	case "subscribe":
		newPatternId := sockMessage.Topics[0]
		updateRequest := &PatternUpdateRequest{
			oldPatternId: c.patternId,
			newPatternId: newPatternId,
			client:       c,
		}
		c.patternId = newPatternId
		c.hub.updateSubscribedPattern <- updateRequest
	case "unsubscribe":
		oldPatternId := sockMessage.Topics[0]
		updateRequest := &PatternUpdateRequest{
			oldPatternId: oldPatternId,
			newPatternId: "",
			client:       c,
		}
		c.patternId = ""
		c.hub.updateSubscribedPattern <- updateRequest
	case "publish":
		c.hub.broadcast <- sockMessage
	}
}

func (c *Client) readMessages() {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		c.processMessage(message)
	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Hub) unsubscribeClientFromPattern(patternId string, client *Client) {
	if h.patterns[patternId] != nil {
		delete(h.patterns[patternId], client)
		if len(h.patterns[patternId]) == 0 {
			delete(h.patterns, patternId)
		}
	}
}

func (h *Hub) subscribeClientToPattern(patternId string, client *Client) {
	if h.patterns[patternId] == nil {
		h.patterns[patternId] = make(map[*Client]bool)
	}
	h.patterns[patternId][client] = true
}

func (h *Hub) closeChannel(ch chan []byte) {
	select {
	case <-ch:
		return
	default:
	}
	close(ch)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.subscribeClientToPattern(client.patternId, client)
		case client := <-h.unregister:
			h.unsubscribeClientFromPattern(client.patternId, client)
			h.closeChannel(client.send)
		case patternUpdateRequest := <-h.updateSubscribedPattern:
			h.unsubscribeClientFromPattern(patternUpdateRequest.oldPatternId, patternUpdateRequest.client)
			h.subscribeClientToPattern(patternUpdateRequest.newPatternId, patternUpdateRequest.client)
		case message := <-h.broadcast:
			if h.patterns[message.Topic] != nil {
				msg, err := json.Marshal(message)
				if err == nil {
					for client := range h.patterns[message.Topic] {
						select {
						case client.send <- msg:
						default:
							delete(h.patterns[message.Topic], client)
							h.closeChannel(client.send)
						}
					}
				}
			}
		}
	}
}

func newHub() *Hub {
	return &Hub{
		patterns:                make(map[string]map[*Client]bool),
		broadcast:               make(chan *SocketMessage),
		register:                make(chan *Client),
		unregister:              make(chan *Client),
		updateSubscribedPattern: make(chan *PatternUpdateRequest),
	}
}

var mainHub *Hub = nil

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{hub: hub, conn: conn, patternId: "", send: make(chan []byte)}
	client.hub.register <- client

	go client.readMessages()
	go client.writeMessages()
}

func (h *Handler) PatternCollabHandler(w http.ResponseWriter, r *http.Request) {
	if mainHub == nil {
		mainHub = newHub()
		go mainHub.run()
	}
	serveWs(mainHub, w, r)
}
