package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	Tables "github.com/arka-cell/ummatest/tables"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
)

var db *sql.DB

var IPNSPublishingKey string

func main() {
	user := flag.String("user", "my_user", "MySQL user")
	password := flag.String("password", "my_password", "MySQL password")
	database := flag.String("database", "my_database", "MySQL database")
	host := flag.String("host", "127.0.0.1", "MySQL host")
	port := flag.String("port", "3306", "MySQL port")
	table := flag.String("table", "messages", "Table Name")
	flag.Parse()
	cfg := mysql.Config{
		User:                 *user,
		Passwd:               *password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", *host, *port),
		DBName:               *database,
		AllowNativePasswords: true,
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	defer db.Close()
	if pingErr := db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connecting to database %s with %s\n", *database, *user)

	sh := shell.NewShell("127.0.0.1:5001")
	fmt.Println("Connecting to IPFS Shell")
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connecting to database %s with %s\n", *database, *user)
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	IPNSPublishingKey = os.Getenv("IPNS_PUBLISH_KEY")

	columns, err := getColumnNamesAndTypes(*table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Column names and types are as follows:\n- %s \n", columns)
	messages, err := getMessages()
	if err != nil {
		log.Fatal(err)
	}
	var ipfsHashes []string
	var ipnsHashes []string
	for i := range messages {
		message, err := json.Marshal(messages[i])
		if err != nil {
			fmt.Println(err)
			return
		}

		ipfsHash, err := getCID(string(message), sh)
		fmt.Printf("IPFS Hash for row with ID %x is: %s \n", messages[i].ID, ipfsHash)
		ipfsHashes = append(ipfsHashes, ipfsHash)
		err = addToIPNS(ipfsHash, sh)
		if err != nil {
			fmt.Println(err)
		}
		ipnsResolve, err := resolveIPNS(sh)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("IPNS Hash to edit row with ID %x is: %s \n", messages[i].ID, ipnsResolve)
		ipnsHashes = append(ipnsHashes, ipnsResolve)
		messages[i].IPNS = ipnsResolve
		messages[i].IPFS = ipfsHash
	}
	jsonData, err := json.MarshalIndent(messages, "", "  ")
	saveJSONToFile("json_result.json", string(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Get the column names and their types for the table message in MySQL
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

// getMessages retrieves all rows from the 'messages' table in MySQL
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
func getCID(data string, sh *shell.Shell) (string, error) {
	ipfsHash, err := sh.Add(strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}

	return ipfsHash, nil
}

func addToIPNS(ipfs string, sh *shell.Shell) error {
	var lifetime time.Duration = 50 * time.Hour
	var ttl time.Duration = 1 * time.Microsecond
	_, err := sh.PublishWithDetails(ipfs, IPNSPublishingKey, lifetime, ttl, true)

	return err
}

func resolveIPNS(sh *shell.Shell) (string, error) {
	return sh.Resolve(IPNSPublishingKey)
}

func saveJSONToFile(filename string, jsonData string) error {
	jsonBytes := []byte(jsonData)
	err := ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
