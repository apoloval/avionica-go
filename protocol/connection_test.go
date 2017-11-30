package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionHeader_Decode(t *testing.T) {
	buf := append([]byte("AVIONICA"), 1, 22)
	var conn connectionHeader
	err := decodeConnectionHeader(bytes.NewReader(buf), &conn)

	assert.NoError(t, err)
	assert.Equal(t, byte(1), conn.VersionMajor)
	assert.Equal(t, byte(22), conn.VersionMinor)
}

func TestConnectionHeader_DecodeFailsIfHeaderIsInvalid(t *testing.T) {
	buf := append([]byte("PATATA01"), 1, 22)
	var conn connectionHeader
	err := decodeConnectionHeader(bytes.NewReader(buf), &conn)

	assert.EqualError(t, err, "invalid connection header: PATATA01")
}
