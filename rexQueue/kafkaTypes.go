package rexQueue

import (
	"github.com/IBM/sarama"
)

type EasyConsumerGroupHandler struct {
	readMsgFunc func(msg *sarama.ConsumerMessage)
}

func (h *EasyConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *EasyConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *EasyConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if h.readMsgFunc != nil {
			h.readMsgFunc(msg)
			continue
		}
		// 标记该消息为已处理，提交 offset
		session.MarkMessage(msg, "")
	}
	return nil
}
