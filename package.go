package fragmento

import (
	"encoding/binary"
	"fmt"
)

func FragmentData(id uint32, data []byte) []Fragment {
	count := len(data) / MAX_PAYLOAD_SIZE
	if len(data)%MAX_PAYLOAD_SIZE != 0 {
		count++
	}

	fragmented := false
	if count > 0 {
		fragmented = true
	}

	var fragments []Fragment

	for i := 0; i < count; i++ {
		header := NewHeader(id, fragmented, uint16(i), uint16(count))
		start := MAX_PAYLOAD_SIZE * i
		end := start + MAX_PAYLOAD_SIZE
		if end > len(data) {
			end = len(data)
		}
		part := data[start:end]

		fragment := NewFragment(part, header)
		fragments = append(fragments, *fragment)
	}

	return fragments
}

func FromFragments(frags []Fragment) []byte {
	var size uint16 = 0
	for _, x := range frags {
		size += x.Size
	}
	bytes := make([]byte, size)

	var offset int
	for _, x := range frags {
		copy(bytes[offset:], x.Payload)
		offset += int(x.Size)
	}

	return bytes
}

func Deserialize(data []byte) (*Fragment, error) {
	if len(data) < HEADER_SIZE+PAYLOADSIZE_SIZE+CHECKSUM_SIZE {
		return nil, fmt.Errorf("data too short")
	}

	id := binary.LittleEndian.Uint32(data[0:4])
	fragmented := (data[4] & 1) != 0
	idx := binary.LittleEndian.Uint16(data[5:7])
	total := binary.LittleEndian.Uint16(data[7:9])
	header := NewHeader(id, fragmented, idx, total)

	size := binary.LittleEndian.Uint16(data[HEADER_SIZE:])
	if len(data) < HEADER_SIZE+PAYLOADSIZE_SIZE+int(size)+CHECKSUM_SIZE {
		return nil, fmt.Errorf("data length mismatch")
	}

	payload := data[HEADER_SIZE+PAYLOADSIZE_SIZE : HEADER_SIZE+PAYLOADSIZE_SIZE+int(size)]
	checksum := binary.LittleEndian.Uint32(data[HEADER_SIZE+PAYLOADSIZE_SIZE+int(size):])

	fragment := &Fragment{
		Header:   header,
		Size:     size,
		Payload:  payload,
		Checksum: checksum,
	}

	return fragment, nil
}
