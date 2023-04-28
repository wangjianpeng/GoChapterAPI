package socketbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HelloTest() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		bodydata, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(string(bodydata))
		//base64 decode
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":14600")
}

func InputWebPage() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./socketbot/inputbtn.html")
	})

	router.POST("/send", func(c *gin.Context) {
		var data struct {
			Text string `json:"text"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		sendMsg(data.Text)
		c.String(http.StatusOK, "Text sent: %s", data.Text)
	})

	router.Run(":14600")
}

func sendMsg(text string) {
	fmt.Println("Sending message:", text)
}

func LogClient() {

}

func LogServer() {
	http.HandleFunc("/ws", HandleLogWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{}

func HandleLogWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		msg := fmt.Sprintf("Log message a %v", time.Now().Unix())
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// support server and client connect by html.
func LogHtmlWebSocket() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "./socketbot/logclient.html")
	})

	r.GET("/ws", func(ctx *gin.Context) {
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			panic(err)
		}

		go func() {
			defer ws.Close()
			id := getNextID()
			fmt.Printf("new client connnected with iD: %d\n", id)
			for {
				messagetype, p, err := ws.ReadMessage()
				if err != nil && !websocket.IsCloseError(err, websocket.CloseGoingAway) {
					// panic(err)
					log.Println(err)
					break
				}
				if messagetype == websocket.TextMessage {
					message := string(p)
					println(message)
				}

				replyMsg := "receive message: " + string(p)
				println(replyMsg)
				err = ws.WriteMessage(websocket.TextMessage, []byte(replyMsg))
				if err != nil && !websocket.IsCloseError(err, websocket.CloseGoingAway) {
					// panic(err)
					log.Println(err)
					break
				}
			}
		}()

	})
	r.Run(":14660")
}

var mu sync.Mutex
var nextId = 1

func getNextID() int {
	mu.Lock()
	defer mu.Unlock()
	id := nextId
	nextId++
	return id
}
