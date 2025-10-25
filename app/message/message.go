package message

type Section interface {
	ToBytes() []byte
}

type Message struct {
	Header    Header
	Questions []Question
	Answers   []Answer
}

func (m Message) ToBytes() []byte {
	bytes := m.Header.ToBytes()

	// questions section
	qBuff := []byte{}
	for _, q := range m.Questions {
		qBuff = append(qBuff, q.ToBytes()...)
	}

	// answers section
	aBuff := []byte{}
	for _, a := range m.Answers {
		aBuff = append(aBuff, a.ToBytes()...)
	}

	// combine all
	bytes = append(bytes, qBuff...)
	bytes = append(bytes, aBuff...)
	return bytes
}
