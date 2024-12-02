package udpxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"

	"golang.org/x/net/ipv4"
)

type Udpxy struct {
	address string
	buffer  int
	Int     string
}

func (upxy *Udpxy) StreamFromUdp(res http.ResponseWriter, req *http.Request) {
	// set headers for MPEG-TS
	res.Header().Set("Content-Type", "video/mp2t")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Connection", "keep-alive")
	res.WriteHeader(http.StatusOK)

	upxy.buffer = 1500

	re := regexp.MustCompile(`udp/(.+)`)
	addressMatch := re.FindStringSubmatch(req.URL.Path)
	if len(addressMatch) != 2 {
		return
	}
	sourceAddress := addressMatch[1]

	conn, err := net.ListenPacket("udp4", sourceAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	packetConn := ipv4.NewPacketConn(conn)

	var joinGroupInt *net.Interface = nil
	if ifi, err := net.InterfaceByName(upxy.Int); err == nil {
		joinGroupInt = ifi
	}

	sourceIP, _, _ := net.SplitHostPort(sourceAddress)
	parsedSourceIP := net.ParseIP(sourceIP)
	if err := packetConn.JoinGroup(joinGroupInt, &net.UDPAddr{IP: parsedSourceIP}); err != nil {
		log.Fatal("Group Join Error", err)
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

		if cm.Dst.Equal(parsedSourceIP) { // check if multicast destination is the source to avoid packet duplication
			_, writeErr := res.Write(buffer[:n])
			if writeErr != nil {
				fmt.Println(writeErr)
				break
			}
		}

		flusher.Flush()
	}
}

//func getInterface(intString string) *net.Interface {
//	inif := net.ParseIP(intString)
//
//	ifaces, err := net.Interfaces()
//	if err != nil {
//		return nil
//	}
//
//	for _, iface := range ifaces {
//		if addrs, err := iface.Addrs(); err == nil {
//			for _, addr := range addrs {
//				if iip, _, err := net.ParseCIDR(addr.String()); err == nil {
//					fmt.Printf("%+v %+v \n", iip, inif)
//					if iip.Equal(inif) {
//						return &iface
//					}
//				}
//			}
//		}
//	}
//	return nil
//}
