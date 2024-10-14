package fragmento

import "encoding/binary"

type Header struct {
	id         uint32
	fragmented bool
	index      uint16
	total      uint16
}

func (h *Header) ID() uint32 {
	return h.id
}

func (h *Header) Fragmented() bool {
	return h.fragmented
}

func (h *Header) Index() uint16 {
	return h.index
}

func (h *Header) Total() uint16 {
	return h.total
}

func (h *Header) PackBits() byte {
	var x byte = 0
	if h.fragmented {
		x |= 1
	}
	return x
}

func (h *Header) Serialize() []byte {
	bytes := make([]byte, HEADER_SIZE)
	binary.LittleEndian.PutUint32(bytes[0:], h.id)
	bytes[4] = h.PackBits()
	binary.LittleEndian.PutUint16(bytes[5:], h.index)
	binary.LittleEndian.PutUint16(bytes[7:], h.total)

	return bytes
}

func NewHeader(id uint32, fragmented bool, idx uint16, total uint16) *Header {
	return &Header{
		id:         id,
		fragmented: fragmented,
		index:      idx,
		total:      total,
	}
}
