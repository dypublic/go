// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// //+build ignore

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"fmt"
	"os"
	"errors"
	"time"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type FileHead struct {
	File_name string `json:"name"`
	File_size int64  `json:"size"`
}
func get_head(c *websocket.Conn) (FileHead, error){
	head := FileHead{}
	err := c.ReadJSON(&head)
	//_, message, err := c.ReadMessage()
	//fmt.Println(message)
	if err != nil {
		fmt.Println("head not correct!")
		return head, err
	}
	if head.File_name == "" {
		return head, errors.New("head is empty!")
	}
	fmt.Println(head)
	return head, nil
}

func recv_done(c *websocket.Conn, file *os.File) error {

	fmt.Println("recv all bytes")
	file.Sync()
	err := c.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Send close fail:", err)
		return err
	}
	return nil
}

func check_recv_done(save_size int64, head *FileHead,
						c *websocket.Conn, file *os.File) (bool, error) {
	if save_size == head.File_size {
		return true, recv_done(c, file)
	} else if save_size > head.File_size {
		err_str := fmt.Sprintln("! recv too much, org: ", head.File_size, " now:", save_size)
		return false, errors.New(err_str)
	}
	return false, nil
}
func recv_piece(c *websocket.Conn) ([]byte, error) {
	c.SetReadDeadline(time.After(time.Second))
	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("Recv:", err)
		return []byte{}, err
	}
	if mt != websocket.BinaryMessage {
		log.Println("Not a BinaryMessage, ", mt)
		return []byte{}, errors.New("Not a BinaryMessage")
	}
	return message, nil
}

func check_len(len int, c *websocket.Conn) error{
	if len > 0 {
		return nil
	}
	log.Printf("recv: %d", len)
	err := c.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(
			websocket.CloseUnsupportedData,
			"Recv zero byte"))
	if err != nil {
		log.Println("write close:", err)
		return err
	}
	return errors.New("recv wrong length message")

}

func save_file(c *websocket.Conn, head FileHead) error {
	file, ferr := os.OpenFile(head.File_name,
							os.O_CREATE | os.O_WRONLY,
							os.ModePerm)
	if ferr != nil {
		fmt.Println("create file fail:", ferr)
		return ferr
	}
	defer file.Close()

	var save_size int64 = 0
	for {
		message, err := recv_piece(c)
		if err != nil {
			return err
		}
		if err := check_len(len(message), c); err != nil{
			return err
		}
		file.Write(message)

		save_size += int64(len(message))
		done, err := check_recv_done(save_size, &head, c, file)
		if err != nil{
			return err
		}
		if done{
			break
		}
	}
	return nil
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in upload")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	head, err := get_head(c)
	if err != nil{
		log.Println("head error, ", err.Error())
		return
	}
	//fmt.Println(head)
	save_file(c, head)
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
