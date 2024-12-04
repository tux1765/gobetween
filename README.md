# gobetween
gobetween is a UDP-toHTTP multicast traffic relay.
It forwards UDP traffic from a given multicast address to the requesting HTTP client.

Install the package with `go install github.com/tux1765/gobetween/gb` and start the service using ```gb``` command

Send a GET request to `http://ip:port/udp/<multicastip:port>` to start the proxy.