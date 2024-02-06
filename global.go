package global

// kafka
const (
	KafkaConsumerGroupID = "my-consumer-group-id"

	KafkaTopicIDA = "topic-A"
	KafkaTopicIDB = "topic-B"
)

// events
const (
	EventAddUser = "add_user"
)

type AddUserReq struct {
	UserId string `json:"user_id" validate:"required"`
}
