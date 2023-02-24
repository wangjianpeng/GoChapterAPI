package learngorilla

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func DoWebHandler(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		if string(message) == "Hello WebSockets!" {
			message = []byte("pong")
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func DoRunWebSocketServer() {
	bindAddress := "192.168.12.57:9998"
	r := gin.Default()
	r.GET("/ping", DoWebHandler)
	r.Run(bindAddress)
}
