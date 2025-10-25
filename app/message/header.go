package message

import (
	"encoding/binary"
	"fmt"
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

func (h Header) String() string {
	return fmt.Sprintf("ID: %d,\nOPCODE:%d,\nAA: %t,\nTC: %t,\nRD: %t,\nRA: %t,\nZ: %d,\nRCODE: %d,\nQDCOUNT: %d,\nANCOUNT: %d,\nNSCOUNT: %d,\nARCOUNT: %d",
		h.ID,
		*h.OPCODE,
		h.AA,
		h.TC,
		h.RD,
		h.RA,
		*h.Z,
		*h.RCODE,
		*h.QDCOUNT,
		*h.ANCOUNT,
		*h.NSCOUNT,
		*h.ARCOUNT,
	)
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
	if h.ARCOUNT != nil {
		binary.BigEndian.PutUint16(bytes[10:], *h.ARCOUNT)
	}
	return bytes
}

func ParseHeader(bytes []byte) Header {
	id := binary.BigEndian.Uint16(bytes)

	// OPCODE is 4-bits long, starting from the 2nd bit of the 3rd byte (index 2)
	opcode := uint8((bytes[2] & 0b01111000) >> 3)

	// AA is 1-bit long, the 6th bit of the 3rd byte (index 2)
	aa := (bytes[2] & 0b00000100) >> 2

	// TC is 1-bit long, the 7th bit of the 3rd byte (index 2)
	tc := (bytes[2] & 0b00000010) >> 1

	// RD is 1-bit long, the least significant bit of the 3rd byte (index 2)
	rd := bytes[2] & 0b00000001

	// RA is 1-bit long, the most significant bit of the 4th byte (index 3)
	ra := (bytes[3] & 0b10000000) >> 7

	// Z is 3-bits long, starting from the 2nd bit of the 4th byte (index 3)
	z := uint8(bytes[3] & 0b01110000)

	// RCODE is 4-bits long, starting from the 5th bit of the 4th byte (index 3)
	rcode := uint8(bytes[3] & 0b00001111)

	// QDCOUNT is 16-bits long (2 bytes) starts at 5th byte (index 4)
	qdcount := binary.BigEndian.Uint16(bytes[4:])

	// ANCOUNT is 16-bits long (2 bytes) starts at 7th byte (index 6)
	ancount := binary.BigEndian.Uint16(bytes[6:])

	// NSCOUNT is 16-bits long (2 bytes) starts at 9th byte (index 8)
	nscount := binary.BigEndian.Uint16(bytes[8:])

	// ARCOUNT is 16-bits long (2 bytes) starts at 11th byte (index 10)
	arcount := binary.BigEndian.Uint16(bytes[10:])

	return Header{
		ID:      id,
		OPCODE:  &opcode,
		AA:      int(aa) == 1,
		TC:      int(tc) == 1,
		RD:      int(rd) == 1,
		RA:      int(ra) == 1,
		Z:       &z,
		RCODE:   &rcode,
		QDCOUNT: &qdcount,
		ANCOUNT: &ancount,
		NSCOUNT: &nscount,
		ARCOUNT: &arcount,
	}
}
