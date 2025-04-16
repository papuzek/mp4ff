package hevc

import (
	"encoding/binary"
	"fmt"
)

// GetNalusFromSample - get nalus by following 4 byte length fields
func GetNalusFromSample(sample []byte) ([][]byte, error) {
	length := len(sample)
	if length < 4 {
		return nil, fmt.Errorf("less than 4 bytes, No NALUs")
	}
	naluList := make([][]byte, 0, 2)
	var pos uint32 = 0
	for pos < uint32(length-4) {
		naluLength := binary.BigEndian.Uint32(sample[pos : pos+4])
		pos += 4
		// Check for potential overflow or invalid length
		if naluLength > uint32(length)-pos {
			return nil, fmt.Errorf("NALU length %d exceeds remaining sample size %d at position %d", naluLength, uint32(length)-pos, pos-4)
		}
		naluList = append(naluList, sample[pos:pos+naluLength])
		pos += naluLength
	}
	// Check if we consumed the whole buffer exactly
	if pos != uint32(length) {
		// This might indicate trailing data or a truncated NALU length/data at the end
		// Depending on strictness required, you might want to return an error here.
		// For now, we'll allow trailing data.
		// fmt.Printf("Warning: Trailing data of size %d left after parsing NALUs\n", uint32(length)-pos)
	}

	return naluList, nil
}
