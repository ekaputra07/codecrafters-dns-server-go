package main

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/app/message"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// return only the header
		_, err = udpConn.WriteToUDP(buildResponse(receivedData), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func buildResponse(request []byte) []byte {
	attl := uint32(60)
	data := "8888"

	header := message.ParseHeader(request[:12])
	question := message.ParseQuestion(request[12:])

	// RCODE: 0 (no error) if OPCODE is 0 (standard query) else 4 (not implemented)
	if *header.OPCODE == uint8(0) {
		header.RCODE = header.OPCODE
	} else {
		rcode := uint8(4)
		header.RCODE = &rcode
	}

	// set qdcount and ancount
	ancount := uint16(1)
	header.ANCOUNT = &ancount

	resp := &message.Message{
		Header:   header,
		Question: question,
		Answer: message.Answer{
			Name:  question.Name,
			Type:  question.Type,
			Class: question.Class,
			TTL:   attl, Data: data,
		},
	}

	return resp.ToBytes()
}
