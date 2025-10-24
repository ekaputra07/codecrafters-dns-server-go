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

func ParseQuestion(bytes []byte) Question {
	var names []string
	lastStartIndex := 0

	// loop for names
	for true {
		b, nextStartIndex := lookupName(bytes, lastStartIndex)
		if len(b) == 0 {
			break
		}
		lastStartIndex = nextStartIndex
		names = append(names, string(b))
	}

	// starts from lastStartIndex
	qtype := binary.BigEndian.Uint16(bytes[lastStartIndex+1:])
	qclass := binary.BigEndian.Uint16(bytes[lastStartIndex+3:])

	return Question{
		Name:  strings.Join(names, "."),
		Type:  qtype,
		Class: qclass,
	}
}

func lookupName(bytes []byte, startIndex int) ([]byte, int) {
	lengthValue := int(bytes[startIndex])

	// 0 lengthValue found, this might be the delimiter (end of the `name`)
	// return empty bytes and original startIndex
	if lengthValue == 0 {
		return []byte{}, startIndex
	}
	ifrom := startIndex + 1 // after the length byte
	ito := ifrom + lengthValue

	return bytes[ifrom:ito], ito
}
