# Shout Server

The socket interface is one of the most successful abstractions of all time. It's the standard, almost universal way to get access to network functionality.

We want to send messages to the shouty server and have it shout it back. Use the socket system:

- Open the socket with the socket system call (localhost)
- Look at Go socket package — use the lowest level abstraction provided by the language

## Resources

- What is a socket: https://csprimer.com/watch/sockets/
- How do I know what socket system calls a library is making?: https://csprimer.com/watch/socket-syscalls/
- The many differences between TCP and UDP: https://csprimer.com/watch/tcp-udp/
- A quick tour of netcat: https://csprimer.com/watch/nc/
- What is the loopback interface (ie "localhost"): https://csprimer.com/watch/localhost/
- What does it mean to bind to a port?: https://csprimer.com/watch/bind-to-port/
- What is a system call? (high level explanation): https://csprimer.com/watch/syscall-basic/
- What is a file descriptor?: https://csprimer.com/watch/fds/

## Brainstorm

Problem: create a server that returns the strings in uppercase (applies to all characters)
- Non-characters: returns as-is
- Characters: returns uppercase if it's undercase, else returns as-is

Example
- hello -> HELLO
- he110 -> HE110
- hWh aA -> HWH AA

Data structures
- Socket interface: receive data
- Strings

Algo
- Server parses input
  - If the character is a lowercase character, uppercase it-> ASCII or UTF-8
  - Append to string or []byte
- Socket allows the server to accept connections
  - What protocols to use?
    - Local connections: syscall.AF_UNIX // AF_INTET
      - UNIX SOCKET: won't be able to connect with HTTP because it doesn't use the loopback functionality, unlike AF_INTET
    - SOCK_DGRAM: avoids having to create a connection, which would require specifying listening and accepting a connection
- Bind the socket to an address? Address being localhost and port being the port where we're sending message
- Receive message from the port, process it, and send it back to the sender
- Sendmsg: when using a TCP connection, Sendto when not (e.g. UDP)

Go comment from Josh:
- Errors are generally treated as values in Go: try/catch isn't a thing, encouraged to return and handle errors as values rather than panic
