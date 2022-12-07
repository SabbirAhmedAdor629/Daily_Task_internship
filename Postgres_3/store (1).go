package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type (
	pgClient struct {
		sync.RWMutex
		*sql.DB
	}

	mockPgClient struct {
		sync.RWMutex
	}

	SslType int
)

const (
	AlwaysSsl SslType = iota
	NoSsl
	VerifyCa
	VerifyFull
)

var (
	SslConnection = []string{
		AlwaysSsl:  "require",
		NoSsl:      "disable",
		VerifyCa:   "verify-ca",
		VerifyFull: "verify-full",
	}
)

func OpenDatabase(host string, port int, user, password, dbName string, dbSsl SslType) (*pgClient, error) {
	// Open a connection to the database
	sqlDb, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, SslConnection[dbSsl]))
	if err != nil {
		return nil, fmt.Errorf("OPEN: %s", err.Error())
	}

	sqlDb.SetConnMaxLifetime(15 * time.Minute)

	return &pgClient{
		RWMutex: sync.RWMutex{},
		DB:      sqlDb,
	}, nil
}

func (d *pgClient) StatusUpdate(recordId int, status string) error {
	query := `UPDATE zendesk_accounts`
	query += ` SET zendesk_status = $1,`
	query += ` zendesk_status_updated_at = $2`
	query += ` WHERE id = $3;`

	timestamp := time.Now()

	_, err := d.Exec(query, status, timestamp, recordId)
	if err != nil {
		return fmt.Errorf("EXECUTION_ERROR: %s", err.Error())
	}

	return nil
}

func (d *pgClient) ZendeskUserIdUpdate(recordId int, zendeskUserId int, status string) error {
	query := `UPDATE zendesk_accounts`
	query += ` SET zendesk_user_id = $1,`
	query += ` zendesk_status = $2,`
	query += ` zendesk_status_updated_at = $3`
	query += ` WHERE id = $4;`

	timestamp := time.Now()

	_, err := d.Exec(query, zendeskUserId, status, timestamp, recordId)
	if err != nil {
		return fmt.Errorf("EXECUTION_ERROR: %s", err.Error())
	}

	return nil
}
