package database

import (
	"os"
	"encoding/json"
	"sync"
	"sort"
)

type DB struct {
	path string
	mux *sync.RWMutex
}

type DBStructure struct {
	Chirps map [int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
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
		initialDB := `{"chirps":{}, "users":{}}`
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
	if err != nil {
		return err
	}
	
	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load existing chirps
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Generate an id for the new chirp
	newID := len(dbStructure.Chirps) + 1

	// Create a new chirp
	newChirp := Chirp{
		ID: newID,
		Body: body,
	}

	// Add the new chirp into the Chirps map
	dbStructure.Chirps[newID] = newChirp

	// Write the update chirps back to the file
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	// Load existing chirps
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	// Sort chirps by ID
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
		})


	return chirps, nil
}

func (db *DB) CreateUser(email string) (*User, error) {
    	db.mux.Lock()	
	defer db.mux.Unlock()
	
	// Load existing users
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	
	// Generate a new id for the user
	newID := len(dbStructure.Users) + 1

	// Create a new user.
	newUser := &User{
		Id:    newID,
		Email: email,
	}

	// Add the new user to the Users map
	dbStructure.Users[newID] = *newUser

	// Write the updated users back to the file
	err = db.writeDB(dbStructure)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

