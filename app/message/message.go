package message

func NewMessage() []byte {
	var m []byte

	header := newHeader()
	header.setID(1234)
	header.setResponse()
	header.setQDCOUNT(1)

	question := newQuestion("codecrafters.io", 1, 1)

	m = append(m, header...)
	m = append(m, question...)
	return m
}
