package protocol

import "io"

// ErrorMessage represents an error
type ErrorMessage struct {
	Message string
}

// Error implements the error interface
func (m ErrorMessage) Error() string {
	return m.Message
}

// Type returns the message type
func (m ErrorMessage) Type() MessageType {
	return MessageTypeError
}

func (m ErrorMessage) encode(w io.Writer) error {
	err := writeMessageType(w, m.Type())
	if err != nil {
		return err
	}

	err = writeString(w, m.Message)
	if err != nil {
		return err
	}

	return nil
}

func (m *ErrorMessage) decode(r io.Reader) error {
	msg, err := readString(r)
	if err != nil {
		return err
	}

	m.Message = msg

	return nil
}
