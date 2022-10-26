package db

import (
	log "github.com/sirupsen/logrus"
)

// NetworkDB contains options common to SQL databases accessed over network.
type NetworkDB struct {
	Database string
	User     string
	Password string
	Host     string
	Port     uint16
	Debug    bool

	ConnectionTimeout int // Seconds

	// database/sql tunables, see
	// https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime and below
	// Note: defaults will be set if these are 0
	MaxOpenConns    int `json:"maxOpenConns"`    // default: 5
	MaxIdleConns    int `json:"maxIdleConns"`    // default: 5
	ConnMaxLifetime int `json:"connMaxLifetime"` // Seconds, default: not set
}

// SSL represents SSL options for network databases.
type SSL struct {
	Mode   string
	CAFile string `json:"cAFile"`
	// Files for client auth.
	KeyFile  string `json:"keyFile"`
	CertFile string `json:"certFile"`
}

type Storage interface {
	Open(logger *log.Logger) (*Conn, error)
}
