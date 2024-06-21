package socknet

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
)

const (
	PROTOCOL    = 0
	SERVER_PORT = 8080
	MTU         = 1500
)

var (
	serverAddress = [4]byte{127, 0, 0, 1}
	sigint        = []byte{255, 244, 255, 253, 6}
)

/*
	Objective: create an HTTP 1.1 connection (i.e. with TCP) that returns
	a JSON with the HTTP headers

	Requirements:
	1. Implement http.ListenAndServe thorugh TCP sockets
	2. Implement http.HandleFunc to read the HTTP headers and return them as JSON
*/

type Socknet struct {
	fd int
}

func NewSocket() (*Socknet, error) {
	// Create the server's socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, PROTOCOL)
	if err != nil {
		return nil, fmt.Errorf("creating a socket: %w", err)
	}
	return &Socknet{fd}, nil
}

func (s *Socknet) Connect(fd int, sockaddr *syscall.SockaddrInet4) error {
	// Bind the socket to an IPv4 address and a port through which to listen to
	err := syscall.Bind(fd, sockaddr)
	if err != nil {
		return fmt.Errorf("binding a socket: %w", err)
	}

	err = syscall.Listen(fd, 10)
	if err != nil {
		return fmt.Errorf("listening for incoming connections: %w", err)
	} else {
		fmt.Printf("listening on port %d\n", SERVER_PORT)
	}

	defer syscall.Close(fd)
	for {
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			return fmt.Errorf("accepting connection: %s", err)
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
		n, err := syscall.Read(fd, b)
		if err != nil {
			log.Fatalf("receiving data from client: %s\n", err)
		}

		// Return and close connection if data is the string "exit" or SIGINT
		// SIGNINT is []byte{} in nc and []byte{255, 244, 255, 253, 6} in telnet
		if bytes.Equal((b[:n]), sigint) || n == 0 || bytes.Equal(bytes.ToLower(b[:n-2]), []byte("exit")) {
			fmt.Println("closing connection")
			return
		}
		// Uppercase all characters and send the resulting []byte to the client
		_, err = syscall.Write(fd, b[:n])
		if err != nil {
			log.Fatalf("sending data to client: %s\n", err)
		}
	}
}
