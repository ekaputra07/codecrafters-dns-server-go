package message

import (
	"encoding/binary"
	"strings"
)

func newQuestion(qname string, qtype, qclass uint16) (q []byte, count uint16) {

	// append labels
	parts := strings.SplitSeq(qname, ".")
	count = 0

	for part := range parts {
		q = append(q, byte(len(part)))
		q = append(q, []byte(part)...)
		count += 1
	}

	// end of labels
	q = append(q, 0)

	// append QTYPE and QCLASS
	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, qtype)

	classBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(classBuf, qclass)

	q = append(q, typeBuf...)
	q = append(q, classBuf...)
	return
}
