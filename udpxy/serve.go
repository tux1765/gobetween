package udpxy

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(address string, port string) {
	serverPort := fmt.Sprintf("%s:%s", address, port)
	http.HandleFunc("/stream", StreamFromUdp)
	log.Printf("Serving MPEG-TS UDP stream on %s/stream", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
