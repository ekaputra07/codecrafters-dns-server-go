package main

import (
	"encoding/binary"
	"fmt"
)

type header []byte

// setID sets the ID field of the header
// ID is 16-bits long (2 bytes)
func (h header) setID(value uint16) {
	binary.BigEndian.PutUint16(h, value)
}

// setR sets the QR field of the header to Response
// QR is 1-bit long, the most significant bit of the 3rd byte (index 2)
func (h header) setResponse() {
	h[2] |= 0b10000000
}

// setOPCODE sets the OPCODE field of the header
// OPCODE is 4-bits long, starting from the 2nd bit of the 3rd byte (index 2)
func (h header) setOPCODE(value uint8) {
	mask := 0b00001111 & value
	h[2] |= mask << 3
}

// setAA sets the AA field of the header
// AA is 1-bit long, the 6th bit of the 3rd byte (index 2)
func (h header) setAA() {
	h[2] |= 0b00000100
}

// setTC sets the TC field of the header
// TC is 1-bit long, the 7th bit of the 3rd byte (index 2)
func (h header) setTC() {
	h[2] |= 0b00000010
}

// setRD sets the RD field of the header
// RD is 1-bit long, the least significant bit of the 3rd byte (index 2)
func (h header) setRD() {
	h[2] |= 0b00000001
}

// setRA sets the RA field of the header
// RA is 1-bit long, the most significant bit of the 4th byte (index 3)
func (h header) setRA() {
	h[3] |= 0b10000000
}

// setZ sets the Z field of the header
// Z is 3-bits long, starting from the 2nd bit of the 4th byte (index 3)
func (h header) setZ(value uint8) {
	mask := 0b00000111 & value
	h[3] |= mask << 4
}

// setRCODE sets the RCODE field of the header
// RCODE is 4-bits long, starting from the 5th bit of the 4th byte (index 3)
func (h header) setRCODE(value uint8) {
	mask := 0b00001111 & value
	h[3] |= mask
}

func (h header) setQDCOUNT(value uint16) {
	binary.BigEndian.PutUint16(h[4:], value)
}

func (h header) setANCOUNT(value uint16) {
	binary.BigEndian.PutUint16(h[6:], value)
}

func (h header) setNSCOUNT(value uint16) {
	binary.BigEndian.PutUint16(h[8:], value)
}

func (h header) setARCOUNT(value uint16) {
	binary.BigEndian.PutUint16(h[10:], value)
}

func (h header) String() string {
	return fmt.Sprintf("%08b", h)
}

func newHeader() header {
	return make(header, 12)
}
