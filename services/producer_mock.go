package services

import (
	"gokafka/events"

	"github.com/stretchr/testify/mock"
)

type eventProducerMock struct {
	mock.Mock
}

func NewEventProducerMock() *eventProducerMock {
	return &eventProducerMock{}
}

func (obj *eventProducerMock) Produce(event events.Event) error {
	args := obj.Called(event)
	return args.Error(0)
}
