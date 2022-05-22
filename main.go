package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id           uint32
	Author_id    uint32
	Title        string
	Publish_date time.Time
	Author       Author
}

type Author struct {
	Id   uint32
	Name string
}

func main() {
	db, err := createConnection()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	err = queryBooks(ctx, db)

	if err != nil {
		panic(err)
	}

	err = addBook(ctx, db, "Go Programming", 1, time.Now())

	if err != nil {
		panic(err)
	}

	err = queryBooks(ctx, db)

	if err != nil {
		panic(err)
	}

	db.Close()
}

func createConnection() (*sql.DB, error) {
	connectionString := "root:123456@tcp(localhost:3306)/go_database_connection?parseTime=True"

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func queryBooks(ctx context.Context, db *sql.DB) error {
	qry := `
	SELECT b.id, b.author_id, b.title, b.publish_data, 
	a.id, a.name
	FROM books AS b
	JOIN autors AS a ON b.author_id = a.id

	`

	rows, err := db.QueryContext(ctx, qry)
	if err != nil {
		return err
	}

	books := []Book{}

	for rows.Next() {
		b := Book{}

		err = rows.Scan(&b.Id, &b.Author_id, &b.Title, &b.Publish_date, &b.Author.Id, &b.Author.Name)
		if err != nil {
			return err
		}
		books = append(books, b)
	}

	fmt.Println(books)

	return nil
}

func addBook(ctx context.Context, db *sql.DB, title string, author_id uint32, publish_data time.Time) error {
	qryadd := `
	INSERT INTO books (title, author_id, publish_data) values (?,?,?)
	`

	result, err := db.ExecContext(ctx, qryadd, title, author_id, publish_data)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println("Book added with id: ", id)

	return nil
}
