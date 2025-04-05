package websocket

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pingPeriod = 30 * time.Second
	pongWait   = 60 * time.Second
	writeWait  = 30 * time.Second
)

func SetHeartbeatConfig(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func PingConnection(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
			if err != nil {
				log.Println("Ping error:", err)
				cancel()
				conn.Close()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
