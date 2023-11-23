package services

import (
	"encoding/json"
	"fmt"
	"gokafka/events"
	"gokafka/repositories"
	"log"
	"reflect"
)

type EventHandler interface {
	Handle(topic string, eventBytes []byte)
}

type accountEventHandler struct {
	accountRepo repositories.AcountRepository
}

func NewAccountHandler(accountRepo repositories.AcountRepository) EventHandler {
	return accountEventHandler{accountRepo}

}

func (obj accountEventHandler) Handle(topic string, eventBytes []byte) {
	switch topic {
	case reflect.TypeOf(events.OpenAccountEvent{}).Name():
		event := &events.OpenAccountEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			fmt.Println("error unmarshal", err)
			return
		}
		fmt.Println("event save", event)
		bankAccount := repositories.BankAccount{
			ID:            event.ID,
			AccountHolder: event.AccountHolder,
			AccountType:   event.AccountType,
			Balance:       event.OpeningBalance,
		}

		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			fmt.Println("error save", err)
			return
		}
		log.Println("account saved")
	case reflect.TypeOf(events.DepositFundEvent{}).Name():
		event := &events.DepositFundEvent{}
		err := json.Unmarshal(eventBytes, event)

		if err != nil {
			fmt.Println("error unmarshal", err)
			return
		}
		bankAccount, err := obj.accountRepo.FindByID(event.ID)

		if err != nil {
			fmt.Println("error find", err)
			return
		}
		bankAccount.Balance += event.Amount
		err = obj.accountRepo.Update(bankAccount)
		if err != nil {
			fmt.Println("error update", err)
			return
		}
		log.Println("account updated", bankAccount)
	case reflect.TypeOf(events.WithdrawFundEvent{}).Name():
		event := &events.WithdrawFundEvent{}
		err := json.Unmarshal(eventBytes, event)

		if err != nil {
			fmt.Println("error unmarshal", err)
			return
		}
		bankAccount, err := obj.accountRepo.FindByID(event.ID)
		if err != nil {
			fmt.Println("error find", err)
			return
		}
		bankAccount.Balance -= event.Amount
		err = obj.accountRepo.Update(bankAccount)
		if err != nil {
			fmt.Println("error update", err)
			return
		}
		log.Println("account updated", bankAccount)
	case reflect.TypeOf(events.CloseAccountEvent{}).Name():
		event := &events.CloseAccountEvent{}
		err := json.Unmarshal(eventBytes, event)
		if err != nil {
			fmt.Println("error unmarshal", err)
		}
		err = obj.accountRepo.Delete(event.ID)
		if err != nil {
			fmt.Println("error delete", err)
		}
		log.Println("account deleted")
	default:
		fmt.Println("unknow topic", topic)
	}
}
