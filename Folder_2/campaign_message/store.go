package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
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

	// mockPgClient struct {
	// 	sync.RWMutex
	// }

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

type (
	Hash []byte
)

func (h Hash) String() string { return hex.EncodeToString(h) }

type CampaignTable struct {
	Id                string
	Active            bool
	ScheduleTime      sql.NullInt32
	EventName         string
	MessageTemplateId sql.NullInt32
	ScheduleDays      ScheduleDays
}

type MessageTemplate struct {
	Id         string
	Active     bool
	DailyCap   sql.NullInt32
	MonthlyCap sql.NullInt32
	TotalCap   sql.NullInt32
}

type ScheduleDays []interface{}

func (a ScheduleDays) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ScheduleDays) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func OpenDatabase(Host string, Port int, User, Password, DbName string, DbSsl SslType) (*pgClient, error) {
	// Open a connection to the database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		Host, Port, User, Password, DbName, SslConnection[DbSsl]))
	if err != nil {
		return nil, &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("OPEN: %v", err),
		}
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("PING: %v", err),
		}
	}

	// createTableAdId := `CREATE TABLE IF NOT EXISTS ` + env.DbTable + ` (hashed_ad_id BYTEA, app_ids TEXT[], created_at TIMESTAMP, updated_at TIMESTAMP, PRIMARY KEY(hashed_ad_id));`
	// if _, err := db.Exec(createTableAdId); err != nil {
	// 	logEntry.Error(Fields{"Create": err.Error()}) // Cannot create the table
	// 	return nil, err
	// }

	return &pgClient{
		RWMutex: sync.RWMutex{},
		DB:      db,
	}, nil
}

func (client *pgClient) UpdateWorkLog(log_chunk_id int, status string) error {
	// Update into the logChunk Table
	timestamp := time.Now()
	query := `UPDATE ` + env.DbLogChunkTable
	query += ` SET status = $1,`
	query += ` status_updated_at = $2`
	query += ` WHERE id = $3;`
	_, err := client.Exec(query, status, timestamp, log_chunk_id)
	if err != nil {
		return fmt.Errorf("UPDATE: %s", err.Error())
	}
	return nil
}

func (client *pgClient) ReadCampaignData(campaign_id int) (*CampaignTable, error) {
	campaignTable := new(CampaignTable)
	sqlStatement := `SELECT schedule_time, event_name, message_template_id, active, schedule_days FROM campaigns WHERE id=$1`

	err := client.QueryRow(sqlStatement, campaign_id).Scan(&campaignTable.ScheduleTime, &campaignTable.EventName, &campaignTable.MessageTemplateId, &campaignTable.Active, &campaignTable.ScheduleDays)
	if err != nil {
		return nil, err
	}
	return campaignTable, nil
}

func (client *pgClient) ReadMessageTemplateData(message_template_id int) (*MessageTemplate, error) {
	messageTemplate := new(MessageTemplate)
	sqlStatement := `SELECT active, daily_cap, monthly_cap, total_cap FROM message_templates WHERE id=$1`

	err := client.QueryRow(sqlStatement, message_template_id).Scan(&messageTemplate.Active, &messageTemplate.DailyCap, &messageTemplate.MonthlyCap, &messageTemplate.TotalCap)
	if err != nil {
		return nil, err
	}
	return messageTemplate, nil
}
