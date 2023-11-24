package services_test

import (
	"gokafka/commands"
	"gokafka/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOpenAccount(t *testing.T) {
	command := commands.OpenAccountCommand{
		AccountHolder:  "name1",
		AccountType:    1,
		OpeningBalance: 1000,
	}

	//Arrange
	mockProducer := services.NewEventProducerMock()
	service := services.NewAccountServiceCommand(mockProducer)
	mockProducer.On("Produce", mock.AnythingOfType("OpenAccountEvent")).Return(nil)

	//Act
	id, err := service.OpenAccount(command) // Call the OpenAccount method

	//Assert
	assert.Nil(t, err)                 // Assert that the returned error is nil
	assert.NotEmpty(t, id)             // Assert that the returned ID is not empty
	mockProducer.AssertExpectations(t) // Assert that the expectations were met

	t.Run("invalid command", func(t *testing.T) {
		command := commands.OpenAccountCommand{
			AccountHolder:  "",
			AccountType:    0,
			OpeningBalance: 0,
		}

		//Act
		id, err := service.OpenAccount(command)

		//Assert
		assert.NotNil(t, err)                     // Assert that the returned error is not nil
		assert.Equal(t, fiber.ErrBadRequest, err) // Assert that the returned error is not nil
		assert.Empty(t, id)                       // Assert that the returned ID is empty
	})
}

func TestDepositFund(t *testing.T) {
	command := commands.DepositFundCommand{
		ID:     "id1",
		Amount: 1000,
	}

	//Arrange
	mockProducer := services.NewEventProducerMock()
	service := services.NewAccountServiceCommand(mockProducer)
	mockProducer.On("Produce", mock.AnythingOfType("DepositFundEvent")).Return(nil)

	//Act
	err := service.DepositFund(command)

	//Assert
	assert.Nil(t, err)                 // Assert that the returned error is nil
	mockProducer.AssertExpectations(t) // Assert that the expectations were met

	t.Run("invalid command", func(t *testing.T) {
		// Define the command to open an account
		command := commands.DepositFundCommand{
			ID:     "",
			Amount: 0,
		}

		// Call the OpenAccount method
		err := service.DepositFund(command)

		// Assert that the returned error is not nil
		assert.NotNil(t, err)

		assert.Equal(t, fiber.ErrBadRequest, err)
	})
}

func TestWithdrawFund(t *testing.T) {
	command := commands.WithdrawFundCommand{
		ID:     "id1",
		Amount: 1000,
	}

	//Arrange
	mockProducer := services.NewEventProducerMock()
	service := services.NewAccountServiceCommand(mockProducer)
	mockProducer.On("Produce", mock.AnythingOfType("WithdrawFundEvent")).Return(nil)

	//Act
	err := service.WithdrawFund(command)

	assert.Nil(t, err) // Assert that the returned error is nil

	mockProducer.AssertExpectations(t) // Assert that the expectations were met

	t.Run("invalid command", func(t *testing.T) {
		command := commands.WithdrawFundCommand{
			ID:     "",
			Amount: 0,
		}

		//Act
		err := service.WithdrawFund(command)

		//Assert
		assert.NotNil(t, err)                     // Assert that the returned error is not nil
		assert.Equal(t, fiber.ErrBadRequest, err) // Assert that the returned error is not nil
	})
}

func TestCloseAccount(t *testing.T) {
	command := commands.CloseAccountCommand{
		ID: "id1",
	}

	//Arrange
	mockProducer := services.NewEventProducerMock()
	mockProducer.On("Produce", mock.AnythingOfType("CloseAccountEvent")).Return(nil)
	service := services.NewAccountServiceCommand(mockProducer)

	//Act
	err := service.CloseAccount(command)

	//Assert
	assert.Nil(t, err)                 // Assert that the returned error is nil
	mockProducer.AssertExpectations(t) // Assert that the expectations were met

	t.Run("invalid command", func(t *testing.T) {
		command := commands.CloseAccountCommand{
			ID: "",
		}

		//Act
		err := service.CloseAccount(command)

		//Assert
		assert.NotNil(t, err)                     // Assert that the returned error is not nil
		assert.Equal(t, fiber.ErrBadRequest, err) // Assert that the returned error is not nil
	})
}

func BenchmarkOpenAccount(b *testing.B) {
	command := commands.OpenAccountCommand{
		AccountHolder:  "name1",
		AccountType:    1,
		OpeningBalance: 1000,
	}

	mockProducer := services.NewEventProducerMock()
	service := services.NewAccountServiceCommand(mockProducer)
	mockProducer.On("Produce", mock.AnythingOfType("OpenAccountEvent")).Return(nil)

	for i := 0; i < b.N; i++ {
		service.OpenAccount(command)
	}
}

func BenchmarkDepositFund(b *testing.B) {
	command := commands.DepositFundCommand{
		ID:     "id1",
		Amount: 1000,
	}

	mockProducer := services.NewEventProducerMock()
	service := services.NewAccountServiceCommand(mockProducer)
	mockProducer.On("Produce", mock.AnythingOfType("DepositFundEvent")).Return(nil)

	for i := 0; i < b.N; i++ {
		service.DepositFund(command)
	}
}
func BenchmarkWithdrawFund(b *testing.B) {
	command := commands.WithdrawFundCommand{
		ID:     "id1",
		Amount: 1000,
	}

	mockProducer := services.NewEventProducerMock()
	service := services.NewAccountServiceCommand(mockProducer)
	mockProducer.On("Produce", mock.AnythingOfType("WithdrawFundEvent")).Return(nil)

	for i := 0; i < b.N; i++ {
		service.WithdrawFund(command)
	}
}

func BenchmarkCloseAccount(b *testing.B) {
	command := commands.CloseAccountCommand{
		ID: "id1",
	}

	mockProducer := services.NewEventProducerMock()
	mockProducer.On("Produce", mock.AnythingOfType("CloseAccountEvent")).Return(nil)
	service := services.NewAccountServiceCommand(mockProducer)

	for i := 0; i < b.N; i++ {
		service.CloseAccount(command)
	}
}
