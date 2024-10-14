package fragmento

import (
	"encoding/binary"
	"hash/crc32"
)

const (
	HEADER_SIZE      int = 9
	PAYLOADSIZE_SIZE int = 2
	CHECKSUM_SIZE    int = 4
	IP_SIZE          int = 40
	UDP_SIZE         int = 8
	RESERVED         int = 100
	DEFAULT_MTU_SIZE int = 1500
	MAX_PAYLOAD_SIZE int = DEFAULT_MTU_SIZE - (RESERVED + HEADER_SIZE + PAYLOADSIZE_SIZE + CHECKSUM_SIZE + UDP_SIZE + IP_SIZE)
)

type Fragment struct {
	Header   *Header
	Size     uint16
	Payload  []byte
	Checksum uint32
}

func (f *Fragment) Serialize() []byte {
	size := HEADER_SIZE + PAYLOADSIZE_SIZE + len(f.Payload) + CHECKSUM_SIZE
	buff := make([]byte, size)

	copy(buff[0:], f.Header.Serialize())
	binary.LittleEndian.PutUint16(buff[HEADER_SIZE:], f.Size)
	copy(buff[HEADER_SIZE+PAYLOADSIZE_SIZE:], f.Payload)
	binary.LittleEndian.PutUint32(buff[HEADER_SIZE+PAYLOADSIZE_SIZE+len(f.Payload):], f.Checksum)

	return buff
}

func calculateChecksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func NewFragment(data []byte, header *Header) *Fragment {
	if len(data) > MAX_PAYLOAD_SIZE {
		panic("recieved too much bytes")
	}

	checksum := calculateChecksum(data)

	return &Fragment{
		Header:   header,
		Size:     uint16(len(data)),
		Payload:  data,
		Checksum: checksum,
	}
}
