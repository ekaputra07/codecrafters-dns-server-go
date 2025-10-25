package message

import (
	"encoding/binary"
	"strings"
)

// Answer format:
// +--+--+--+--+--+
// | Name         | vary in bytes
// +--+--+--+--+--+
// | Type         | 2 bytes
// +--+--+--+--+--+
// | Class        | 2 bytes
// +--+--+--+--+--+

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (q Question) ToBytes() []byte {
	var bytes []byte

	// append labels
	parts := strings.SplitSeq(q.Name, ".")

	for part := range parts {
		bytes = append(bytes, byte(len(part)))
		bytes = append(bytes, []byte(part)...)
	}

	// end of labels
	bytes = append(bytes, 0)

	// append QTYPE and QCLASS
	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, q.Type)

	classBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(classBuf, q.Class)

	bytes = append(bytes, typeBuf...)
	bytes = append(bytes, classBuf...)
	return bytes
}

func ParseQuestions(count uint16, bytes []byte) []Question {
	// bytes are starting from question section (allBytes[12:])
	// so compression pointer should be offset by 12

	var questions []Question
	lastStartIndex := 0

	// loop for questions
	for range count {
		// loop for domains
		var domains []string
		for {
			b, nextStartIndex := lookupName(bytes, lastStartIndex)
			if len(b) == 0 {
				break
			}
			lastStartIndex = nextStartIndex
			domains = append(domains, string(b))
		}

		// starts from lastStartIndex
		qtype := binary.BigEndian.Uint16(bytes[lastStartIndex+1:])
		qclass := binary.BigEndian.Uint16(bytes[lastStartIndex+3:])

		q := Question{
			Name:  strings.Join(domains, "."),
			Type:  qtype,
			Class: qclass,
		}
		questions = append(questions, q)

		// move past QTYPE and QCLASS
		lastStartIndex += 5
	}
	return questions
}

func lookupName(bytes []byte, startIndex int) ([]byte, int) {
	lengthValue := int(bytes[startIndex])

	if lengthValue == 0 {
		// 0 lengthValue found, this might be the delimiter (end of the `name`)
		// return empty bytes and original startIndex
		return []byte{}, startIndex
	}

	if lengthValue == 192 {
		// pointer found (decimal 192 == binary 11000000 == hex C0)
		pointerOffset := int(binary.BigEndian.Uint16(bytes[startIndex:])) & 0x3FFF
		return lookupName(bytes, pointerOffset-12) // -12: we're already at question section
	}

	ifrom := startIndex + 1 // after the length byte
	ito := ifrom + lengthValue

	return bytes[ifrom:ito], ito
}
