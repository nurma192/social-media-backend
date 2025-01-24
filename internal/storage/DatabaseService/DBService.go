package DatabaseService

import (
	"database/sql"
)

type DBService struct {
	DB *sql.DB
}

func NewDBService(db *sql.DB) *DBService {
	return &DBService{DB: db}
}
