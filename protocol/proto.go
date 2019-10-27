package protocol

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type MessageType byte

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeError
	MessageTypeChat
)

// String provide a human readable representation of a MessageType byte
// value
func (t MessageType) String() string {
	switch t {
	case MessageTypeError:
		return "error"
	case MessageTypeChat:
		return "chat"

	case MessageTypeUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

// Message represent a message read or to be written to a connection
// The Type() method determines its type
type Message interface {
	Type() MessageType

	// private interface methods
	decode(r io.Reader) error
	encode(w io.Writer) error
}

// Read reads and deocdes a message from a reader
func Read(r io.Reader) (Message, error) {
	t, err := readMessageType(r)
	if err != nil {
		return nil, err
	}

	// Select the correct message struct from the message type
	var msg Message
	switch t {
	case MessageTypeError:
		msg = &ErrorMessage{}
	case MessageTypeChat:
		msg = &ChatMessage{}

	case MessageTypeUnknown:
		// We should never read an unknown message type, fallback to error case
		fallthrough
	default:
		return nil, fmt.Errorf("unhandled message type %d", t)
	}

	// Call the message implementation to read and decode the actual content of
	// the message
	err = msg.decode(r)
	if err != nil {
		return nil, errors.Wrapf(err, "error decoding message %q", t)
	}

	return msg, nil
}

// Write write an enoded message to a writer
func Write(w io.Writer, msg Message) error {
	return msg.encode(w)
}
