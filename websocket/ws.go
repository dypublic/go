// ws.go
package main

import (
	"golang.org/x/net/websocket"
	//    "io"
	"fmt"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {

	for {
		receivedtext := make([]byte, 100)

		n, err := ws.Read(receivedtext)

		if err != nil {
			fmt.Printf("Received: %d bytes\n", n)
		}

		s := string(receivedtext[:n])
		fmt.Printf("Received: %d bytes: %s\n", n, s)
		//io.Copy(ws, ws)
		//fmt.Printf("Sent: %s\n",s)
	}
}

func main() {
	s := websocket.Server{Handler: echoHandler}
	http.Handle("/echo", s)
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
