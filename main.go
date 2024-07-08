package main

import (
	"fmt"
	"log"
	"net"
	nc "netcat/funcs"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Conn net.Conn
	Name string
}

// variable to store chat history and mutex for shared resources
var chatHistory []string
var mu sync.Mutex

// map of 10 clients so we can keep track of clients and stop connections if server is full
var clients = make(map[string]Client)

// func init that runs before the main function to reset the log file everytime the server is run
func init() {
	// Open the file in write mode and truncate it to 0 bytes
	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		log.Fatal(err)
	}
	defer file.Close()
}

func main() {

	//change port based on number of arguments - default is 8989
	if len(os.Args) > 2 {
		fmt.Println("Usage: go run main.go [port]")
		return
	}

	port := "8989"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	portInt, errInt := nc.Atoi(port)
	if errInt {
		fmt.Println("Error:", errInt)
		return
	}

	if portInt < 8000 || portInt > 9999 {
		fmt.Println("Error: port must be between 8000 and 9999")
		return
	}
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port", port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// Deferred close of connection when function exits in case of an error
	defer conn.Close()

	// Check length of clients map and exit (close client's connection) if more than 10
	mu.Lock()
	if len(clients) >= 10 {
		mu.Unlock()
		_, err := conn.Write([]byte("[SERVER IS FULL, PLEASE TRY AGAIN LATER]\n"))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		return
	}
	mu.Unlock()

	// Send welcome message to client
	data := []byte(nc.PrintWelcome())
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create buffer to read data from client
	buffer := make([]byte, 1024)

	var clientName string // to store client's name

	// Prompt client for name
	_, err = conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for {

		// Read client's name
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		clientName = strings.TrimSpace(string(buffer[:n]))

		// Check if the name is valid, if not reprompt otherwise continue
		if !nc.IsValidName(clientName) {
			_, err := conn.Write([]byte("[INVALID NAME, PLEASE ENTER A VALID NAME]: "))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			continue
		}

		// Check if the name is taken
		mu.Lock()
		nameTaken := isNameTaken(clientName)
		mu.Unlock()

		// If name is taken, reprompt
		if nameTaken {
			_, err := conn.Write([]byte("[NAME IS USED, PLEASE ENTER ANOTHER NAME]: "))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			continue
		}

		// If name is valid and not taken, break the loop
		break
	}

	// Create a new Client struct and store in map
	client := Client{
		Conn: conn,
		Name: clientName,
	}
	mu.Lock()
	// Add client to map
	clients[clientName] = client
	mu.Unlock()

	// Send chat history to new client
	mu.Lock()
	for _, message := range chatHistory {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error:", err)
			mu.Unlock()
			return
		}
	}
	mu.Unlock()

	// Inform other clients that a new client has joine
	joinedMessage := fmt.Sprintf("%s Has Joined The Chat", clientName)
	broadcastMessage("", joinedMessage)

	// Handle client messages
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Client", clientName, "has disconnected.")
			// Remove client from map
			mu.Lock()
			delete(clients, clientName)
			mu.Unlock()
			// Inform other clients that this client has left
			leftMessage := fmt.Sprintf("%s Has Left The Chat", clientName)
			broadcastMessage("", leftMessage)
			return
		}

		// Check if client wants to change name
		if strings.TrimSpace(string(buffer[:n])) == "/name" {
			//call changeName function and store the new name
			newName := changeName(clientName, conn, buffer)
			if newName != "" {
				clientName = newName // Update the local variable to the new name
			}
			continue
		}

		// Limit message to 1024 bytes
		if n > 1024 {
			_, err := conn.Write([]byte("[YOU HAVE EXCEEDED THE MESSAGE LIMIT]"))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			continue
		}

		// Trim white spaces and broadcast message only if message is not empty
		if trimmed := strings.TrimSpace(string(buffer[:n])); trimmed != "" {
			broadcastMessage(clientName, trimmed)
		}
	}
}

func broadcastMessage(sender string, message string) {
	// Get time stamp for the message
	time := time.Now()
	timeNow := time.Format("[2006-01-02 15:04:05]")

	// variable to store formatted message based on sender
	var formattedMessage string

	// loop over clients in the map and send message to all clients except the sender
	mu.Lock()
	for name, client := range clients {
		// if sender is not empty send to all clients except the sender
		if name != sender && sender != "" {
			// Format message with time stamp and sender name
			formattedMessage = fmt.Sprintf("%s [%s]: %s\n", timeNow, sender, message)
			_, err := client.Conn.Write([]byte(formattedMessage))
			if err != nil {
				fmt.Printf("Error broadcasting message to client %s: %s\n", name, err)
			}
			// otherwise if sender is empty meaning its a joined/left message we send to all clients with different formatting
		} else if sender == "" {
			formattedMessage = fmt.Sprintf("%s : %s\n", timeNow, message)
			_, err := client.Conn.Write([]byte(formattedMessage))
			if err != nil {
				fmt.Printf("Error broadcasting message to client %s: %s\n", name, err)
			}
		}
	}

	// adding any messages to chat history so we can show them to new clients
	chatHistory = append(chatHistory, formattedMessage)
	// call logging function to log messages in file
	nc.Logging(formattedMessage)

	mu.Unlock()
}

// function to check if user name is taken or not
func isNameTaken(name string) bool {
	// check if name (key) exists in the map
	_, exists := clients[name]
	// return true if name exists in the map
	return exists
}

func changeName(oldName string, conn net.Conn, buffer []byte) string {
	// Prompt client to enter their new name
	_, err := conn.Write([]byte("[ENTER YOUR NEW NAME]: "))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	//infinite loop until name is valid and not taken
	for {

		// Read name from client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}

		newName := strings.TrimSpace(string(buffer[:n]))

		// Check if the name is valid
		if !nc.IsValidName(newName) {
			_, err := conn.Write([]byte("[INVALID NAME, PLEASE ENTER A VALID NAME]: "))
			if err != nil {
				fmt.Println("Error:", err)
				return ""
			}
			continue
		}

		// Check if the name is taken
		mu.Lock()
		nameTaken := isNameTaken(newName)
		mu.Unlock()

		if nameTaken {
			_, err := conn.Write([]byte("[NAME IS USED, PLEASE ENTER ANOTHER NAME]: "))
			if err != nil {
				fmt.Println("Error:", err)
				return ""
			}
			continue
		}

		// Name is valid and not taken, update client name in map
		mu.Lock()
		//delete old name from map
		client := clients[oldName]
		delete(clients, oldName)

		//add new name to map
		client.Name = newName
		clients[newName] = client
		mu.Unlock()

		// Inform other clients that the name has changed
		broadcastMessage("", fmt.Sprintf("%s changed their name to %s", oldName, newName))

		// Return the new name
		return newName
	}
}
