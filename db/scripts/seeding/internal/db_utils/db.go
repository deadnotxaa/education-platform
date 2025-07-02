package db_utils

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "seeder/config"
    "time"
)

type DB struct {
    Conn *sql.DB
    cfg  *config.Config
}

func NewDB(cfg *config.Config) (*DB, error) {
    db := &DB{
        cfg: cfg,
    }

    err := db.connectToDB()
    if err != nil {
        return nil, err
    }

    return db, nil
}

func (db *DB) CloseConnection() {
    err := db.Conn.Close()
    if err != nil {
        log.Printf("Error closing database connection: %v", err)
    }
}

func (db *DB) connectToDB() error {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        db.cfg.DBConfig.DbHost,
        db.cfg.DBConfig.DbPort,
        db.cfg.DBConfig.DbUser,
        db.cfg.DBConfig.DbPassword,
        db.cfg.DBConfig.DbName,
    )

    log.Printf("Connecting to DB with connection string: host=%s port=%s dbname=%s user=%s",
        db.cfg.DBConfig.DbHost, db.cfg.DBConfig.DbPort, db.cfg.DBConfig.DbName, db.cfg.DBConfig.DbUser)

    for i := 0; i < 60; i++ {
        conn, err := sql.Open("postgres", connStr)

        if err != nil {
            log.Printf("Failed to open DB: %v, retrying...", err)
            time.Sleep(1 * time.Second)
            continue
        }

        err = conn.Ping()
        if err == nil {
            db.Conn = conn
            return nil
        }

        log.Printf("DB not ready: %v, retrying...", err)
        err = conn.Close()
        if err != nil {
            return err
        }

        time.Sleep(1 * time.Second)
    }

    return fmt.Errorf("database not ready after 60 attempts")
}
