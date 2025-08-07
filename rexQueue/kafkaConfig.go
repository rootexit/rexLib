package rexQueue

import (
	"github.com/IBM/sarama"
	"time"
)

type ConsumerMode string
type ProducerMode string

const (
	ModeProducerOnly      ConsumerMode = "producer"
	ModeConsumerGroup     ConsumerMode = "group"
	ModeSimpleConsumer    ConsumerMode = "simple"
	ModePartitionConsumer ConsumerMode = "partition"

	ProducerModeSync  ProducerMode = "sync"
	ProducerModeAsync ProducerMode = "async"
)

type KafkaConfig struct {
	ConsumerMode ConsumerMode
	ProducerMode ProducerMode
	Brokers      []string `json:",default=[localhost:29092]"`
	Topics       []string `json:",default=[]"`
	GroupId      string   `json:",default=default_group"`
	*sarama.Config
}

func Default(brokers []string, groupId string, topics []string) *KafkaConfig {
	c := sarama.NewConfig()
	c.Version = sarama.V4_0_0_0

	c.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最早开始
	c.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	c.Consumer.Offsets.AutoCommit.Enable = true
	c.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	c.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	c.Producer.Retry.Max = 5                    // 重试次数
	c.Producer.Return.Successes = true          // 返回成功确认

	c.Producer.Timeout = 5 * time.Second // 单条消息超时时间
	c.Net.WriteTimeout = 5 * time.Second // 网络层写超时
	c.Net.DialTimeout = 5 * time.Second
	c.Net.ReadTimeout = 5 * time.Second

	return &KafkaConfig{
		ConsumerMode: ModeConsumerGroup,
		ProducerMode: ProducerModeSync,
		Brokers:      brokers,
		Topics:       topics,
		GroupId:      groupId,
		Config:       c,
	}
}

func NewConfig(consumerMode ConsumerMode, producerMode ProducerMode, brokers []string, groupId string, topics []string, c *sarama.Config) *KafkaConfig {
	return &KafkaConfig{
		ConsumerMode: consumerMode,
		ProducerMode: producerMode,
		Brokers:      brokers,
		Topics:       topics,
		GroupId:      groupId,
		Config:       c,
	}
}

func (c *KafkaConfig) With(conf *sarama.Config) *KafkaConfig {
	c.Config = conf
	return c
}

func (c *KafkaConfig) WithConsumerMode(mode ConsumerMode) *KafkaConfig {
	c.ConsumerMode = mode
	return c
}

func (c *KafkaConfig) WithProducerMode(mode ProducerMode) *KafkaConfig {
	c.ProducerMode = mode
	return c
}
