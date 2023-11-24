package services_test

import (
	"errors"
	"gokafka/events"
	"gokafka/services"
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

type MarshalFail struct{}

func (m MarshalFail) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal error")
}
func TestNewEventProducer(t *testing.T) {
	//Arrange
	mockProducer := mocks.NewSyncProducer(t, nil)

	//Act
	producer := services.NewEventProducer(mockProducer)

	//Assert
	assert.NotNil(t, producer)
}

func TestProduceSuccess(t *testing.T) {
	topics := events.Topics
	mockProducer := mocks.NewSyncProducer(t, nil)
	producer := services.NewEventProducer(mockProducer)

	for _, topic := range topics {
		mockProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
			return nil
		})

		err := producer.Produce(topic)
		assert.Nil(t, err)
	}

}

func TestProduceFailure(t *testing.T) {
	topics := events.Topics
	mockProducer := mocks.NewSyncProducer(t, nil)
	producer := services.NewEventProducer(mockProducer)

	for _, topic := range topics {
		mockProducer.ExpectSendMessageAndFail(errors.New("expected error"))

		err := producer.Produce(topic)
		assert.NotNil(t, err)
	}

	t.Run("invalid event", func(t *testing.T) {

		mockProducer := mocks.NewSyncProducer(t, nil)
		producer := services.NewEventProducer(mockProducer)

		event := MarshalFail{}

		err := producer.Produce(event)
		assert.NotNil(t, err)
		assert.Equal(t, "marshal error", err.Error())
	})
}

func BenchmarkNewEventProducer(b *testing.B) {
	mockProducer := sarama.SyncProducer(nil)

	for i := 0; i < b.N; i++ {
		services.NewEventProducer(mockProducer)
	}
}
