package event_reader

import (
	global "async_course/main"
	"context"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	switch string(m.Key) {
	case global.EventAddUser:
		return er.handleAddUser(m)
	}
	return nil
}

func (er *EventReader) handleAddUser(m kafka.Message) error {
	payload, err := validatePayload[global.AddUserReq](m)
	if err != nil {
		return err
	}
	return er.s.AddUser(context.Background(), payload.UserId)
}
