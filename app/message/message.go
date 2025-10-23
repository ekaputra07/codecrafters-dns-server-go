package message

func NewMessage() []byte {
	m := make([]byte, 512)

	header := newHeader()
	header.setID(1234)
	header.setResponse()

	question, count := newQuestion("codecrafters.io", 1, 1)
	header.setQDCOUNT(count)

	m = append(m, header...)
	m = append(m, question...)
	return m
}
