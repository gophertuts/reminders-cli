package repositories

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/gophertuts/reminders-cli/server/models"
)

// dbConfig represents the config which is used when DB is initialized
type dbConfig struct {
	ID       int    `json:"id"`
	Checksum string `json:"checksum"`
}

// DB represents the application server database (json file)
type DB struct {
	dbPath    string
	dbCfgPath string
	cfg       dbConfig
	db        []byte
}

// NewDB creates a new instance of application file DB
func NewDB(dbPath, dbCfgPath string) *DB {
	db := &DB{
		dbPath:    dbPath,
		dbCfgPath: dbCfgPath,
	}
	return db
}

// Start starts and initializes the file database
func (d *DB) Start() error {
	bs, err := d.read(d.dbCfgPath)
	if err != nil {
		return models.WrapError("could not read db config contents", err)
	}
	var cfg dbConfig
	if len(bs) == 0 {
		bs = []byte("{}")
	}
	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		return models.WrapError("could not unmarshal db config", err)
	}

	bs, err = d.read(d.dbPath)
	if err != nil {
		return models.WrapError("could not read db contents", err)
	}
	d.db = bs
	if d.cfg.Checksum == "" {
		checksum, err := genChecksum(bytes.NewReader(bs))
		if err != nil {
			return err
		}
		cfg.Checksum = checksum
	}
	d.cfg = cfg

	return nil
}

// Read fetches a list of reminders by given ids
func (d *DB) Read(bs []byte) (int, error) {
	n, err := bytes.NewReader(d.db).Read(bs)
	if err != nil && err != io.EOF {
		return 0, models.WrapError("could not read db file bytes", err)
	}
	return n, nil
}

// Write writes a list of reminders to DB
func (d *DB) Write(bs []byte) (int, error) {
	bs = append(bs, '\n')
	checksum, err := genChecksum(bytes.NewReader(bs))
	if err != nil {
		return 0, err
	}
	if d.cfg.Checksum == checksum {
		return 0, nil
	}
	d.cfg.Checksum = checksum

	if err := d.writeDBCfg(); err != nil {
		return 0, err
	}
	n, err := d.write(d.dbPath, bs)
	if err != nil {
		return 0, err
	}
	d.db = bs

	return n, nil
}

// Size retrieves the current size of the database
func (d *DB) Size() int {
	if len(d.db) == 0 {
		d.db = []byte("[]")
	}
	return len(d.db)
}

// GenerateID generates the next AUTOINCREMENT id for a reminder
func (d *DB) GenerateID() int {
	d.cfg.ID++
	return d.cfg.ID
}

// Stop shuts down properly the file database by saving metadata to config file
func (d DB) Stop() error {
	log.Println("shutting down the database")
	_, errDB := os.Open(d.dbPath)
	_, errDBCfg := os.Open(d.dbCfgPath)
	if errors.Is(errDB, os.ErrNotExist) {
		_, err := d.write(d.dbPath, d.db)
		if err != nil {
			return err
		}
	}
	if errors.Is(errDBCfg, os.ErrNotExist) {
		if err := d.writeDBCfg(); err != nil {
			return err
		}
	}
	log.Println("database was successfully shut down")
	return nil
}

// genCheckSum generates check sum for a reader
func genChecksum(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", models.WrapError("could not copy db contents", err)
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

// readContents reads the contents of a db file
func (d *DB) read(path string) ([]byte, error) {
	dbFile, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if errors.Is(err, os.ErrNotExist) {
		dbFile, err = os.Create(path)
	}
	if err != nil {
		return nil, models.WrapError("could not open or create db file", err)
	}
	return ioutil.ReadAll(dbFile)
}

// writeDBCfg writes db config to file
func (d *DB) writeDBCfg() error {
	bs, err := json.Marshal(d.cfg)
	if err != nil {
		return models.WrapError("could not marshal db config", err)
	}
	bs = append(bs, '\n')
	_, err = d.write(d.dbCfgPath, bs)
	if err != nil {
		return models.WrapError("could not write to db cfg file", err)
	}
	return nil
}

func (d *DB) write(path string, bs []byte) (int, error) {
	dbFile, err := os.Create(path)
	if err != nil {
		return 0, models.WrapError("could not create file", err)
	}
	defer d.close(dbFile)

	n, err := dbFile.Write(bs)
	if err == nil {
		log.Printf("successfully wrote %d byte(s) to %s file", n, dbFile.Name())
	}
	return n, err
}

// close closes an open db file
func (d *DB) close(f *os.File) {
	if err := f.Close(); err != nil {
		log.Printf("could not close file '%s': %v", f.Name(), err)
	}
}
