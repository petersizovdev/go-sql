package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type user struct {
	ID int64
	Name string
	Email string
	Password string
	RegisteredAt time.Time
}

func main () {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	var (
		DB_HOST     = os.Getenv("DB_HOST")
		DB_PORT     = os.Getenv("DB_PORT")
		DB_USER     = os.Getenv("DB_USER")
		DB_NAME     = os.Getenv("DB_NAME")
	)

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
	DB_HOST, DB_PORT, DB_USER, DB_NAME)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = insertUser(db, user{
		Name:"Name",
		Email:"a@a.com",
		Password: "sdjfkdfsjfvcxvikx",
	})

	users, err := getUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

	err = updateUser(db, 4, user{
		Name: "Anton",
		Email: "e@e.com",
	})

	users, err = getUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)
}

func getUsers(db *sql.DB) ([]user, error){
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	users := make([]user, 0)
	for rows.Next() {
		u := user{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}


func getUserById(db *sql.DB, id int) (user, error){
	var u user
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", 2).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
	return u, err
}

func insertUser(db *sql.DB, u user) error {
	tx, err := db.Begin() 
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (name, email, password) values ($1, $2, $3)", u.Name, u.Email, u.Password)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO logs (entity, action) values ($1, $2)", "user", "created")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func updateUser(db *sql.DB, id int, u user) error {
	_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
	return err
}

func deleteUser(db *sql.DB, id int) error{
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

