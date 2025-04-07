package learn_go_db_test

import (
	"context"
	"database/sql"
	"fmt"
	dbconn "learn-go-db/db"
	"strconv"
	"testing"
)

func TestEmpty(t *testing.T) {

}

func TestOpenConnection(t *testing.T) {
	db := dbconn.GetConnection()
	defer db.Close()

	ctx := context.Background()
	query := "INSERT INTO customer(id, name) values ('1', 'alwi')"

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
}

func TestSelect(t *testing.T) {
	db := dbconn.GetConnection()
	defer db.Close()

	ctx := context.Background()
	query := "select id, name, email, balance, rating, created_at, birth_date, married from customer"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		/*
			e.g. sql.NullString

			type NullString struct {
				String string
				Valid bool --> true if string is not null
			}
		*/

		var id, name, email sql.NullString // for nullable column
		var balance sql.NullInt64
		var rating float64
		var created_at, birth_date sql.NullTime
		var married sql.NullBool

		err := rows.Scan(&id, &name, &email, &balance, &rating, &created_at, &birth_date, &married)
		if err != nil {
			panic(err)
		}

		fmt.Println(id, name, email, balance, rating, created_at, birth_date, married)
	}

	defer rows.Close()
}

func TestSqlParam(t *testing.T) {
	db := dbconn.GetConnection()
	defer db.Close()

	username := "alwi"
	password := "alwi"

	ctx := context.Background()
	//query := "select username from public.user where username = '" + username + "' and password = '" + password + "' limit 1"
	query := "select username from public.user where username = $1 and password = $2 limit 1"
	fmt.Println(query)

	rows, err := db.QueryContext(ctx, query, username, password)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string

		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("logged as", username)
	} else {
		fmt.Println("Failed")
	}
}

// Result.LastInsertId() does not support PostgreSQL driver
func TestAutoIncrement(t *testing.T) {
	db := dbconn.GetConnection()
	defer db.Close()

	comment := "test"
	email := "alwi@email.com"

	ctx := context.Background()
	query := "INSERT INTO comments(email,  comment) values ($1, $2) returning id"

	var lastInsertId int
	res := db.QueryRowContext(ctx, query, comment, email)

	if res.Err() != nil {
		panic(res.Err())
	}

	res.Scan(&lastInsertId)
	fmt.Println(lastInsertId)
}

func TestPrepareStatement(t *testing.T) {
	db := dbconn.GetConnection()
	defer db.Close()

	ctx := context.Background()
	query := "INSERT INTO comments(email,  comment) values ($1, $2) returning id"

	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "alwi" + strconv.Itoa(i) + "email.com"
		comment := "comment ke " + strconv.Itoa(i)

		res, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("res", lastId)
	}
}

func TestTransaction(t *testing.T) {
	db := dbconn.GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email,  comment) values ($1, $2) returning id"

	for i := 0; i < 10; i++ {
		email := "alwi" + strconv.Itoa(i) + "@email.com"
		comment := "comment ke " + strconv.Itoa(i)

		res, err := tx.ExecContext(ctx, query, email, comment)
		if err != nil {
			panic(err)
		}

		// lastId, err := res.LastInsertId()
		// if err != nil {
		// 	panic(err)
		// }

		fmt.Println("res", res)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
