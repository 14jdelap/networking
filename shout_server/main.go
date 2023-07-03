package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"syscall"
)

const (
	PROTOCOL    = 0
	SERVER_PORT = 8080
)

var (
	serverAddress = [4]byte{127, 0, 0, 1}
)

func main() {
	// Create the server's socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, PROTOCOL)
	if err != nil {
		log.Fatalf("creating socket: %s", err)
	}

	// Bind the socket to an IPv4 address and a port through which to listen to
	err = syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: serverAddress,
	})
	if err != nil {
		log.Fatalf("binding socket: %s\n", err)
	} else {
		fmt.Printf("receiving incoming packets in %s:%d\n", net.IP(serverAddress[:]).String(), SERVER_PORT)
	}

	defer syscall.Close(fd)
	for {
		// Instantiate a 1024 byte slice and copy incoming data into the slice
		b := make([]byte, 1024)
		n, from, err := syscall.Recvfrom(fd, b, 0)
		if err != nil {
			log.Fatalf("reading data from client: %s\n", err)
		}
		if from != nil {
			// Uppercase all characters and send the resulting []byte to the client
			err = syscall.Sendto(fd, bytes.ToUpper(b[:n]), 0, from)
			if err != nil {
				log.Fatalf("sending data to client: %s\n", err)
			}
			fmt.Printf("received %b, sent %b", b[:n], bytes.ToUpper(b[:n]))
		}
	}
}
