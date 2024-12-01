package udpxy

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/ipv4"
)

func StreamFromUdp(res http.ResponseWriter, req *http.Request) {
	// set headers for MPEG-TS
	res.Header().Set("Content-Type", "video/mp2t")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Connection", "keep-alive")
	res.WriteHeader(http.StatusOK)

	// connect to udp source
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8208")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%+v \n", addr)

	ifi, _ := net.InterfaceByName("eno1")
	fmt.Printf("%+v \n", ifi)

	conn, err := net.ListenPacket("udp4", "239.1.1.8:8208")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	src := net.IPv4(239, 1, 1, 8)

	packetConn := ipv4.NewPacketConn(conn)

	if err := packetConn.JoinGroup(ifi, &net.UDPAddr{IP: src}); err != nil {
		fmt.Println("Group Join Error", err)
		return
	}

	if err := packetConn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		fmt.Println("Control message error", err)
	}

	flusher, ok := res.(http.Flusher)
	if !ok {
		http.Error(res, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	buffer := make([]byte, 1500)

	for {
		n, cm, _, err := packetConn.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}

		if cm.Dst.Equal(src) { // check if multicast destination is the source to avoid packet duplication
			_, writeErr := res.Write(buffer[:n])
			if writeErr != nil {
				fmt.Println(writeErr)
				break
			}
		}

		flusher.Flush()
	}
}
