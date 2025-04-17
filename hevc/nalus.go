package hevc

import (
	"encoding/binary"
	"fmt"
)

// GetNalusFromSample - extracts NALUs from a sample based on the specified NALU length size.
// naluLengthSize should be 1, 2, or 4.
func GetNalusFromSample(sample []byte, naluLengthSize uint32) ([][]byte, error) {
	length := len(sample)
	if length < int(naluLengthSize) { // Check if sample is smaller than the length prefix size
		return nil, fmt.Errorf("sample size %d is less than naluLengthSize %d", length, naluLengthSize)
	}

	naluList := make([][]byte, 0, 2) // Pre-allocate assuming at least a few NALUs
	var pos uint32 = 0

	for pos < uint32(length) {
		// Check if there are enough bytes left for the length prefix
		if pos+naluLengthSize > uint32(length) {
			// This likely indicates trailing data or an incomplete last NALU prefix
			break // Stop processing
		}

		var naluLength uint32
		switch naluLengthSize {
		case 1:
			naluLength = uint32(sample[pos])
		case 2:
			naluLength = uint32(binary.BigEndian.Uint16(sample[pos : pos+naluLengthSize]))
		case 4:
			naluLength = binary.BigEndian.Uint32(sample[pos : pos+naluLengthSize])
		default:
			return nil, fmt.Errorf("invalid naluLengthSize: %d (must be 1, 2, or 4)", naluLengthSize)
		}
		pos += naluLengthSize

		// Check for potential overflow or invalid length
		if naluLength == 0 {
			// Some encoders might add 0-length NALUs as padding, treat as warning or error?
			// For now, let's skip them. Add logging if desired.
			// log.Printf("Warning: Skipping 0-length NALU at offset %d", pos-naluLengthSize)
			continue
		}
		if pos+naluLength > uint32(length) {
			return nil, fmt.Errorf("NALU length %d exceeds remaining sample size %d at position %d", naluLength, uint32(length)-pos, pos-naluLengthSize)
		}

		// Extract NALU data
		naluData := sample[pos : pos+naluLength]
		naluList = append(naluList, naluData)
		pos += naluLength
	}

	// Check if we consumed the whole buffer exactly (optional, depending on strictness)
	if pos != uint32(length) {
		// fmt.Printf("Warning: Trailing data of size %d left after parsing NALUs\n", uint32(length)-pos)
	}

	if len(naluList) == 0 && length > 0 {
		// If we didn't find any NALUs but had data, it might indicate an issue
		// or maybe the sample format is different than expected.
		// log.Printf("Warning: No NALUs extracted from sample of size %d", length)
	}

	return naluList, nil
}
