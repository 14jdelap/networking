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
	address = [4]byte{127, 0, 0, 1}
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, PROTOCOL)
	if err != nil {
		panic(err)
	}

	syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: SERVER_PORT,
		Addr: address,
	})

	defer syscall.Close(fd)
	for {
		b := make([]byte, 1024)
		n, from, err := syscall.Recvfrom(fd, b, 0)
		if err != nil {
			return
		}
		if from != nil {
			err = syscall.Sendto(fd, bytes.ToUpper(b[:n]), 0, from)
			if err != nil {
				return
			}
		}
	}
}
