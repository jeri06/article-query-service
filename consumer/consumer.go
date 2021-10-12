package consumer

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	eventHandler EventHandler
}

func NewConsumerHandler(eventHandler EventHandler) Consumer {
	return Consumer{eventHandler: eventHandler}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		consumer.claim(context.Background(), message)
		session.MarkMessage(message, "")

	}

	return nil
}

func (consumer *Consumer) claim(ctx context.Context, message *sarama.ConsumerMessage) {

	if err := consumer.eventHandler.Handle(ctx, message); err != nil {
		return
	}
	return
}

type EventHandler interface {
	Handle(ctx context.Context, message interface{}) (err error)
}
