package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Command struct {
	Action string `json:"action"`
	Room   string `json:"room"`
	Id     string `json:"id"`
	Data   string `json:"data"`
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Server is not ready")
		os.Exit(0)
	}
	defer conn.Close()
	go run(conn)

	for {
		var com Command = Command{}
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Action: ")
		com.Action = readLine(reader)
		fmt.Print("Room: ")
		com.Room = readLine(reader)
		fmt.Print("Id: ")
		com.Id = readLine(reader)
		fmt.Print("Data: ")
		com.Data = readLine(reader)
		data, _ := json.Marshal(com)
		fmt.Println(string(data))
		fmt.Fprint(conn, string(data)+"\n")
	}
}

func readLine(reader *bufio.Reader) string {
	data, _ := reader.ReadString('\n')
	return strings.TrimSuffix(data, "\n")
}

func run(conn net.Conn) {
	for {
		// read in input from stdin
		reader := bufio.NewReader(conn)
		text, err := reader.ReadString('\n')
		// send to socket
		if err != nil {
			fmt.Println("Server is close")
			conn.Close()
			os.Exit(0)
		}
		fmt.Println("Message from server: " + text)
	}
}
