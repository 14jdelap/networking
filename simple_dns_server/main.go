package main

import (
	"math/rand"
	"strconv"
	"strings"
)

const (
	PROTOCOL = 0
	// SERVER_PORT = 8080
	// MTU         = 1500
)

/*
UDP socket:
1. Create socket
2. Bind fd to a port
3. Create packet payload
4. Send packet to DNS server
5. Receive back the DNS response
*/

type typecode [2]byte
type class [2]byte

type Message struct {
	header     Header
	question   Question
	answer     Answer
	authority  Authority
	additional Additional
}

type Header struct {
	id      [2]byte
	options [2]byte
	qdcount [2]byte
	ancount [2]byte
	nscount [2]byte
	arcount [2]byte
}

type Question struct {
	name []byte // up to 256 bits
	typecode
	class
}

type Answer struct {
	resources []RelatedResource
}

type Authority struct {
	resources []RelatedResource
}

type Additional struct {
	resources []RelatedResource
}

type RelatedResource struct {
	name     []byte // up to 256 bits
	ttl      [2]byte
	rdlength [2]byte
	rdata    [2]byte
}

func (m *Message) constructDNSQuery() []byte {
	payload := []byte{}
	for _, field := range m {
		for _, b := range field {

		}
	}
	return payload
}

func main() {
	// args := os.Args[1:]
	// if len(args) != 1 {
	// 	fmt.Println("please pass only one name to query")
	// } else {
	// 	fmt.Println(args[0])
	// }

	// fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, PROTOCOL)
	// if err != nil {
	// 	log.Fatalf("creating a socket: %s\n", err)
	// }

	// err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 8080})
	// if err != nil {
	// 	log.Fatalf("binding socket %d: %s\n", fd, err)
	// }

	// syscall.Sendto(fd, p)

	h := Header{
		id:      generateID(),
		options: [2]byte{0x00, 0x00},
		qdcount: [2]byte{0x00, 0x01},
		ancount: [2]byte{0x00, 0x00},
		nscount: [2]byte{0x00, 0x00},
		arcount: [2]byte{0x00, 0x00},
	}
	q := Question{
		name:     []byte{},
		typecode: [2]byte{0x00, 0x01}, // 1 is the value for A Records
		class:    [2]byte{0x00, 0x01}, // 1 should be the value for internet
	}

	name := "google.com" // Protocol (e.g. http) and hostname (e.g. www)
	splitString := strings.Split(name, ".")
	for _, str := range splitString {
		// Create length byte and succeed it with the label in bytes
		b := make([]byte, 0, len(str)+1)
		strLength := strconv.Itoa(len(str))
		b = append(b, []byte(strLength+str)...)
		// Append each label to q.name
		q.name = append(q.name, b...)
	}
	// Append the zero byte to indicate end of name
	q.name = append(q.name, 0x00)

	m := Message{
		header:     h,
		question:   q,
		answer:     Answer{},
		authority:  Authority{},
		additional: Additional{},
	}
}

// generateID converts a uint16 number to a byte slice.
func generateID() [2]byte {
	num := uint16(rand.Intn(1 << 16))
	return [2]byte{byte(num >> 8), byte(num & 0xFF)}
}

// DOES IT MAKE SENSE TO HAVE SO MANY STRUCTS IF I CAN'T THEN CREATE THE BYTE SLICE EASILY AT THE END?
