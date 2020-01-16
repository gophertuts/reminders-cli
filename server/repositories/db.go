package repositories

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gophertuts/reminders-cli/server/models"
)

// dbConfig represents the config which is used when DB is initialized
type dbConfig struct {
	ID int `json:"id"`
}

// DB represents the application server database (json file)
type DB struct {
	DB       *os.File
	DBCfg    *os.File
	ID       int
	Checksum string
}

// NewDB creates a new instance of application file DB
func NewDB(dbPath, dbCfgPath string) *DB {
	dbF := openOrCreate(dbPath)
	dbCfgF := openOrCreate(dbCfgPath)
	var dbCfg dbConfig
	err := json.NewDecoder(dbCfgF).Decode(&dbCfg)
	if err != io.EOF && err != nil {
		log.Fatalf("could not decode db config file: %v", err)
	}
	return &DB{
		DB:       dbF,
		DBCfg:    dbCfgF,
		ID:       dbCfg.ID,
		Checksum: genChecksum(dbF),
	}
}

// Read fetches a list of reminders by given ids
func (d *DB) Read(ids []string) []models.Reminder {
	log.Printf("successfully read: %d record(s)\n", len(ids))
	return []models.Reminder{}
}

// ReadAll fetches a list of all reminders
func (d *DB) ReadAll() []models.Reminder {
	resetFilePointer(d.DB)
	var reminders []models.Reminder
	err := json.NewDecoder(d.DB).Decode(&reminders)
	if err != nil && err != io.EOF {
		log.Fatalf("could not decode json from db file: %v", err)
	}
	log.Printf("successfully read: %d record(s)\n", len(reminders))
	return reminders
}

// Write writes a list of reminders to DB
func (d *DB) Write(reminders []models.Reminder) (int, error) {
	resetFilePointer(d.DB)
	bs, err := json.Marshal(reminders)
	if err != nil {
		log.Fatalf("could not marshal json: %v", err)
	}
	bs = append(bs, '\n')
	sum := genChecksum(bytes.NewReader(bs))
	if d.Checksum == sum {
		return 0, nil
	}

	d.Checksum = sum
	newDB, err := os.Create(d.DB.Name())
	if err != nil {
		log.Fatalf("could not create the new db.json: %v", err)
	}
	d.DB = newDB
	err = json.NewEncoder(d.DB).Encode(&reminders)
	if err != nil {
		log.Fatalf("could not encode db json: %v", err)
	}
	newDBCfg, err := os.Create(d.DBCfg.Name())
	if err != nil {
		log.Fatalf("could not create the new .db.config.json: %v", err)
	}
	d.DBCfg = newDBCfg
	dbConfig := dbConfig{
		ID: d.ID,
	}
	err = json.NewEncoder(d.DBCfg).Encode(&dbConfig)
	if err != nil {
		log.Fatalf("could not encode db config json: %v", err)
	}

	log.Printf("successfully wrote: %d record(s)\n", len(reminders))
	return len(reminders), nil
}

// Kill shuts down properly the file database by saving metadata to config file
func (d DB) Kill() {
	defer func() {
		log.Println("closing db files")
		if err := d.DB.Close(); err != nil {
			log.Fatalf("could not close db file: %v", err)
		}
		if err := d.DBCfg.Close(); err != nil {
			log.Fatalf("could not close db file: %v", err)
		}
		log.Println("database was successfully shut down")
	}()
	log.Println("shutting down the database")
	dbCfg := dbConfig{ID: d.ID}
	newDBCfg, err := os.Create(d.DBCfg.Name())
	if err != nil {
		log.Fatalf("could not create new db config file: %v", err)
	}
	err = json.NewEncoder(newDBCfg).Encode(&dbCfg)
	if err != nil {
		log.Fatalf("could not save db config before shutting down")
	}
}

// genCheckSum generates check sum for a reader
func genChecksum(r io.Reader) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		log.Fatalf("could not copy file: %v", err)
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

// GenerateID generates the next AUTOINCREMENT id for a reminder
func (d *DB) GenerateID() int {
	d.ID++
	return d.ID
}

// resetFilePointer resets the file pointer to allow future readings
func resetFilePointer(s io.Seeker) {
	_, err := s.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatalf("could not reset seek pointer: %v", err)
	}
}

// openOrCreate opens a file if it exists, or creates it if it doesn't
func openOrCreate(filePath string) *os.File {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		file, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Fatalf("could not open filw: %v", err)
		}
		return file
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	return file
}
