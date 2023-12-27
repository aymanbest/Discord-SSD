package storage

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)


var DB *sql.DB

func GetDB() *sql.DB {
    return DB
}

type FileTable struct {
	Files   map[string]string `json:"files"`
	IdChain map[string]string `json:"id_chain"`
}



func ConnectToDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := os.Getenv("PSQLINFO")
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
		files JSONB,
		id_chain JSONB
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
