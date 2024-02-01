package database

import (
	"os"
	"encoding/json"
	"sync"
)

type DB struct {
	path string
	mux *sync.RWMutex
}

type DBStructure struct {
	Chirps map [int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
  db := &DB{
    path: path,
    mux:  &sync.RWMutex{},
  }
  
  // Ensure DB file exists
  if err := db.ensureDB(); err != nil {
    return nil, err
  }

  return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		// Initialize an empty database
		initialDB = `{"chirps":{}}`
		err := os.WriteFile("database.json", []byte(initialDB), 0666)
		if err != nil {
			return err
		}
	} 
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	// Loads the file
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	// Parse the data
	var dbStructure DBStructure
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, nil
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	// Marshal dbStructure to JSON
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	// Write the data to the file
	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}
	
	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

}

