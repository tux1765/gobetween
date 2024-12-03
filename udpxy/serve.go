package udpxy

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(address string, port string, listenInterface string, packetBufferSize int) {
	serverPort := fmt.Sprintf("%s:%s", address, port)

	myUdpxy := &Udpxy{Int: listenInterface, Buffer: packetBufferSize}
	http.HandleFunc("/udp/{udpAddress}", myUdpxy.StreamFromUdp)
	log.Printf("Serving MPEG-TS UDP stream on %s/udp", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
