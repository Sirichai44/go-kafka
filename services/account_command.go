package services

import (
	"gokafka/commands"
	"gokafka/events"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AccountService interface {
	OpenAccount(command commands.OpenAccountCommand) (id string, err error)
	DepositFund(command commands.DepositFundCommand) error
	WithdrawFund(command commands.WithdrawFundCommand) error
	CloseAccount(command commands.CloseAccountCommand) error
}

type accountService struct {
	eventProducer EventProducer
}

func NewAccountServiceCommand(eventProducer EventProducer) AccountService {
	return accountService{eventProducer}
}

func (obj accountService) OpenAccount(command commands.OpenAccountCommand) (id string, err error) {

	if command.AccountHolder == "" || command.AccountType == 0 || command.OpeningBalance == 0 {
		return "", fiber.ErrBadRequest
	}
	event := events.OpenAccountEvent{
		ID:             uuid.NewString(),
		AccountHolder:  command.AccountHolder,
		AccountType:    command.AccountType,
		OpeningBalance: command.OpeningBalance,
	}

	log.Printf("%#v", event)
	return event.ID, obj.eventProducer.Produce(event)
}

func (obj accountService) DepositFund(command commands.DepositFundCommand) error {
	if command.ID == "" || command.Amount == 0 {
		return fiber.ErrBadRequest
	}
	event := events.DepositFundEvent{
		ID:     command.ID,
		Amount: command.Amount,
	}

	log.Printf("%#v", event)
	return obj.eventProducer.Produce(event)
}

func (obj accountService) WithdrawFund(command commands.WithdrawFundCommand) error {
	if command.ID == "" || command.Amount == 0 {
		return fiber.ErrBadRequest
	}
	event := events.WithdrawFundEvent{
		ID:     command.ID,
		Amount: command.Amount,
	}

	log.Printf("%#v", event)
	return obj.eventProducer.Produce(event)
}

func (obj accountService) CloseAccount(command commands.CloseAccountCommand) error {
	if command.ID == "" {
		return fiber.ErrBadRequest
	}
	event := events.CloseAccountEvent{
		ID: command.ID,
	}

	log.Printf("%#v", event)
	return obj.eventProducer.Produce(event)
}
