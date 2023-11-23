package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BankAccount struct {
	ID            string  `json:"id"`
	AccountHolder string  `json:"account_holder"`
	AccountType   int     `json:"account_type"`
	Balance       float64 `json:"balance"`
}

type AcountRepository interface {
	Save(bankAccount BankAccount) error
	Delete(id string) error
	FindAll() (bankAccounts []BankAccount, err error)
	FindByID(id string) (bankAccount BankAccount, err error)
	Update(bankAccount BankAccount) error
}

type accountRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

func NewAccountRepository(db *mongo.Database) AcountRepository {
	collection := db.Collection("account")
	return accountRepository{db: db, col: collection}

}

func (r accountRepository) Save(bankAccount BankAccount) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.col.InsertOne(ctx, bankAccount)
	if err != nil {
		return err
	}

	return err
}

func (r accountRepository) Delete(id string) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := r.col.DeleteOne(context.Background(), filter)
	return err
}

func (r accountRepository) FindAll() (bankAccounts []BankAccount, err error) {
	_, err = r.col.Find(context.TODO(), &bankAccounts)
	return bankAccounts, err
}

func (r accountRepository) FindByID(id string) (bankAccount BankAccount, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err != nil {
		return bankAccount, err
	}

	filter := bson.D{{Key: "id", Value: id}}
	err = r.col.FindOne(ctx, filter).Decode(&bankAccount)
	return bankAccount, err
}

func (r accountRepository) Update(bankAccount BankAccount) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "id", Value: bankAccount.ID}}
	update := bson.D{{Key: "$set", Value: bankAccount}}

	_, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}
