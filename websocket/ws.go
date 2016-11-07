// ws.go
package main

import (
	"fmt"
	//	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"golang.org/x/net/websocket"
)

func echoHandler(ws *websocket.Conn) {

	for {
		//		receivedtext := make([]byte, 100)

		//		n, err := ws.Read(receivedtext)

		//		if err != nil {
		//			fmt.Printf("Err Received: %d bytes\n", n)
		//		}

		//		s := string(receivedtext[:n])
		//		fmt.Printf("Received: %d bytes: %s\n", n, s)
		//		io.Copy(ws, ws)

		var in []byte
		if err := websocket.Message.Receive(ws, &in); err != nil {
			return
		}
		fmt.Printf("Received: %s\n", string(in))
		websocket.Message.Send(ws, in)
		//fmt.Printf("Sent: %s\n",s)
	}
}

func main() {
	s := websocket.Server{Handler: echoHandler}
	http.Handle("/echo", s)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
