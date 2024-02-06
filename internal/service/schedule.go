package service

import (
	global "async_course/main"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func (s *Service) ScheduleSendMessages() {
	tickerA := time.NewTicker(time.Second * 2)
	go func() {
		for range tickerA.C {
			req := global.AddUserReq{UserId: fmt.Sprint(rand.Intn(100))}
			s.ew.TopicAWriter.WriteJSON(context.Background(), global.EventAddUser, req)
		}
	}()

	tickerB := time.NewTicker(time.Second * 3)
	go func() {
		for range tickerB.C {
			value := fmt.Sprintf("test value %s", time.Now().Format(time.RFC3339))
			s.ew.TopicBWriter.WriteString(context.Background(), "key_for_b", value)
		}
	}()

}
