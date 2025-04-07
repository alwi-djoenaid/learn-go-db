package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	dbconn "learn-go-db/db"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Customer struct {
	Id         string         `db:"id"`
	Name       sql.NullString `db:"name"`
	Email      sql.NullString `db:"email"`
	Balance    sql.NullInt64  `db:"balance"`
	Rating     float64        `db:"rating"`
	Created_at sql.NullTime   `db:"created_at"`
	Birth_date sql.NullTime   `db:"birth_date"`
	Married    sql.NullBool   `db:"married"`
}

func GetData() [][]string {
	result := [][]string{{"id", "name", "email", "balance", "rating", "created_at", "birth_date", "married"}}
	c := []Customer{}
	db := dbconn.GetxConnection()
	defer db.Close()

	ctx := context.Background()
	query := "select id, name, email, balance, rating, created_at, birth_date, married from customer"

	err := db.SelectContext(ctx, &c, query)
	if err != nil {
		panic(err)
	}

	for _, cust := range c {
		cusStr := []string{
			cust.Id,
			cust.Name.String,
			cust.Email.String,
			strconv.FormatInt(cust.Balance.Int64, 10),
			strconv.FormatFloat(cust.Rating, 'E', -1, 64),
			cust.Created_at.Time.String(),
			cust.Birth_date.Time.String(),
			strconv.FormatBool(cust.Married.Bool),
		}

		result = append(result, cusStr)
	}

	return result
}

func main() {
	customer := GetData()

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Path:", path)

	csvFile, err := os.Create("customers.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	w := csv.NewWriter(csvFile)
	w.WriteAll(customer)

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	defer csvFile.Close()

	if err := csvFile.Close(); err != nil {
		panic(err)
	}

	pathFile := filepath.Join(path, csvFile.Name())
	fmt.Println(pathFile)

	time.Sleep(2 * time.Second)
	if err := os.Remove(pathFile); err != nil {
		panic(err)
	}
}
