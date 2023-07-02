package main

import (
	"bytes"
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
		panic(err)
	}

	// Bind the socket to an IPv4 address and a port through which to listen to
	syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: serverAddress,
	})

	defer syscall.Close(fd)
	for {
		// Instantiate a 1024 byte slice and copy incoming data into the slice
		b := make([]byte, 1024)
		n, from, err := syscall.Recvfrom(fd, b, 0)
		if err != nil {
			return
		}
		if from != nil {
			// Uppercase all characters and send the resulting []byte to the client
			err = syscall.Sendto(fd, bytes.ToUpper(b[:n]), 0, from)
			if err != nil {
				return
			}
		}
	}
}
