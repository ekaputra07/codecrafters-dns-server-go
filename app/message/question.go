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
