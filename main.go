package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var server = "107.0.0.1"
var port = 3306
var user = "appuser"
var password = "Sofsof123"
var database = "passwords"

var db *sql.DB

const letterCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberCharset = "0123456789"
const symbolCharset = "!@#$%^&*()_+[]{}|;:,.<>?/~`"

func generatePassword(length int, useNumbers, useSymbols bool) string {
	charset := letterCharset
	if useNumbers {
		charset += numberCharset
	}
	if useSymbols {
		charset += symbolCharset
	}

	rand.Seed(time.Now().UnixNano())

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

func checkAndUpdateDatabase(password string) error {
	// Open database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database))
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the password already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM passwords WHERE password = ?", password).Scan(&count)
	if err != nil {
		return err
	}

	// If the password exists, generate a new one and recursively check again
	if count > 0 {
		newPassword := generatePassword(len(password), true, true)
		return checkAndUpdateDatabase(newPassword)
	}

	// If the password doesn't exist, insert it into the database
	_, err = db.Exec("INSERT INTO passwords (password) VALUES (?)", password)
	return err
}

func main() {
	var length int
	var useNumbers, useSymbols bool

	flag.IntVar(&length, "length", 12, "Length of the generated password")
	flag.BoolVar(&useNumbers, "numbers", false, "Include numbers in the generated password")
	flag.BoolVar(&useSymbols, "symbols", false, "Include symbols in the generated password")

	flag.Parse()

	password := generatePassword(length, useNumbers, useSymbols)

	// Check and update the database
	err := checkAndUpdateDatabase(password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated password: %s\n", password)
}
