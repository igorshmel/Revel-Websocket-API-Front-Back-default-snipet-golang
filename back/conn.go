package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stixi/back/api"
	"stixi/back/hlp"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"stixi/back/bom"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type IniMySQL struct {
	Database string
	User     string
	Pass     string
	Host     string
	Port     string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The whlpebsocket connection.
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		if err := c.ws.Close(); err != nil {
			hlp.Add2Log("conn.go", "error: readPump: ws.Close failed")
		}
		hlp.Add2Log("conn.go", "ws: close")
	}()
	c.ws.SetReadLimit(maxMessageSize)
	if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		hlp.Add2Log("conn.go", "error: ws.SetReadDeadline failed")
	}
	c.ws.SetPongHandler(func(string) error {
		if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return err
		}
		return nil
	})
	for {
		if _, message, err := c.ws.ReadMessage(); err == nil {
			req(message, c)
		} else {
			break
		}
	}
}

func req(mess []byte, co *connection) {
	var msql IniMySQL
	// Чтение переменных из файла настроек ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	i := hlp.InitS{File: "config.ini", Sect: "DB"}
	if err := i.InitF(&msql); err != nil {
		hlp.Add2Log("conn.go", "fail to read ini: "+fmt.Sprint(err))
	}

	if db, err := gorm.Open("mysql", msql.User+":"+msql.Pass+"@/"+msql.Database); err == nil {
		defer func() {
			if err := db.Close(); err != nil {
				hlp.Add2Log("conn.go", "error close db "+fmt.Sprint(err))
			}
		}()

		var dat map[string]interface{}
		if err := json.Unmarshal(mess, &dat); err == nil {
			hlp.Add2Log("okey", fmt.Sprint(dat))
			point := dat["point"].(string)

			if len(point) > 2 {
				h.unregister <- co
				if err = co.ws.Close(); err != nil {
					hlp.Add2Log("conn.go", "error: default: ws.Close failed")
				}
			} else {
				req := api.ReqApiDefaultStruct{}
				req.D = mess
				rpl := api.RplApiDefaultStruct{}

				back := bom.Backend()
				back.DoIt(hlp.AnyToByte(req), &rpl, point, db)

				co.send <- rpl.Rpl
			}

		} else {
			hlp.Add2Log("conn.go", "unmarshal err: "+fmt.Sprint(err))
			h.unregister <- co
			if err = co.ws.Close(); err != nil {
				hlp.Add2Log("conn.go", "error: ws.Close failed")
			}
		}
	} else {
		hlp.Add2Log("conn.go", "failed to connect database")
		h.unregister <- co
		if err = co.ws.Close(); err != nil {
			hlp.Add2Log("conn.go", "error: ws.Close failed")
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		hlp.Add2Log("conn.go", "error: ws.SetWriteDeadline failed")
	}
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.ws.Close(); err != nil {
			hlp.Add2Log("conn.go", "error: writePump ws.Close failed")
		}
		hlp.Add2Log("conn.go", "ws & ticker: close")
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
					hlp.Add2Log("conn.go", "error: write websocket.CloseMessage failed")
				}
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	if ws, err := upgrader.Upgrade(w, r, nil); err == nil {
		c := &connection{send: make(chan []byte, 256), ws: ws}
		h.register <- c
		go c.writePump()
		c.readPump()
	} else {
		hlp.Add2Log("conn.go", fmt.Sprint(err))
		return
	}
}
