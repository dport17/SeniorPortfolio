/* Simple EchoServer in GoLang by Phu Phung, customized by <YOUR NAME> for SecAD*/
package main

import (
	"fmt"
	"net"
	"os"
)

var allClients_conns = make(map[net.Conn]string)
var lostClient = make(chan net.Conn)

const BUFFERSIZE int = 1024

func main() {
	newclient := make(chan net.Conn)

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(0)
	}
	port := os.Args[1]
	if len(port) > 5 {
		fmt.Println("Invalid port value. Try again!")
		os.Exit(1)
	}
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Cannot listen on port '" + port + "'!\n")
		os.Exit(2)
	}
	fmt.Println("EchoServer in GoLang developed by Phu Phung, SecAD, revised by Devin Porter")
	fmt.Printf("EchoServer is listening on port '%s' ...\n", port)
	go func() {
		for {
			client_conn, _ := server.Accept()
			newclient <- client_conn
		}
	}()
	for {
		select {
		case client_conn := <-newclient:
			allClients_conns[client_conn] = client_conn.RemoteAddr().String()
			go client_goroutine(client_conn)
		case client_conn := <-lostClient:
			delete(allClients_conns, client_conn)
			byemessage := fmt.Sprintf("Client %s is DISCONNECTED\n# of clients is now %d\n", client_conn.RemoteAddr().String(), len(allClients_conns))
			fmt.Println(byemessage)
			go sendToAll([]byte(byemessage))
		}
	}

}

func client_goroutine(client_conn net.Conn) {
	welcomemessage := fmt.Sprintf("A new client '%s' connected!\n# of connected clients: %d\n", client_conn.RemoteAddr().String(), len(allClients_conns))
	fmt.Println(welcomemessage)
	go sendToAll([]byte(welcomemessage))

	var buffer [BUFFERSIZE]byte
	go func() {
		for {
			byte_received, read_err := client_conn.Read(buffer[0:])
			if read_err != nil {
				fmt.Println("Error in receiving...")
				lostClient <- client_conn
				return
			}

			clientdata := buffer[0 : byte_received-2]
			fmt.Printf("Received data: %s, len=%d\n", clientdata, len(clientdata))
			result1 := string(clientdata)

			if len(clientdata) > 5 && string(clientdata[0:5]) == "login" {
				fmt.Println("Potential attack detected: ignoring input.")
			} else if result1 == "login" {
				loginmessage := "login data\n"
				client_conn.Write([]byte(loginmessage))
			} else {
				loginmessage := "non-login data\n"
				client_conn.Write([]byte(loginmessage))
			}
			go sendToAll(buffer[0:byte_received])
		}
	}()
}

func sendToAll(data []byte) {
	for client_conn, _ := range allClients_conns {
		_, write_err := client_conn.Write(data)
		if write_err != nil {
			fmt.Println("Error in sending...")
			continue
		}
	}
	fmt.Printf("Send data: %s Sent to all connected clients!\n", data)
}
