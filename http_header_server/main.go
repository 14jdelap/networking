package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"syscall"
)

const (
	PROTOCOL    = 0
	SERVER_PORT = 8080
	MTU         = 1500
	COLON       = 58
)

var (
	serverAddress     = [4]byte{127, 0, 0, 1}
	lineSeparator     = []byte{13, 10}
	httpBodySeparator = []byte{13, 10, 13, 10}
)

func main() {
	// Create the server's socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, PROTOCOL)
	if err != nil {
		log.Fatalf("creating a socket: %s\n", err)
	}

	// Bind the socket to an IPv4 address and a port through which to listen to
	syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: serverAddress,
	})

	// Listen for incoming connections
	err = syscall.Listen(fd, 10)
	if err != nil {
		log.Fatalf("listening for incoming connections: %s\n", err)
	} else {
		fmt.Printf("listening on port %d\n", SERVER_PORT)
	}

	defer syscall.Close(fd)
	for {
		// Accept incoming connection requests
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			log.Fatalf("accepting connection: %s\n", err)
		}
		fmt.Println("connection established")

		// Handle each accepted connection concurrently
		go handleConnection(nfd)
	}
}

func handleConnection(fd int) {
	// Close the connection when the function returns
	defer syscall.Close(fd)
	for {
		// Instantiate a 1024 byte slice and copy incoming data into the slice
		b := make([]byte, MTU)
		_, err := syscall.Read(fd, b)
		if err != nil {
			log.Fatalf("receiving data from client: %s\n", err)
		}

		// Split by divisor between header and body
		splitRequest := bytes.Split(b, httpBodySeparator)

		// If no request, return a string asking for a request so the user knows he's used the server wrong.
		// This if defines (wrongly) an http request as a string with a single instance of \r\n\r\n
		if len(splitRequest) != 2 {
			_, err = syscall.Write(fd, []byte("http request not sent: please send one to receive your http response as json\r\n"))
			if err != nil {
				log.Fatalf("sending data to client: %s\n", err)
			}
			// Return and close connection
			return
		}

		// Split header by lines
		splitHeader := bytes.Split(splitRequest[0], lineSeparator)
		headerValuesByName := map[string]string{}
		for _, line := range splitHeader[1:] {
			// Split header lines by colon to separate keys and values
			splitLine := bytes.Split(line, []byte{COLON})
			headerValuesByName[string(splitLine[0])] = strings.TrimLeft(string(splitLine[1]), " ")
		}
		jsonData, err := json.MarshalIndent(headerValuesByName, "", "  ")
		if err != nil {
			fmt.Printf("marshaling http header: %s", err)
		}

		// Construct HTTP response with JSON payload and write to client
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: application/json\r\n\r\n%s", len(jsonData), string(jsonData))
		_, err = syscall.Write(fd, []byte(response))
		if err != nil {
			log.Fatalf("sending data to client: %s\n", err)
		}
	}
}
