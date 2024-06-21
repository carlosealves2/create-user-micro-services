package kafkaclient

import (
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type KafkaClient struct {
	config   *KafkaConfig
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
}

func NewKafkaClient(config *KafkaConfig) (*KafkaClient, error) {
	// Configuração do produtor
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.Brokers, producerConfig)
	if err != nil {
		return nil, err
	}

	// Configuração do consumidor
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, consumerConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		config:   config,
		producer: producer,
		consumer: consumer,
	}, nil
}

func (kc *KafkaClient) Publish(message string) error {
	msg := &sarama.ProducerMessage{
		Topic: kc.config.Topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := kc.producer.SendMessage(msg)
	return err
}

func (kc *KafkaClient) Consume(handler func(message *sarama.ConsumerMessage)) {
	consumer := kc.consumer
	ready := make(chan bool)

	go func() {
		for {
			if err := consumer.Consume(nil, []string{kc.config.Topic}, &consumerGroupHandler{handler: handler, ready: ready}); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
		}
	}()

	<-ready
	log.Println("Consumer up and running...")

	// Aguardando sinal para parar
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm

	if err := consumer.Close(); err != nil {
		log.Panicf("Error closing consumer: %v", err)
	}
}

type consumerGroupHandler struct {
	handler func(message *sarama.ConsumerMessage)
	ready   chan bool
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handler(message)
		session.MarkMessage(message, "")
	}
	return nil
}
