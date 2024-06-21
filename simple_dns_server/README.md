# README

Given a hostname make a DNS query for an A record and print out the response.

Working with sockets
Working with DNS
DNS is a binary-encoded protocol, so we need to know how to pack, parse, etc bytes well
Look for the key diagrams in the RFC that explain how a DNS request is structured

I'll need a destination to send my DNS query to (e.g. `8.8.8.8`)

https://www.ietf.org/rfc/rfc1035.txt

`dig wikipedia.org`

https://medium.com/@openmohan/dns-basics-and-building-simple-dns-server-in-go-6cb8e1cfe461

Key questions:

- Is DNS text-based like HTTP? Seems to be based on bytes given the message structure.
- Do I have to construct in my app code the DNS message or does the OS do that?

I need to write to a byte slice or buffer the data I need to construct (i.e. the DNS message). DNS is BIG ENDIAN!

Header

- ID: random unsigned 16 bits
- QR: `0`
- `OPCODE`: `0`
- `AA`: 0?
- TC: 0?
- RD: 0?
- RA: 0?
- Z: 0
- RCODE: 0? -> 4 bits!
- QDCOUNT: 1? -> note that it's a 16 bit signed integer!
- ANCOUNT: 0? -> note that it's a 16 bit unsigned integer!
- NSCOUNT: 0? -> note that it's a 16 bit unsigned integer!
- ARCOUNT: 0? -> note that it's a 16 bit unsigned integer!

Question

- QNAME: "google", "com", "pe" -> split by "."
- QTYPE: 1 (i.e. A records)
- QCLASS: 1 (Internet)

Empty sections

- Answer
- Authority
- Additional

Steps

- Take in user argument through CLI (can be postponed)
- Create the DNS query by creating a byte slices -> start here
- Make the DNS query with a UDP socket to a DNS server
