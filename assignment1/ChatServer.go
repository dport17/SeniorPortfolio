/* Simple EchoServer in GoLang by Phu Phung, customized by <YOUR NAME> for SecAD*/
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

var allClients_conns = make(map[net.Conn]string)
var lostClient = make(chan net.Conn)
var allAuthClients = make(map[net.Conn]string)

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
			delete(allAuthClients, client_conn)
			byemessage := fmt.Sprintf("Client %s is DISCONNECTED\n# of clients is now %d\n", client_conn.RemoteAddr().String(), len(allClients_conns))
			fmt.Println(byemessage)
			go sendToAll([]byte(byemessage))
		}
	}

}

func client_goroutine(client_conn net.Conn) {
	welcomemessage := fmt.Sprintf("A new client '%s' connected!\n# of connected clients: %d\n", client_conn.RemoteAddr().String(), len(allClients_conns))
	fmt.Println(welcomemessage)
	go sendTo(client_conn, []byte(welcomemessage))

	var buffer [BUFFERSIZE]byte
	go func() {
		for {
			byte_received, read_err := client_conn.Read(buffer[0:])
			if read_err != nil {
				fmt.Println("Error in receiving...")
				lostClient <- client_conn
				return
			}

			clientdata := buffer[0:byte_received]
			fmt.Printf("Received data: %s, len=%d\n", clientdata, len(clientdata))
			//compare the data
			result1 := string(clientdata)
			fmt.Printf("The clientdata as a string is:'%s'\n", result1)

			if len(clientdata) >= 5 && result1[0:5] == "login" { //handle login attempts.

				fmt.Printf("The JSON received is: %s\n", result1[7:len(clientdata)])

				//Make JSON string parseable.
				var result map[string]interface{}
				json.Unmarshal([]byte(result1[7:]), &result)
				checkLogin(client_conn, result)
			} else { //handle all non-login, non-user list request messages.

				if _, present := allAuthClients[client_conn]; present {
					fmt.Println("Sending user is authenticated.")
					sendingUser := allAuthClients[client_conn]
					//Make JSON string parseable.
					var result map[string]interface{}
					json.Unmarshal([]byte(result1), &result)
					mode := fmt.Sprintf("%v", result["mode"])
					user := fmt.Sprintf("%v", result["user"])
					msg := fmt.Sprintf("%v", result["msg"])
					if msg == "users" {
						sendTo(client_conn, []byte("Here are the users: "))
						go sendUserData(client_conn)
					} else if mode == "pub" {
						msg = "Message from " + sendingUser + ": " + msg
						go sendPublicMessage(msg)
					} else {
						msg = "PM from " + sendingUser + ": " + msg
						success := sendPrivMessage(result["user"], msg)
						if !success {
							go sendTo(client_conn, []byte("Cannot send message to offline user "+user+": "+msg))
						}
					}
				} else {
					sendTo(client_conn, []byte("You cannot send a message until you are authenticated."))
					go sendTo(client_conn, []byte("LF"))
				}
			}
		}
	}()
}

func sendToAll(data []byte) {
	for client_conn, _ := range allAuthClients {
		_, write_err := client_conn.Write(data)
		if write_err != nil {
			fmt.Println("Error in sending...")
			continue
		}
	}
	fmt.Printf("Send data: %s Sent to all connected clients!\n", data)
}

func sendTo(client_conn net.Conn, data []byte) {
	client_conn.Write(data)
}

func sendUserData(client_conn net.Conn) {
	var usersString = "\nUsers:\n"
	for key := range allAuthClients {
		var element = allAuthClients[key]
		usersString = usersString + element + "\n"
	}
	sendTo(client_conn, []byte(usersString))
}

func findUserSocket(username interface{}) net.Conn {
	for key := range allAuthClients {
		if allAuthClients[key] == username {
			return key
		}
	}
	return nil
}

func sendPublicMessage(msg string) {
	go sendToAll([]byte(msg))
}

func sendPrivMessage(user interface{}, msg string) bool {
	client_conn := findUserSocket(user)
	if client_conn != nil {
		go sendTo(client_conn, []byte(msg))
		return true
	}
	return false
}

func checkLogin(client_conn net.Conn, result map[string]interface{}) {
	username := fmt.Sprintf("%v", result["username"])
	go func() {
		if (result["username"] == "devin" && result["password"] == "porter") || (result["username"] == "chloe" && result["password"] == "richie") {
			allAuthClients[client_conn] = username
			go sendTo(client_conn, []byte("\nWelcome to the chatserver!\n"))
			sendToAll([]byte("\n"+username + " has joined the ChatServer!"))
			for key := range allAuthClients {
				sendUserData(key)
			}

		} else {
			go sendTo(client_conn, []byte("LF"))
		}
	}()
}
