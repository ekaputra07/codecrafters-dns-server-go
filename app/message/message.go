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
