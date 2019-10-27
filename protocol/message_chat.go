package protocol

import (
	"io"
)

// ChatMessage represents a message sent to the chat
type ChatMessage struct {
	UserName string
	Message  string
}

// Type returns the message type
func (m ChatMessage) Type() MessageType {
	return MessageTypeChat
}

func (m ChatMessage) encode(w io.Writer) error {
	err := writeMessageType(w, m.Type())
	if err != nil {
		return err
	}

	err = writeString(w, m.UserName)
	if err != nil {
		return err
	}

	err = writeString(w, m.Message)
	if err != nil {
		return err
	}

	return nil
}

func (m *ChatMessage) decode(r io.Reader) error {
	userName, err := readString(r)
	if err != nil {
		return err
	}

	msg, err := readString(r)
	if err != nil {
		return err
	}

	m.UserName = userName
	m.Message = msg

	return nil
}
