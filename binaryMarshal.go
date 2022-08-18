package sftp

import "encoding/binary"

func MarshalBinary(r RequestResponse) ([]byte, error) {
	b := make([]byte, 0, 128)

	// encode version(1-byte) and (ResponseCode 1-byte)
	b = append(b,
		byte(r.Version()),
		byte(r.Code()),
	)

	// encode Header len
	hLen := make([]byte, 2)
	l := uint16(len(r.Header()))
	binary.BigEndian.PutUint16(hLen, l)
	b = append(b, hLen...)

	// encode Header
	b = append(b, r.Header()...)
	return b, nil
}
