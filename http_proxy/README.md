# HTTP Proxy

HTTP is so successful that most SWE can't tell you the difference between the web and the internet — the web is the internet's killer app. (see from 0:50m)

This is about implementing a reverse proxy.

Receive HTTP requests for an upstream server and forward the server's response back to the client.

Log the number of bytes received/forwarded in each step.

Only support HTTP/1.1.
