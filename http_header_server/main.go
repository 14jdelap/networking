package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
)

const (
	PROTOCOL    = 0
	SERVER_PORT = 8080
)

var (
	serverAddress = [4]byte{127, 0, 0, 1}
	sigint        = []byte{255, 244, 255, 253, 6}
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

	err = syscall.Listen(fd, 10)
	if err != nil {
		log.Fatalf("listening for incoming connections: %s\n", err)
	} else {
		fmt.Printf("listening on port %d\n", SERVER_PORT)
	}

	defer syscall.Close(fd)
	for {
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
		b := make([]byte, 1024)
		n, err := syscall.Read(fd, b)
		if err != nil {
			log.Fatalf("receiving data from client: %s\n", err)
		}

		// Return and close connection if data is the string "exit" or SIGINT
		if bytes.Equal(bytes.ToLower(b[:n-2]), []byte("exit")) || bytes.Equal((b[:n]), sigint) {
			fmt.Println("closing connection")
			return
		}
		// Uppercase all characters and send the resulting []byte to the client
		_, err = syscall.Write(fd, bytes.ToUpper(b[:n]))
		if err != nil {
			log.Fatalf("sending data to client: %s\n", err)
		}
	}
}
