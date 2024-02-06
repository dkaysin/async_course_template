package eventbus

import (
	"context"
	"fmt"
	"time"
)

func ScheduleSendMessages(ew *EventWriter) {
	tickerA := time.NewTicker(time.Second * 2)
	go func() {
		for range tickerA.C {
			value := fmt.Sprintf("test value %s", time.Now().Format(time.RFC3339))
			ew.TopicAWriter.WriteMessage(context.Background(), "key_for_a", value)
		}
	}()

	tickerB := time.NewTicker(time.Second * 3)
	go func() {
		for range tickerB.C {
			value := fmt.Sprintf("test value %s", time.Now().Format(time.RFC3339))
			ew.TopicBWriter.WriteMessage(context.Background(), "key_for_b", value)
		}
	}()

}
