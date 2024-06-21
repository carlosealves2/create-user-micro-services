package kafkaclient

import "testing"

func TestNewClient(t *testing.T) {
	options := KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "create-user",
		GroupID: "user-consumer",
	}
	instance, err := NewKafkaClient(&options)
	if err != nil {
		t.Error(err)
		return
	}

	err = instance.Publish(`{"id": 1}`)
	if err != nil {
		t.Error(err)
	}
}
