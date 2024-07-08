# net-cat

# Information
- This project recreates the NetCat utility in a Server-Client architecture. It can run in server mode on a specified port, listening for incoming connections, or in client mode, connecting to a specified port and transmitting information to the server.

- NetCat (nc system command) is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, and is able to open TCP connections, send UDP packets, listen on arbitrary TCP and UDP ports, and more.

# Features
- TCP Connection: Establishes a TCP connection between the server and multiple clients (one-to-many relationship).
- Username Requirement: Clients must choose a unique and valid username to join the chat.
- Connection Control: Limits the number of concurrent client connections to 10.
- Message Broadcasting: Clients can send messages to the chat, and all connected clients will receive them.
- Message Format: Messages are tagged with a timestamp and the sender's username in the format: [YYYY-MM-DD HH:MM:SS][client.name]: [client.message].
- No Empty Messages: The server does not broadcast empty messages from clients.
- Name Change: Clients can change their username without having to leave the server.
- Chat History: New clients receive the full chat history upon joining.
- History Log: A log file is started every time the server runs and is cleared when the server is closed.
- Join/Leave Notifications: All clients are notified when a new client joins or leaves the chat.
- Persistent Connections: Clients remain connected even if another client leaves the chat.
- Default Port: The server defaults to port 8989 if no port is specified. If a port is specified incorrectly, the server responds with a usage message: [USAGE]: ./TCPChat $port.

# Usage Instructions

1- Install Go: Ensure Go is installed on your device.

2- Navigate to the Directory: Use cd to go into the project directory.

3- Start the Server:

To start the server on the default port: go run main.go
To start the server on a specified port: go run main.go [port]

4- Connect to the Server:
Use NetCat to connect: nc [IP] [port]

5- Change Username:

Type /name followed by entering a valid username.