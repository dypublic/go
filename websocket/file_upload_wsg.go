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
	"io"
	"sync"
)

type HeadInfo struct {
	name string `json:"name"`
	size int64 `json:"size"`
}

func main() {
	addr := flag.String("addr", "192.168.1.108:8080", "ip:port")
	flag.Parse()
	filename := flag.Arg(0)
	//filename := os.Args[1]
	u := url.URL{Scheme:"ws", Host:*addr, Path:"/upload"}
	log.Printf("target addr: %s", u.String())

	//file.
	info, err := os.Stat(filename)
	if err != nil{
		log.Println("file info not get, ", err.Error())
		os.Exit(1)
	}
	headinfo := HeadInfo{}
	headinfo.name = filename
	headinfo.size = info.Size()

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

		ws_conn.WriteJSON(headinfo)
		var buffer [1024*8]byte
		file_end := false
		for !file_end {
			len, err := reader.Read(buffer[0:0])
			if err == io.EOF{
				//file_end = true
				break
			}else if err != nil {
				log.Println("read file and get error, ", err.Error())
			}
			piece := buffer[:len]
			ws_err := ws_conn.WriteMessage(websocket.BinaryMessage, piece)
			if ws_err != nil{
				log.Println("ws write error:", ws_err.Error())
				break
			}


			select {
			case <- interrupt:
				log.Print("user interrupt")
				break
			default:
			}
		}
		wg.Done()

	}()

	wg.Wait()

}
