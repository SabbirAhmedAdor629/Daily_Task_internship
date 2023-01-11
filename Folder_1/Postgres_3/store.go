package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
)package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
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

func OpenDatabase(Host string, Port int, User, Password, DbName string, DbSsl SslType) (*pgClient, error) {
						
									// OPEN A CONNECTION TO THE DATABASE

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

func (d *pgClient) AddHashedAppId(hash hashBuffer, appId string) error {
	// d.Lock()
	// defer d.Unlock()

	var appIds []string
	// var AdId Hash = make(Hash, 32)
	// copy(AdId, hash[:])

	AdId := hex.EncodeToString(hash[:])
	// if !utf8.ValidString(AdId) {
	// 	opCounters.incUtf8Error(1)
	// 	return &DatabaseError{
	// 		Type: "UTF8",
	// 		Err:  errors.New(fmt.Sprintf("UTF8")),
	// 	}
	// }

	// Create the transaction
	tx1, err := d.Begin()
	if err != nil {
		opCounters.incTransactionError(1)
		return &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("Create Transaction Error: %v", err),
		}
	}
	// defer tx.Rollback()
	defer func() {
		_ = tx1.Rollback()
		// if err != nil {
		// 	logger.
		// 		WithFields(Fields{"Database": err.Error()}).
		// 		Error()
		// }
	}()

	// Retrieve any AppIDs associated with the specified AdID
	selectTableAdId := `SELECT app_ids FROM ` + env.DbTable + ` WHERE hashed_ad_id=$1`
	row := tx1.QueryRow(selectTableAdId, AdId) // SELECT

	if err := row.Scan(pq.Array(&appIds)); err != nil {
		if err == sql.ErrNoRows { // Found no AdIDs, and therefore no AppIds
			insertTableAdId :=
				`INSERT INTO ` + env.DbTable +
					` (hashed_ad_id, app_ids, created_at, updated_at) VALUES ($1, $2, $3, $4) ` +
					` ON CONFLICT DO NOTHING`
			_, err := tx1.Exec(insertTableAdId, AdId, pq.Array(append(appIds, appId)), time.Now(), time.Now()) // INSERT
			if err != nil {
				opCounters.incInsertError(1)
				return &DatabaseError{
					Type: "DATABASE",
					Err:  fmt.Errorf("INSERT: %v", err),
				}
			}

			if !env.SkipLogTable {
											// INSERT INTO THE SUPPRESSION LOG TABLE
				insertTableAdId =
					`INSERT INTO ` + env.DbSuppressionLogTable +
						` (hashed_ad_id, app_id, source, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) ` +
						` ON CONFLICT DO NOTHING`
				_, err = tx1.Exec(insertTableAdId, AdId, appId, DEFAULT_POSTBACK_SOURCE, time.Now(), time.Now()) // INSERT
				if err != nil {
					opCounters.incSuppressionLogInsertError(1)
					return &DatabaseError{
						Type: "DATABASE",
						Err:  fmt.Errorf("INSERT: %v", err),
					}
				}
			}

			if err = tx1.Commit(); err != nil {
				opCounters.incTransactionError(1)
				return &DatabaseError{
					Type: "DATABASE",
					Err:  fmt.Errorf("Commit Transaction Error: %v", err),
				}
			}

			opCounters.incNewAdIdCount(1)
			return nil // INSERT successful
		}
		if err != nil {
			opCounters.incSelectError(1)
			return &DatabaseError{
				Type: "DATABASE",
				Err:  fmt.Errorf("SELECT: %v", err),
			}
		}
	}

	if appIdFound, _ := ContainsAppId(appId, appIds); appIdFound {
		opCounters.incExistingAppIdCount(1)
		return nil
	}

	// // Create the transaction
	// tx2, err := d.Begin()
	// if err != nil {
	// 	opCounters.incTransactionError(1)
	// 	return &DatabaseError{
	// 		Type: "DATABASE",
	// 		Err:  fmt.Errorf("Create Transaction Error: %v", err),
	// 	}
	// }
	// // defer tx.Rollback()
	// defer func() {
	// 	err := tx2.Rollback()
	// 	if err != nil {
	// 		logger.
	// 			WithFields(Fields{"Database": err.Error()}).
	// 			Error()
	// 	}
	// }()
	
												//  UPDATE						
	updateTableAdId := `UPDATE ` + env.DbTable + ` SET app_ids=$2, updated_at=$3 WHERE hashed_ad_id=$1`
	_, err = tx1.Exec(updateTableAdId, AdId, pq.Array(append(appIds, appId)), time.Now()) // UPDATE

	if err != nil {
		opCounters.incUpdateError(1)
		return &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("UPDATE: %v", err),
		}
	}
	if !env.SkipLogTable {
		// Insert into the Suppression Log Table
		insertTableAdId :=`INSERT INTO ` + env.DbSuppressionLogTable +
				` (hashed_ad_id, app_id, source, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) ` +
				` ON CONFLICT DO NOTHING`
		_, err = tx1.Exec(insertTableAdId, AdId, appId, DEFAULT_POSTBACK_SOURCE, time.Now(), time.Now()) // INSERT
		if err != nil {
			opCounters.incSuppressionLogInsertError(1)
			return &DatabaseError{
				Type: "DATABASE",
				Err:  fmt.Errorf("INSERT: %v", err),
			}
		}
	}
	if err = tx1.Commit(); err != nil {
		opCounters.incTransactionError(1)
		return &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("Commit Transaction Error: %v", err),
		}
	}

	opCounters.incNewAppIdCount(1)
	return nil
}

// Returns true if AppId is found and the index. Returns false and 0 if AppIDs
// is nil. Returns false and len(AppIDs) if AppID is not found
func ContainsAppId(appId string, appIds []string) (bool, int) {
	if appIds == nil {
		return false, 0
	}
	for i, v := range appIds {
		if appId == v {
			return true, i
		}
	}
	return false, len(appIds)
}

// // Returns a list of App IDs for the specified Ad ID.
// // Returns nil if the Ad ID is not found.
// func (d *DbRepo) Get(adId [32]byte) ([]string, error) {
// 	d.Lock()
// 	defer d.Unlock()

// 	row := d.QueryRow(`SELECT "ad_id", "app_ids" FROM "adids" WHERE "adid"=$1`, adId)

// 	var AdId string
// 	var AppIds []string

// 	if err := row.Scan(&AdId, &AppIds); err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // No matching ad_ids found
// 		}
// 		return nil, err // Something went wrong with the query
// 	}
// 	// var result [][32]byte = make([][32]byte, len(AppIds))
// 	// for i, v := range AppIds {
// 	// 	copy(result[i][:], []byte(v))
// 	// }
// 	return AppIds, nil // matching ad_id found. AppIds may be nil.
// }
//

// func (d *DbRepo) AddHashedAppId(hash hashBuffer, appId string) error {
// 	d.Lock()
// 	defer d.Unlock()

// 	var appIds []string

// 	// We need to covert the byte array into a byte slice
// 	var AdIdSlice []byte = make([]byte, 32)
// 	copy(AdIdSlice, hash[:])

// 	query := `INSERT INTO adids (ad_id, app_ids) ` +
// 		`VALUES ($1, $2) ` +
// 		`ON CONFLICT (ad_id) DO ` +
// 		`UPDATE SET app_ids = (SELECT array_agg(distinct app_ids) FROM unnest(app_ids && $2)) ` +
// 		`WHERE NOT app_ids @> $2`

// 	_, err := d.Exec(query, AdIdSlice, pq.Array(append(appIds, appId)))
// 	if err != nil {
// 		l.database("INSERT Error", err.Error())
// 		return err // Something went wrong with the INSERT
// 	}

// 	return nil
// }


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

func (d *pgClient) AddHashedAppId(hash hashBuffer, appId string) error {
	// d.Lock()
	// defer d.Unlock()

	var appIds []string
	// var AdId Hash = make(Hash, 32)
	// copy(AdId, hash[:])

	AdId := hex.EncodeToString(hash[:])
	// if !utf8.ValidString(AdId) {
	// 	opCounters.incUtf8Error(1)
	// 	return &DatabaseError{
	// 		Type: "UTF8",
	// 		Err:  errors.New(fmt.Sprintf("UTF8")),
	// 	}
	// }

	// Create the transaction
	tx1, err := d.Begin()
	if err != nil {
		opCounters.incTransactionError(1)
		return &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("Create Transaction Error: %v", err),
		}
	}
	// defer tx.Rollback()
	defer func() {
		_ = tx1.Rollback()
		// if err != nil {
		// 	logger.
		// 		WithFields(Fields{"Database": err.Error()}).
		// 		Error()
		// }
	}()
							      // QUERY 
	// Retrieve any AppIDs associated with the specified AdID
	selectTableAdId := `SELECT app_ids FROM ` + env.DbTable + ` WHERE hashed_ad_id=$1`
	row := tx1.QueryRow(selectTableAdId, AdId) // SELECT

	if err := row.Scan(pq.Array(&appIds)); err != nil {
		if err == sql.ErrNoRows { // Found no AdIDs, and therefore no AppIds
			insertTableAdId :=
				`INSERT INTO ` + env.DbTable +
					` (hashed_ad_id, app_ids, created_at, updated_at) VALUES ($1, $2, $3, $4) ` +
					` ON CONFLICT DO NOTHING`
			_, err := tx1.Exec(insertTableAdId, AdId, pq.Array(append(appIds, appId)), time.Now(), time.Now()) // INSERT
			if err != nil {
				opCounters.incInsertError(1)
				return &DatabaseError{
					Type: "DATABASE",
					Err:  fmt.Errorf("INSERT: %v", err),
				}
			}

			if !env.SkipLogTable {
				// Insert into the Suppression Log Table
				insertTableAdId =
					`INSERT INTO ` + env.DbSuppressionLogTable +
						` (hashed_ad_id, app_id, source, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) ` +
						` ON CONFLICT DO NOTHING`
				_, err = tx1.Exec(insertTableAdId, AdId, appId, DEFAULT_POSTBACK_SOURCE, time.Now(), time.Now()) // INSERT
				if err != nil {
					opCounters.incSuppressionLogInsertError(1)
					return &DatabaseError{
						Type: "DATABASE",
						Err:  fmt.Errorf("INSERT: %v", err),
					}
				}
			}

			if err = tx1.Commit(); err != nil {
				opCounters.incTransactionError(1)
				return &DatabaseError{
					Type: "DATABASE",
					Err:  fmt.Errorf("Commit Transaction Error: %v", err),
				}
			}

			opCounters.incNewAdIdCount(1)
			return nil // INSERT successful
		}
		if err != nil {
			opCounters.incSelectError(1)
			return &DatabaseError{
				Type: "DATABASE",
				Err:  fmt.Errorf("SELECT: %v", err),
			}
		}
	}

	if appIdFound, _ := ContainsAppId(appId, appIds); appIdFound {
		opCounters.incExistingAppIdCount(1)
		return nil
	}

	// // Create the transaction
	// tx2, err := d.Begin()
	// if err != nil {
	// 	opCounters.incTransactionError(1)
	// 	return &DatabaseError{
	// 		Type: "DATABASE",
	// 		Err:  fmt.Errorf("Create Transaction Error: %v", err),
	// 	}
	// }
	// // defer tx.Rollback()
	// defer func() {
	// 	err := tx2.Rollback()
	// 	if err != nil {
	// 		logger.
	// 			WithFields(Fields{"Database": err.Error()}).
	// 			Error()
	// 	}
	// }()

	updateTableAdId :=
		`UPDATE ` + env.DbTable + ` SET app_ids=$2, updated_at=$3 WHERE hashed_ad_id=$1`
	_, err = tx1.Exec(updateTableAdId, AdId, pq.Array(append(appIds, appId)), time.Now()) // UPDATE
	if err != nil {
		opCounters.incUpdateError(1)
		return &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("UPDATE: %v", err),
		}
	}
	if !env.SkipLogTable {
										// INSERT INTO THE SUPPRESSION LOG TABLE
		insertTableAdId :=
			`INSERT INTO ` + env.DbSuppressionLogTable +
				` (hashed_ad_id, app_id, source, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) ` +
				` ON CONFLICT DO NOTHING`

		_, err = tx1.Exec(insertTableAdId, AdId, appId, DEFAULT_POSTBACK_SOURCE, time.Now(), time.Now()) // INSERT
		if err != nil {
			opCounters.incSuppressionLogInsertError(1)
			return &DatabaseError{
				Type: "DATABASE",
				Err:  fmt.Errorf("INSERT: %v", err),
			}
		}
	}
	if err = tx1.Commit(); err != nil {
		opCounters.incTransactionError(1)
		return &DatabaseError{
			Type: "DATABASE",
			Err:  fmt.Errorf("Commit Transaction Error: %v", err),
		}
	}

	opCounters.incNewAppIdCount(1)
	return nil
}

// Returns true if AppId is found and the index. Returns false and 0 if AppIDs
// is nil. Returns false and len(AppIDs) if AppID is not found
func ContainsAppId(appId string, appIds []string) (bool, int) {
	if appIds == nil {
		return false, 0
	}
	for i, v := range appIds {
		if appId == v {
			return true, i
		}
	}
	return false, len(appIds)
}

// // Returns a list of App IDs for the specified Ad ID.
// // Returns nil if the Ad ID is not found.
// func (d *DbRepo) Get(adId [32]byte) ([]string, error) {
// 	d.Lock()
// 	defer d.Unlock()

// 	row := d.QueryRow(`SELECT "ad_id", "app_ids" FROM "adids" WHERE "adid"=$1`, adId)

// 	var AdId string
// 	var AppIds []string

// 	if err := row.Scan(&AdId, &AppIds); err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // No matching ad_ids found
// 		}
// 		return nil, err // Something went wrong with the query
// 	}
// 	// var result [][32]byte = make([][32]byte, len(AppIds))
// 	// for i, v := range AppIds {
// 	// 	copy(result[i][:], []byte(v))
// 	// }
// 	return AppIds, nil // matching ad_id found. AppIds may be nil.
// }
//

// func (d *DbRepo) AddHashedAppId(hash hashBuffer, appId string) error {
// 	d.Lock()
// 	defer d.Unlock()

// 	var appIds []string

// 	// We need to covert the byte array into a byte slice
// 	var AdIdSlice []byte = make([]byte, 32)
// 	copy(AdIdSlice, hash[:])

// 	query := `INSERT INTO adids (ad_id, app_ids) ` +
// 		`VALUES ($1, $2) ` +
// 		`ON CONFLICT (ad_id) DO ` +
// 		`UPDATE SET app_ids = (SELECT array_agg(distinct app_ids) FROM unnest(app_ids && $2)) ` +
// 		`WHERE NOT app_ids @> $2`

// 	_, err := d.Exec(query, AdIdSlice, pq.Array(append(appIds, appId)))
// 	if err != nil {
// 		l.database("INSERT Error", err.Error())
// 		return err // Something went wrong with the INSERT
// 	}

// 	return nil
// }
