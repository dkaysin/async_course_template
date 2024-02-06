package eventbus

import (
	"context"
	"fmt"
	"time"
)

func ScheduleSendMessages(ew *EventWriter) {
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			value := fmt.Sprintf("test value %s", time.Now().Format(time.RFC3339))
			ew.TopicAWriter.WriteMessage(context.Background(), "test_key", value)
		}
	}()

}
