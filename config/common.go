package config

import (
	"encoding/json"
	"fmt"

	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/db/postgres"
	"github.com/oklog/run"
)

var (
	G   run.Group
	Cfg *Config
)

type Config struct {
	// Storage Storage
	AutoMigrate bool   `json:"autoMigrate"`
	Log         string `json:"log"`
	Server      Server
	Proxy       []Proxy
}

var storages = map[string]func() db.Storage{
	"postgres": func() db.Storage { return new(postgres.PostgresStorage) },
	// "mysql":      func() db.Storage { return new(sql.MySQL) },
	// "sqlite3": func() db.Storage { return new(sqlite.SQLite3) },
}

type Proxy struct {
	Name    string  `json:"name"`
	Desc    string  `json:"desc"`
	Storage Storage `json:"storage"`
}

type Storage struct {
	Type   string     `json:"type"`
	Config db.Storage `json:"config"`
}

// UnmarshalJSON allows Storage to implement the unmarshaler interface to
// dynamically determine the type of the storage config.
func (s *Storage) UnmarshalJSON(b []byte) error {
	var store struct {
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config"`
	}
	if err := json.Unmarshal(b, &store); err != nil {
		return fmt.Errorf("parse storage: %v", err)
	}
	f, ok := storages[store.Type]
	if !ok {
		return fmt.Errorf("unknown storage type %q", store.Type)
	}

	storageConfig := f()
	if len(store.Config) != 0 {
		data := []byte(store.Config)
		if err := json.Unmarshal(data, storageConfig); err != nil {
			return fmt.Errorf("parse storage config: %v", err)
		}
	}
	*s = Storage{
		Type:   store.Type,
		Config: storageConfig,
	}
	return nil
}

type Server struct {
	HTTP           string   `json:"http"`
	HTTPS          string   `json:"https"`
	Prefix         string   `json:"prefix"`
	TLSCert        string   `json:"tlsCert"`
	TLSKey         string   `json:"tlsKey"`
	AllowedOrigins []string `json:"allowedOrigins"`
}
