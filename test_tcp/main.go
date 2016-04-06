package main

import (
	"net"
	"os"
	"errors"
)
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func echo(conn net.Conn, log chan string) {
	defer conn.Close()
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n') // output message received
		fmt.Print("Message Received:", string(message))      // sample process for string received
		log <- message
		newmessage := strings.ToUpper(message)               // send new string back to client
		conn.Write([]byte(newmessage + "\n"))
		if message == "exit\n" {
			break
		}
	}
}

func collection(message_chan chan string, err_chan chan error) {
	looping := true
	for looping {
		select {
		case str := <-message_chan:
			fmt.Println("All: ", str)
		case err := <-err_chan:
			fmt.Println("error:", err)
			looping = false
		}
	}

}
func server() {
	fmt.Println("Launching server...")  // listen on all interfaces
	collec_chan := make(chan string)
	collec_error_chan := make(chan error)
	go collection(collec_chan, collec_error_chan)
	err := errors.New("exit")
	errfn := func() {collec_error_chan <- err}
	defer errfn()
	ln, _ := net.Listen("tcp", ":8081") // accept connection on port


	// run loop forever (or until ctrl-c)
	for {
		conn, _ := ln.Accept()
		go echo(conn, collec_chan)
	}
	
}

func client() {
	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text+"\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("s for server, c for client: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	if text == "s\n" {
		server()
	} else if text == "c\n" {
		client()
	} else {
		fmt.Println(text)
	}
}
