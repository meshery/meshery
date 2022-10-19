package handlers

import (
	"encoding/json"
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
	connected bool
	patternID string
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
	oldPatternID string
	newPatternID string
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
		if len(sockMessage.Topics) == 0 {
			return
		}
		newPatternID := sockMessage.Topics[0]
		updateRequest := &PatternUpdateRequest{
			oldPatternID: c.patternID,
			newPatternID: newPatternID,
			client:       c,
		}
		c.patternID = newPatternID
		c.hub.updateSubscribedPattern <- updateRequest
	case "unsubscribe":
		if len(sockMessage.Topics) == 0 {
			return
		}
		oldPatternID := sockMessage.Topics[0]
		updateRequest := &PatternUpdateRequest{
			oldPatternID: oldPatternID,
			newPatternID: "",
			client:       c,
		}
		c.patternID = ""
		c.hub.updateSubscribedPattern <- updateRequest
	case "publish":
		c.hub.broadcast <- sockMessage
	}
}

func (c *Client) readMessages() {
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.processMessage(message)
	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					return
				}

				_, err = w.Write(<-c.send)
				if err != nil {
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Hub) unsubscribeClientFromPattern(patternID string, client *Client) {
	if h.patterns[patternID] != nil {
		delete(h.patterns[patternID], client)
		if len(h.patterns[patternID]) == 0 {
			delete(h.patterns, patternID)
		}
	}
}

func (h *Hub) subscribeClientToPattern(patternID string, client *Client) {
	if h.patterns[patternID] == nil {
		h.patterns[patternID] = make(map[*Client]bool)
	}
	h.patterns[patternID][client] = true
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
			h.subscribeClientToPattern(client.patternID, client)
		case client := <-h.unregister:
			h.unsubscribeClientFromPattern(client.patternID, client)
			client.connected = false
			h.closeChannel(client.send)
		case patternUpdateRequest := <-h.updateSubscribedPattern:
			h.unsubscribeClientFromPattern(patternUpdateRequest.oldPatternID, patternUpdateRequest.client)
			h.subscribeClientToPattern(patternUpdateRequest.newPatternID, patternUpdateRequest.client)
		case message := <-h.broadcast:
			if h.patterns[message.Topic] != nil {
				msg, err := json.Marshal(message)
				if err == nil {
					for client := range h.patterns[message.Topic] {
						if client.connected {
							client.send <- msg
						} else {
							delete(h.patterns[message.Topic], client)
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
		broadcast:               make(chan *SocketMessage, 10),
		register:                make(chan *Client, 10),
		unregister:              make(chan *Client, 10),
		updateSubscribedPattern: make(chan *PatternUpdateRequest, 10),
	}
}

var mainHub *Hub

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{hub: hub, conn: conn, patternID: "", send: make(chan []byte, 100), connected: true}
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
