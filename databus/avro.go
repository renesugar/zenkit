package databus

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
)

var (
	magicByte = []byte{0}

	// ErrAvroSerialization is thrown when a message does not include a proper Avro header
	ErrAvroSerialization = errors.New("message did not include a proper Avro header")
)

// AvroSerialize prefaces a message with a proper Avro header including the
// schema registry ID
func AvroSerialize(msg []byte, schemaID int) []byte {
	var b bytes.Buffer
	buf := &b
	buf.Write(magicByte)
	idSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(idSlice, uint32(schemaID))
	buf.Write(idSlice)
	buf.Write(msg)
	return buf.Bytes()
}

// AvroDeserialize deserializes an Avro message header and returns the schema
// registry ID and the rest of the message
func AvroDeserialize(msg []byte) (int, []byte, error) {
	if len(msg) < 5 || msg[0] != magicByte[0] {
		return 0, nil, ErrAvroSerialization
	}
	schemaID := int(binary.BigEndian.Uint32(msg[1:5]))
	return schemaID, msg[5:], nil
}
