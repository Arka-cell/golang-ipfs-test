package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	Tables "github.com/arka-cell/ummatest/tables"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	shell "github.com/ipfs/go-ipfs-api"
)

var db *sql.DB

func main() {
	// Define command-line flags
	user := flag.String("user", "my_user", "MySQL user")
	password := flag.String("password", "my_password", "MySQL password")
	database := flag.String("database", "my_database", "MySQL database")
	host := flag.String("host", "127.0.0.1", "MySQL host")
	port := flag.String("port", "3306", "MySQL port")
	table := flag.String("table", "messages", "Table Name")

	// Parse the command-line arguments
	flag.Parse()

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 *user,
		Passwd:               *password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", *host, *port),
		DBName:               *database,
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	sh := shell.NewShell("localhost:5001")
	fmt.Println(sh)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connecting to database %s with %s\n", *database, *user)

	columns, err := getColumnNamesAndTypes(*table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Column names and types are as follows:\n- %s \n", columns)
	messages, err := getMessages()
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(messages, "", "  ")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(jsonData))
	fmt.Printf(string(jsonData[0]))

	for i := 0; i < len(messages); i++ {

		ipfsHash, err := interactWithIPFS(string(jsonData[i]), sh)
		if err != nil {
			return
		}
		fmt.Printf("IPFS Data is: %s \n", string(jsonData[i]))
		fmt.Printf("IPFS Hash is: %s \n", ipfsHash)
	}

	if err != nil {
		log.Fatal(err)
	}

	// chainOfHashes := buildChainOfHashes(ipfsCID, columns)
	// fmt.Println("Chain of Hashes:", chainOfHashes)
}

func getColumnNamesAndTypes(tableName string) (map[string]string, error) {
	columns := make(map[string]string)

	rows, err := db.Query(fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s';", tableName))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, columnType string
		if err := rows.Scan(&columnName, &columnType); err != nil {
			return nil, err
		}
		columns[columnName] = columnType
	}

	return columns, nil
}

// getMessages retrieves all rows from the 'messages' table
func getMessages() ([]Tables.Message, error) {
	var messages []Tables.Message

	rows, err := db.Query("SELECT * FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Tables.Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Text, &message.CreatedAt, &message.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// interactWithIPFS interacts with IPFS to add data and obtain CID
func interactWithIPFS(data string, sh *shell.Shell) (string, error) {
	// Placeholder for IPFS interaction
	// You need to implement this based on your actual IPFS setup
	// This might involve using an IPFS library, API, or command-line interface
	// Return the obtained CID
	ipfsHash, err := sh.Add(strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}

	return ipfsHash, nil
}

// buildChainOfHashes builds a chain of hashes based on the IPFS CID and column names
func buildChainOfHashes(ipfsCID string, columns map[string]string) string {
	// Placeholder for building a chain of hashes
	// Customize this based on your actual requirements
	// You might use a hashing algorithm (e.g., SHA-256) on the IPFS CID and column names
	// to create a chain of hashes
	// Return the resulting chain of hashes
	return "YourChainOfHashesHere"
}
