# HTTP Header Server

Problem:

- Create a server that receives HTTP/1.1 connection requests
- When a request is received, parse the headers, reformat them as JSON, and return it in the response body
- Stretch goal: make sure I've received the entire request
  - Sometimes a request (or even its headers) doesn't fit in a single packet
  - The challenge here is to keep calling `Read` until there's no more incoming packets
