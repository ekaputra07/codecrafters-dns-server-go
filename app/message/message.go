package message

import "fmt"

func NewMessage() []byte {
	var m []byte

	header := newHeader()
	header.setID(1234)
	header.setResponse()

	question, count := newQuestion("codecrafters.io", 1, 1)
	header.setQDCOUNT(count)

	m = append(m, header...)
	m = append(m, question...)
	m = append(m, make([]byte, 512-33)...)
	fmt.Println(m)
	fmt.Println(len(m))
	return m
}
