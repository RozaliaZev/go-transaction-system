package db

import (
	"context"
	"fmt"
	ts "go-tr-syst/pkg/gen-tr-syst"
	"log"

	"os"
	"github.com/jackc/pgx/v4"
)


var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var sslmode = os.Getenv("SSLMODE")
var dbname = os.Getenv("DBNAME")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

//var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", "localhost", "5432", "postgres", "Password10", "base", "disable")

func ConnectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}
	return conn, nil
}

func CreateTable() error {
	conn, err := ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `
  CREATE TABLE IF NOT EXISTS transactions (
   id SERIAL PRIMARY KEY,
   amount NUMERIC(10, 2) NOT NULL,
   currency VARCHAR(3) NOT NULL,
   sender_card_number VARCHAR(16) NOT NULL,
   recipient_phone_number VARCHAR(11) NOT NULL,
   date_time TIMESTAMP NOT NULL DEFAULT NOW()
  );
 `)
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
		return err
	}

	return nil
}

func WriteTransactionData(req *ts.CreateTransactionRequest) (int64, error) {
	conn, err := ConnectDB()
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())

	amount := req.GetAmount()
	cerrency := req.GetCurrency()
	senderCardNumber := req.GetSenderCardNumber()
	recipientPhoneNumber := req.GetRecipientPhoneNumber()

	var id int64
	err = conn.QueryRow(context.Background(), `
  INSERT INTO transactions (amount, currency, sender_card_number, recipient_phone_number)
  VALUES ($1, $2, $3, $4)
  RETURNING id
 `, amount, cerrency, senderCardNumber, recipientPhoneNumber).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to insert data into table: %v\n", err)
		return 0, err
	}

	return id, nil
}

type Transaction struct {
	ID                   int64   `json:"id"`
	Amount               string `json:"amount,omitempty"`
	Currency             string  `json:"currency,omitempty"`
	SenderCardNumber     string  `json:"sender_card_number,omitempty"`
	RecipientPhoneNumber string  `json:"recipient_phone_number,omitempty"`
	DateTime             string  `json:"date_time,omitempty"`
}

func GetTransactionByID(req *ts.GetTransactionRequest) (Transaction, error) {
	conn, err := ConnectDB()
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return Transaction{}, err
	}
	defer conn.Close(context.Background())

	id := req.Id
	row := conn.QueryRow(context.Background(), `
	SELECT id, amount, currency, sender_card_number, recipient_phone_number, date_time
	FROM transactions
	WHERE id = $1
	LIMIT 1
   `, id)

	var transaction Transaction
	err = row.Scan(&transaction.ID, &transaction.Amount, &transaction.Currency, &transaction.SenderCardNumber, &transaction.RecipientPhoneNumber, &transaction.DateTime)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

