package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/pkg/errors"
)

const (
	maxReadLength = 2048
)

// writeMessageType writes a single byte as the value for MessageType
func writeMessageType(w io.Writer, t MessageType) error {
	byteType := byte(t)
	_, err := w.Write([]byte{byteType})
	if err != nil {
		return errors.Wrap(err, "error writing message type byte")
	}

	return nil
}

// readMessageType reads a single byte as the value for MessageType
func readMessageType(r io.Reader) (MessageType, error) {
	// allocate buffer
	b := make([]byte, 1)

	// we use io.ReadFull to ensure thatwe read all the bytes we need
	_, err := io.ReadFull(r, b)
	if err != nil {
		return MessageTypeUnknown, errors.Wrap(err, "error reading message type byte")
	}

	return MessageType(b[0]), nil
}

// writeString writes a string to a writer as
// * The string length as a little endian unsigned 32 bits integer (4 bytes)
// * The string itself as UTF-8 raw bytes
func writeString(w io.Writer, s string) error {
	// allocate buffer for string length
	lengthBuf := make([]byte, 4)
	// encode it as a little endian uint32
	binary.LittleEndian.PutUint32(lengthBuf, uint32(len(s)))

	// write the string length
	_, err := w.Write(lengthBuf)
	if err != nil {
		return errors.Wrap(err, "error writing string length")
	}

	// write the string itself
	_, err = w.Write([]byte(s))
	if err != nil {
		return errors.Wrap(err, "error writing string")
	}

	return nil
}

// readString reads a string from a reader by
// * Reading a 4 bytes little endian unsigned int to determninedthe string length
// * Reading the number of bytes previously determined to read the string
func readString(r io.Reader) (string, error) {
	// allocate buffer for string length
	lengthBuf := make([]byte, 4)

	// read the bytes to the buffer
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return "", errors.Wrap(err, "error reading string length")
	}

	// decode the uint32
	len := binary.LittleEndian.Uint32(lengthBuf)

	// protect against too large strings
	if len > maxReadLength {
		return "", fmt.Errorf("max read length exceeded: %d > %d", len, maxReadLength)
	}

	// allocate buffer for string
	buf := make([]byte, len)

	// read the actual string
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", errors.Wrap(err, "error reading string")
	}

	// check that what we read is actually a valid utf8 string
	if !utf8.Valid(buf) {
		return "", fmt.Errorf("error reading string: invalid utf8 string: %q", string(buf))
	}

	return string(buf), nil
}
