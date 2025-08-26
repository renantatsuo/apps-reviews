package apps

import (
	"database/sql"

	"github.com/renantatsuo/app-review/server/pkg/apple"
)

type AppsClient struct {
	db          *sql.DB
	appleClient *apple.AppleClient
}

func New(db *sql.DB, appleClient *apple.AppleClient) *AppsClient {
	return &AppsClient{db: db, appleClient: appleClient}
}
