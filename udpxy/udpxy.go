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
	Buffer  int
	Int     string
}

func (udpxy *Udpxy) StreamFromUdp(res http.ResponseWriter, req *http.Request) {
	// set headers for MPEG-TS
	res.Header().Set("Content-Type", "video/mp2t")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Connection", "keep-alive")
	res.WriteHeader(http.StatusOK)

	re := regexp.MustCompile(`udp/(.+)`)
	addressMatch := re.FindStringSubmatch(req.URL.Path)
	if len(addressMatch) != 2 {
		return
	}
	multicastSourceAddr := addressMatch[1]

	conn, err := net.ListenPacket("udp4", multicastSourceAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	packetConn := ipv4.NewPacketConn(conn)

	var joinGroupInt *net.Interface = nil
	if ifi, err := net.InterfaceByName(udpxy.Int); err == nil {
		joinGroupInt = ifi
	}

	unparsedMulticastIP, _, _ := net.SplitHostPort(multicastSourceAddr)
	parsedSourceMulticastIP := net.ParseIP(unparsedMulticastIP)
	if err := packetConn.JoinGroup(joinGroupInt, &net.UDPAddr{IP: parsedSourceMulticastIP}); err != nil {
		log.Fatal("Group Join Error", err)
		return
	}

	if err := packetConn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		log.Fatal("Control message error", err)
		return
	}

	flusher, ok := res.(http.Flusher)
	if !ok {
		http.Error(res, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	buffer := make([]byte, udpxy.Buffer)

	for {
		n, cm, _, err := packetConn.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}

		// if int to listen on is supplied only get packets from that int
		if joinGroupInt != nil && joinGroupInt.Index != cm.IfIndex {
			continue
		}

		// verify multicast destination is the source
		if cm.Dst.Equal(parsedSourceMulticastIP) {
			_, writeErr := res.Write(buffer[:n])
			if writeErr != nil {
				fmt.Println(writeErr)
				break
			}
		}

		flusher.Flush()
	}
}

// TODO: implement this function to accept interface IP addresses
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
