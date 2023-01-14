package db

import (
	"fmt"
	"os"

	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func DbConnect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("EMAIL_API_DB_USER"),
		Passwd:               os.Getenv("EMAIL_API_DB_PASS"),
		Net:                  "tcp",
		Addr:                 "ksadebiekorf.be",
		DBName:               "k16461ks_mailing",
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	return db, err
}

type EmailEntry struct {
	FirstName string
	LastName  string
	Email     string
}

func dbQueryEmail(db *sql.DB, query string) ([]EmailEntry, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []EmailEntry
	for rows.Next() {
		var fn string
		var ln string
		var eml string
		if err := rows.Scan(&fn, &ln, &eml); err != nil {
			return nil, err
		}

		entry := EmailEntry{fn, ln, eml}
		entries = append(entries, entry)
	}

	return entries, nil
}

func FetchBan(db *sql.DB, ban int) ([]EmailEntry, error) {
	query := "SELECT Persons.first_name, Persons.last_name, Emails.email FROM Persons " +
		"INNER JOIN Emails ON Persons.id = Emails.person_id " +
		"INNER JOIN Bannen ON Persons.id = Bannen.person_id " +
		fmt.Sprintf("WHERE Bannen.ban = %d;", ban)

	return dbQueryEmail(db, query)
}
