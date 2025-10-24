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
// | TTL          | 4 bytes
// +--+--+--+--+--+
// | Length       | 2 bytes
// +--+--+--+--+--+
// | Data         | 4 bytes
// +--+--+--+--+--+

type Answer struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  string
}

func (a Answer) ToBytes() []byte {
	var bytes []byte

	// append labels
	parts := strings.SplitSeq(a.Name, ".")

	for part := range parts {
		bytes = append(bytes, byte(len(part)))
		bytes = append(bytes, []byte(part)...)
	}

	// end of labels
	bytes = append(bytes, 0)

	// append QTYPE and QCLASS
	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, a.Type)

	classBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(classBuf, a.Class)

	bytes = append(bytes, typeBuf...)
	bytes = append(bytes, classBuf...)

	// append TTL
	ttlBuf := make([]byte, 4)
	binary.BigEndian.AppendUint32(ttlBuf, a.TTL)
	bytes = append(bytes, ttlBuf...)

	dataBytes := []byte(a.Data)

	// append length
	lengthBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBuf, uint16(len(dataBytes)))
	bytes = append(bytes, lengthBuf...)

	// append DATA
	bytes = append(bytes, []byte(a.Data)...)
	return bytes
}
