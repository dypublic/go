package main

import(
	"os"
	"flag"
	"net/url"
	"log"
	//"errors"
	"os/signal"
	"bufio"
	"github.com/gorilla/websocket"
	//"io"
	"sync"
	"fmt"
	//"encoding/json"
	"path"
	"time"
)

type HeadInfo struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
}

func main() {
	addr := flag.String("addr", "192.168.1.108:8080", "ip:port")
	flag.Parse()
	filename := flag.Arg(0)
	//filename := os.Args[1]
	u := url.URL{Scheme:"ws", Host: *addr, Path: "/upload"}
	log.Printf("target addr: %s", u.String())

	//file.
	info, err := os.Stat(filename)
	if err != nil{
		log.Println("file info not get, ", err.Error())
		os.Exit(1)
	}
	name := path.Base(filename)
	headinfo := &HeadInfo{
					Name: name,
					Size: info.Size(),
	}

	file, err := os.Open(filename)
	if err != nil{
		log.Panic("file not open: ", filename, err.Error())
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	ws_conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws_conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ws_conn.WriteJSON(headinfo)
		fmt.Println(headinfo)
		//buffer := make([]byte, 1024*8)
		var buffer [1024*8]byte

SendLoop:
		for {
			piece := buffer[:]
			len, err := reader.Read(piece)
			if err != nil {
				log.Println("read error, ", err.Error())
				break
			}
			piece = buffer[:len]
			ws_err := ws_conn.WriteMessage(websocket.BinaryMessage, piece)
			if ws_err != nil{
				log.Println("ws write error:", ws_err.Error())
				break
			}
			select {
			case <- interrupt:
				log.Print("user interrupt")
				ws_err := ws_conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseGoingAway,"user cancel"),
					time.Now().Add(100*time.Millisecond))
				if ws_err != nil{
					log.Println("close message send fail", ws_err.Error())
				}
				break SendLoop
			default:
				//log.Println("no interrupt")
				//time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	wg.Wait()

}
