package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func NewUUID() (string, error) {
	var (
		uuid [16]byte
		b    [36]byte
	)
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	hex.Encode(b[:], uuid[:4])
	b[8] = '-'
	hex.Encode(b[9:13], uuid[4:6])
	b[13] = '-'
	hex.Encode(b[14:18], uuid[6:8])
	b[18] = '-'
	hex.Encode(b[19:23], uuid[8:10])
	b[23] = '-'
	hex.Encode(b[24:], uuid[10:])
	return string(b[:]), nil
}
