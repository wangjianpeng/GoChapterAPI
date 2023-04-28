package socketbot

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID         string
	ClientType string
	Conn       *websocket.Conn
}

var clients sync.Map

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func MutipleServe() {
	router := gin.Default()
	router.GET("/web", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "./socketbot/logclient_web.html")
	})

	router.GET("/unity", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "./socketbot/logclient_unity.html")
	})

	router.GET("/ws", func(ctx *gin.Context) {
		conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Connect msg ", string(message))
		client := &Client{
			ID:   "client_" + conn.RemoteAddr().String(),
			Conn: conn,
		}
		// client_type:unity
		// client_type:web
		if string(message[0:12]) == "client_type:" {
			client.ClientType = string(message[12:])
			fmt.Println("Connect mss type ", client.ClientType)
		}

		clients.Store(client.ID, client)
		go func() {
			defer func() {
				client.Conn.Close()
				clients.Delete(client)
				log.Printf("Client %s disconnected\n", client.ID)
			}()

			for {
				_, message, err := client.Conn.ReadMessage()
				if err != nil {
					log.Println(err)
					break
				}

				log.Printf("Received message %s from type  %s , client ID %s\n", message, client.ClientType, client.ID)
				//check unnique id
				if client.ClientType == "unity" {
					replyMsg := fmt.Sprintf("Received msg %s from unity %s", message, client.ID)
					clients.Range(func(key, value interface{}) bool {
						if key.(string) != client.ID {
							replyMsg = "Write to web <----- " + replyMsg
							err = value.(*Client).Conn.WriteMessage(websocket.TextMessage, []byte(replyMsg))
							if err != nil && !websocket.IsCloseError(err, websocket.CloseGoingAway) {
								// panic(err)
								log.Println(err)
							}
						}
						return true
					})
				} else if client.ClientType == "web" {
					log.Printf("Received msg from web %s", client.ID)
				} else {
					log.Printf("Received invalid msg from other %s", client.ID)
				}
			}
		}()
	})

	err := router.Run(":14660")
	if err != nil {
		log.Println(err)
	}
}
