package message

import (
	"encoding/binary"
)

// Header format (12 bytes total):
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// | ID                                          | 2 bytes
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// | QR | OPCODE | AA | TC | RD | RA | Z | RCODE | 2 bytes
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// | QDCOUNT                                     | 2 bytes
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// | ANCOUNT                                     | 2 bytes
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// | NSCOUNT                                     | 2 bytes
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// | ARCOUNT                                     | 2 bytes
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Header struct {
	ID      uint16
	OPCODE  *uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       *uint8
	RCODE   *uint8
	QDCOUNT *uint16
	ANCOUNT *uint16
	NSCOUNT *uint16
	ARCOUNT *uint16
}

func (h Header) ToBytes() []byte {
	bytes := make([]byte, 12)

	// set ID
	binary.BigEndian.PutUint16(bytes, h.ID)

	// QR: type always Response (R)
	// QR is 1-bit long, the most significant bit of the 3rd byte (index 2)
	bytes[2] |= 0b10000000

	// OPCODE is 4-bits long, starting from the 2nd bit of the 3rd byte (index 2)
	if h.OPCODE != nil {
		mask := 0b00001111 & *h.OPCODE
		bytes[2] |= mask << 3
	}

	// AA is 1-bit long, the 6th bit of the 3rd byte (index 2)
	if h.AA {
		bytes[2] |= 0b00000100
	}

	// TC is 1-bit long, the 7th bit of the 3rd byte (index 2)
	if h.TC {
		bytes[2] |= 0b00000010
	}

	// RD is 1-bit long, the least significant bit of the 3rd byte (index 2)
	if h.RD {
		bytes[2] |= 0b00000001
	}

	// RA is 1-bit long, the most significant bit of the 4th byte (index 3)
	if h.RA {
		bytes[3] |= 0b10000000
	}

	// Z is 3-bits long, starting from the 2nd bit of the 4th byte (index 3)
	if h.Z != nil {
		mask := 0b00000111 & *h.Z
		bytes[3] |= mask << 4
	}

	// RCODE is 4-bits long, starting from the 5th bit of the 4th byte (index 3)
	if h.RCODE != nil {
		mask := 0b00001111 & *h.RCODE
		bytes[3] |= mask
	}

	// QDCOUNT is 16-bits long (2 bytes) starts at 5th byte (index 4)
	if h.QDCOUNT != nil {
		binary.BigEndian.PutUint16(bytes[4:], *h.QDCOUNT)
	}

	// ANCOUNT is 16-bits long (2 bytes) starts at 7th byte (index 6)
	if h.ANCOUNT != nil {
		binary.BigEndian.PutUint16(bytes[6:], *h.ANCOUNT)
	}

	// NSCOUNT is 16-bits long (2 bytes) starts at 9th byte (index 8)
	if h.NSCOUNT != nil {
		binary.BigEndian.PutUint16(bytes[8:], *h.NSCOUNT)
	}

	// ARCOUNT is 16-bits long (2 bytes) starts at 11th byte (index 10)
	if h.NSCOUNT != nil {
		binary.BigEndian.PutUint16(bytes[10:], *h.NSCOUNT)
	}
	return bytes
}
