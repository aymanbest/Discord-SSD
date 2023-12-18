package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


var DB *sql.DB

func GetDB() *sql.DB {
    return DB
}

type FileTable struct {
	Files   string `db:"files"`
	IdChain string `db:"id_chain"`
}



func ConnectToDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    psqlInfo := os.Getenv("DATABASE_URL")
    if psqlInfo == "" {
        log.Fatal("DATABASE_URL not set")
    }
    DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    err = DB.Ping()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Successfully connected!")
}

func InitTable() {
	query := `CREATE TABLE IF NOT EXISTS discssd (
		files TEXT,
		id_chain TEXT
	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FileTable) AddFile(names string, id string) {
	query := `INSERT INTO discssd (files, id_chain) VALUES ($1, $2)`
	_, err := DB.Exec(query, names, id)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FileTable) AddToChain(parentId string, childId string) {
	query := `UPDATE discssd SET id_chain = $1 WHERE files = $2`
	_, err := DB.Exec(query, childId, parentId)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FileTable) DeleteFile(name string) {
	query := `DELETE FROM discssd WHERE files = $1`
	_, err := DB.Exec(query, name)
	if err != nil {
		log.Fatal(err)
	}
}