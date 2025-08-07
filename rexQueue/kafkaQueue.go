package rexQueue

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logc"
)

type (
	Queue interface {
		WithAsyncProducerErrFunc(asyncProducerErrFunc func(err error)) *defaultQueue
		GetASyncProducer() sarama.AsyncProducer
		GetSyncProducer() sarama.SyncProducer
		GetConsumer() sarama.Consumer
		GetConsumerGroup() sarama.ConsumerGroup
		Close() error
		CatchAsyncErr(asyncProducerErrFunc func(err error))
		Consume(ctx context.Context, consumerGroup sarama.ConsumerGroup, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler)
		EasyConsume(ctx context.Context, consumerGroup sarama.ConsumerGroup, topics []string, readMsgFunc func(msg *sarama.ConsumerMessage))
		SyncSendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
		SyncSendMessageCtx(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
		EasySyncSendMessage(topic, key, value string) (partition int32, offset int64, err error)
		EasySyncSendMessageCtx(ctx context.Context, topic, key, value string) (partition int32, offset int64, err error)
		AsyncSendMessage(msg *sarama.ProducerMessage) (err error)
		AsyncSendMessageCtx(ctx context.Context, msg *sarama.ProducerMessage) (err error)
		EasyAsyncSendMessage(topic, key, value string) (err error)
		EasyAsyncSendMessageCtx(ctx context.Context, topic, key, value string) (err error)
	}
	defaultQueue struct {
		conf                 *KafkaConfig
		consumerMode         ConsumerMode
		producerMode         ProducerMode
		syncProducer         sarama.SyncProducer
		asyncProducer        sarama.AsyncProducer
		asyncProducerErrFunc func(err error)
		consumer             sarama.Consumer
		consumerGroup        sarama.ConsumerGroup
	}
)

func NewQueue(conf *KafkaConfig) (*defaultQueue, error) {
	q := defaultQueue{
		consumerMode: conf.ConsumerMode,
		producerMode: conf.ProducerMode,
		conf:         conf,
	}

	// 初始化 consumer
	switch conf.ProducerMode {
	case ProducerModeSync:
		// 创建 producer
		producer, err := sarama.NewSyncProducer(conf.Brokers, conf.Config)
		if err != nil {
			return nil, fmt.Errorf("producer 初始化失败: %w", err)
		}
		q.syncProducer = producer
	case ProducerModeAsync:
		// 创建 producer
		producer, err := sarama.NewAsyncProducer(conf.Brokers, conf.Config)
		if err != nil {
			return nil, fmt.Errorf("producer 初始化失败: %w", err)
		}
		q.asyncProducer = producer
	default:
		return nil, fmt.Errorf("未知 consumer 类型: %s", conf.ConsumerMode)
	}
	// 初始化 consumer
	switch conf.ConsumerMode {
	case ModeConsumerGroup:
		cg, err := sarama.NewConsumerGroup(conf.Brokers, conf.GroupId, conf.Config)
		if err != nil {
			return nil, fmt.Errorf("consumer group 初始化失败: %w", err)
		}
		q.consumerGroup = cg
	case ModeSimpleConsumer:
		c, err := sarama.NewConsumer(conf.Brokers, conf.Config)
		if err != nil {
			return nil, fmt.Errorf("consumer 初始化失败: %w", err)
		}
		q.consumer = c
	case ModeProducerOnly:
		// do nothing
	default:
		return nil, fmt.Errorf("未知 consumer 类型: %s", conf.ConsumerMode)
	}
	return &q, nil
}

func (q *defaultQueue) WithAsyncProducerErrFunc(asyncProducerErrFunc func(err error)) *defaultQueue {
	q.asyncProducerErrFunc = asyncProducerErrFunc
	return q
}

func (q *defaultQueue) GetASyncProducer() sarama.AsyncProducer {
	return q.asyncProducer
}

func (q *defaultQueue) GetSyncProducer() sarama.SyncProducer {
	return q.syncProducer
}

func (q *defaultQueue) GetConsumer() sarama.Consumer {
	return q.consumer
}

func (q *defaultQueue) GetConsumerGroup() sarama.ConsumerGroup {
	return q.consumerGroup
}

func (q *defaultQueue) Close() error {
	if q.consumer != nil {
		if err := q.consumer.Close(); err != nil {
			return fmt.Errorf("关闭 consumer 失败: %w", err)
		}
	}
	if q.consumerGroup != nil {
		if err := q.consumerGroup.Close(); err != nil {
			return fmt.Errorf("关闭 consumer group 失败: %w", err)
		}
	}
	if q.syncProducer != nil {
		if err := q.syncProducer.Close(); err != nil {
			return err
		}
	}
	if q.asyncProducer != nil {
		if err := q.asyncProducer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (q *defaultQueue) CatchAsyncErr(asyncProducerErrFunc func(err error)) {
	q.asyncProducerErrFunc = asyncProducerErrFunc
	go func() {
		for err := range q.asyncProducer.Errors() {
			if q.asyncProducerErrFunc != nil {
				q.asyncProducerErrFunc(err)
			}
		}
	}()
}

func (q *defaultQueue) Consume(ctx context.Context, consumerGroup sarama.ConsumerGroup, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler) {
	for {
		if err := consumerGroup.Consume(ctx, topics, consumerGroupHandler); err != nil {
			logc.Errorf(ctx, "Error consuming: %v", err)
			break
		}
		if ctx.Err() != nil {
			logc.Errorf(ctx, "ctx exit: %v", ctx.Err())
			break
		}
	}
}

func (q *defaultQueue) EasyConsume(ctx context.Context, consumerGroup sarama.ConsumerGroup, topics []string, readMsgFunc func(msg *sarama.ConsumerMessage)) {
	for {
		if err := consumerGroup.Consume(ctx, topics, &EasyConsumerGroupHandler{
			readMsgFunc: readMsgFunc,
		}); err != nil {
			logc.Errorf(ctx, "Error consuming: %v", err)
			break
		}
		if ctx.Err() != nil {
			logc.Errorf(ctx, "ctx exit: %v", ctx.Err())
			break
		}
	}
}

func (q *defaultQueue) SyncSendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if q.syncProducer == nil {
		return 0, 0, fmt.Errorf("sync producer 未初始化")
	}
	partition, offset, err = q.syncProducer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("send msg failed: %w", err)
	}
	//fmt.Printf("消息发送成功，分区: %d, 偏移量: %d\n", partition, offset)
	return partition, offset, nil
}

func (q *defaultQueue) SyncSendMessageCtx(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	done := make(chan error, 1)

	go func() {
		partition, offset, err = q.SyncSendMessage(msg)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return 0, 0, ctx.Err() // context 超时 or 取消
	case err = <-done:
		return partition, offset, err
	}
}

func (q *defaultQueue) EasySyncSendMessage(topic, key, value string) (partition int32, offset int64, err error) {
	return q.SyncSendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	})
}

func (q *defaultQueue) EasySyncSendMessageCtx(ctx context.Context, topic, key, value string) (partition int32, offset int64, err error) {
	done := make(chan error, 1)

	go func() {
		partition, offset, err = q.SyncSendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(value),
		})
		done <- err
	}()

	select {
	case <-ctx.Done():
		return 0, 0, ctx.Err() // context 超时 or 取消
	case err = <-done:
		return partition, offset, err
	}
}

func (q *defaultQueue) AsyncSendMessage(msg *sarama.ProducerMessage) (err error) {
	if q.asyncProducer == nil {
		return fmt.Errorf("async producer 未初始化")
	}
	q.asyncProducer.Input() <- msg
	return nil
}

func (q *defaultQueue) AsyncSendMessageCtx(ctx context.Context, msg *sarama.ProducerMessage) (err error) {
	if q.asyncProducer == nil {
		return fmt.Errorf("async producer 未初始化")
	}
	select {
	case q.asyncProducer.Input() <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *defaultQueue) EasyAsyncSendMessage(topic, key, value string) (err error) {
	if q.asyncProducer == nil {
		return fmt.Errorf("async producer 未初始化")
	}
	q.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
	return nil
}

func (q *defaultQueue) EasyAsyncSendMessageCtx(ctx context.Context, topic, key, value string) (err error) {
	if q.asyncProducer == nil {
		return fmt.Errorf("async producer 未初始化")
	}
	select {
	case q.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
