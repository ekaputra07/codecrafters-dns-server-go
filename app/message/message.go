package message

type Section interface {
	ToBytes() []byte
}

type Message struct {
	Header   Section
	Question Section
	Answer   Section
}

func (m Message) ToBytes() []byte {
	bytes := m.Header.ToBytes()
	bytes = append(bytes, m.Question.ToBytes()...)
	bytes = append(bytes, m.Answer.ToBytes()...)
	return bytes
}

func NewMessage(ID uint16, qName string, qType, qClass uint16, aTTL uint32, data string) []byte {
	qdcount := uint16(1)
	ancount := uint16(1)

	message := Message{
		Header:   Header{ID: ID, QDCOUNT: &qdcount, ANCOUNT: &ancount},
		Question: Question{Name: qName, Type: qType, Class: qClass},
		Answer:   Answer{Name: qName, Type: qType, Class: qClass, TTL: aTTL, Data: data},
	}

	return message.ToBytes()
}
