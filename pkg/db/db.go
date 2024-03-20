package db

import (
	"database/sql"
	"fmt"
	"log"
)

// InitDB Initiates database file with table subdomaains if not exist
func InitDB(dbFileName string) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS subdomains (id INTEGER PRIMARY KEY, subdomain TEXT)")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("[+] Database initialized successfully.")
	}
}

// AddSubdomains is a function to add subdomaains to a sqlite database fole
func AddSubdmomains(subdomains []string, dbFileName string) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO subdomains (subdomain) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, subdomain := range subdomains {
		_, err = stmt.Exec(subdomain)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] Added %d new subdomains to the database.\n", len(subdomains))
}

// DBsubdomains is Function to get a slice of subdomains from database file
func Getsubdomains(dbFileName string) []string {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT subdomain FROM subdomains")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var subdomains []string

	for rows.Next() {
		var subdomain string
		err := rows.Scan(&subdomain)
		if err != nil {
			log.Fatal(err)
		}
		subdomains = append(subdomains, subdomain)
	}

	fmt.Println("[+] Retrieved subdomains from the database.")

	return subdomains
}
