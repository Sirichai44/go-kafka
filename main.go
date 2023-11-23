package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gokafka/controller"
	"gokafka/events"
	"gokafka/repositories"
	"gokafka/services"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func initDatabase() *mongo.Database {
	username := url.QueryEscape(viper.GetString("db.username"))
	password := url.QueryEscape(viper.GetString("db.password"))

	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		username,
		password,
		viper.GetString("db.host"),
		viper.GetString("db.port"),
	)
	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client.Database(viper.GetString("db.database"))
}

func main() {
	// Initialize the database
	db := initDatabase()
	accountRepo := repositories.NewAccountRepository(db)
	accountEventHandler := services.NewAccountHandler(accountRepo)
	accountConsumerHandler := services.NewConsumerHandler(accountEventHandler)

	// Create the consumer
	consumer, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.servers"), viper.GetString("kafka.group"), nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Create the producer
	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	eventProducer := services.NewEventProducer(producer)
	accountService := services.NewAccountServiceCommand(eventProducer)
	accountController := controller.NewAccountController(accountService)

	// Create the Fiber app
	app := fiber.New()

	app.Post("/account", accountController.OpenAccount)
	app.Post("/deposit", accountController.DepositFund)
	app.Post("/withdraw", accountController.WithdrawFund)
	app.Post("/close", accountController.CloseAccount)

	// Run the Fiber app in a separate goroutine
	go func() {
		log.Println("Fiber app start...")
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Error starting Fiber app: %v", err)
		}
	}()

	// Run the consumer in a separate goroutine
	go func() {
		log.Println("AccountConsumer start...")
		for {
			err := consumer.Consume(context.Background(), events.Topics, accountConsumerHandler)
			if err != nil {
				log.Printf("Error consuming messages: %v", err)
			}
		}
	}()

	// Run the producer in a separate goroutine
	go func() {
		log.Println("AccountProducer start...")
		for {
			// Replace this with your producer logic
			msg := &sarama.ProducerMessage{
				Topic: "your-topic",
				Value: sarama.StringEncoder("your-message"),
			}
			_, _, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Prevent the main function from exiting immediately
	select {}
}
