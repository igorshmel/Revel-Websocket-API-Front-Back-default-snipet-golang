package controllers

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/revel/revel"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

type App struct {
	*revel.Controller
}

/////// изменить на актуальный адрес и порт на котором висит webspcket //////////
var addr = flag.String("addr", "www.mysite.com:myport", "http service address")

/////// изменить на актуальный адрес и порт на котором висит webspcket  //////////

func (c App) Index() revel.Result {

	// Структуры для запросов-----------------------------------------------------------------------------------------------------------------------
	type Read struct {
		Title    string `json:"Title"`
		Body     string `json:"Body"`
		AuthorID int16  `json:"AuthorID"`
		NickName string `json:"NickName"`
	}
	type Write struct {
		Title string `json:"Title"`
		Id    int32  `json:"Id"`
	}

	var read Read
	var write Write

	// Запросы к API -------------------------------------------------------------------------------------------------------------------------------

	reqRead := `{"point":"read","title":"read API"}`
	ansRead := WebSocket(reqRead)
	json.Unmarshal([]byte(ansRead), &read)

	reqWrite := `{"point":"write","title":"write API"}`
	ansWrite := WebSocket(reqWrite)
	json.Unmarshal([]byte(ansWrite), &write)

	// Инициализация переменных отображения --------------------------------------------------------------------------------------------------------
	c.ViewArgs["TitleRead"] = read.Title
	c.ViewArgs["TitleWrite"] = write.Title

	// Логи ----------------------------------------------------------------------------------------------------------------------------------------
	Add2Log("ReadTxt", ":b:Title:-:", fmt.Sprint(read.Title))

	return c.Render()
}

func WebSocket(request string) string {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	websocket.DefaultDialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
		/////// изменить на актуальный адрес //////////
		ServerName: "www.mysite.com",
		/////// изменить на актуальный адрес //////////
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(request))
	if err != nil {
		log.Println("write:", err)
		return ""
	}

	done := make(chan struct{})
	var message []byte
	go func() {
		defer close(done)
		for {
			_, message, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			return
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return string(message)
		case <-ticker.C:
			log.Println("ticker.C")
			err := c.WriteMessage(websocket.TextMessage, []byte(`{"point":""}`))
			if err != nil {
				log.Println("write:", err)
				return ""
			}
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return ""
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return ""
		}
	}
}

func (c App) urlArg(arg int) string {
	return strings.Split(c.Request.URL.Path, "/")[arg]
}
