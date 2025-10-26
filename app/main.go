package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/app/message"
)

func main() {
	var forwardAddr *net.UDPAddr
	if len(os.Args) == 3 {
		fwdAddr, err := net.ResolveUDPAddr("udp", os.Args[2])
		if err != nil {
			log.Println("Failed to resolve upstream DNS server address:", err)
		}
		fmt.Println("Upstream address: ", fwdAddr)
		forwardAddr = fwdAddr
	}

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

		// either forward the request or build a response
		var response []byte
		if forwardAddr != nil {
			response = forwardRequest(receivedData, forwardAddr)
		} else {
			response = buildResponse(receivedData)
		}

		// send response
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func forwardRequest(request []byte, addr *net.UDPAddr) []byte {
	header := message.ParseHeader(request[:12])
	questions := message.ParseQuestions(*header.QDCOUNT, request[12:])

	if len(questions) == 1 {
		// single question, forward as is
		fmt.Printf("Forwarding request for question: %+v\n", questions[0])
		return makeForwardRequest(request, addr)

	} else if len(questions) > 1 {
		fmt.Printf("Forwarding request for multiple questions: %+v\n", questions)

		// send each question separately then combine responses
		qdcount := uint16(1)
		header.QDCOUNT = &qdcount

		answersBytes := []byte{}

		for _, q := range questions {
			fmt.Printf("- Forwarding request for question: %+v\n", q)
			singleQuestionRequest := append(header.ToBytes(), q.ToBytes()...)
			singleResponse := makeForwardRequest(singleQuestionRequest, addr)
			singleQuestionLength := len(singleQuestionRequest)

			// append answers from singleResponse to the main response
			if len(singleResponse) > singleQuestionLength {
				answer := singleResponse[singleQuestionLength:]
				answersBytes = append(answersBytes, answer...)
			}
		}
		// build final response
		header.Response = true
		ancount := uint16(len(questions))
		header.QDCOUNT = &ancount // set QDCOUNT to number of questions
		header.ANCOUNT = &ancount // set ANCOUNT to number of answers == number of questions
		respBytes := message.Message{
			Header:    header,
			Questions: questions,
		}.ToBytes()
		respBytes = append(respBytes, answersBytes...)
		return respBytes
	}
	fmt.Println("No questions found in the request to forward.")
	return request
}

func makeForwardRequest(request []byte, addr *net.UDPAddr) []byte {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Failed to dial upstream DNS server:", err)
		return nil
	}
	defer conn.Close()

	_, err = conn.Write(request)
	if err != nil {
		fmt.Println("Failed to send request to upstream DNS server:", err)
		return nil
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read response from upstream DNS server:", err)
		return nil
	}
	return buf[:n]
}

func buildResponse(request []byte) []byte {
	attl := uint32(60)
	data := "8888"

	header := message.ParseHeader(request[:12])
	questions := message.ParseQuestions(*header.QDCOUNT, request[12:])

	var answers []message.Answer
	for _, q := range questions {
		answer := message.Answer{
			Name:  q.Name,
			Type:  q.Type,
			Class: q.Class,
			TTL:   attl,
			Data:  data,
		}
		answers = append(answers, answer)
	}

	// RCODE: 0 (no error) if OPCODE is 0 (standard query) else 4 (not implemented)
	if *header.OPCODE == uint8(0) {
		header.RCODE = header.OPCODE
	} else {
		rcode := uint8(4)
		header.RCODE = &rcode
	}

	// set ancount
	ancount := uint16(len(answers))
	header.ANCOUNT = &ancount

	resp := &message.Message{
		Header:    header,
		Questions: questions,
		Answers:   answers,
	}

	return resp.ToBytes()
}
