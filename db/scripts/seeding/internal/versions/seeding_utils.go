package versions

import (
    "database/sql"
    "fmt"
    "github.com/lib/pq"
    "strings"
)

func trackSeedingStart(db *sql.DB, version int, tableName string) error {
    query := `INSERT INTO seeding_status 
              (migration_version, table_name, success, seeded_at) 
              VALUES ($1, $2, false, CURRENT_TIMESTAMP)`

    _, err := db.Exec(query, version, tableName)
    if err != nil {
        return fmt.Errorf("failed to track seeding start for table %s: %v", tableName, err)
    }

    return nil
}

func trackSeedingSuccess(db *sql.DB, version int, tableName string) error {
    query := `UPDATE seeding_status 
              SET success = true, seeded_at = CURRENT_TIMESTAMP 
              WHERE migration_version = $1 AND table_name = $2`

    _, err := db.Exec(query, version, tableName)
    if err != nil {
        return fmt.Errorf("failed to track seeding success for table %s: %v", tableName, err)
    }

    return nil
}

func isTableAlreadySeeded(db *sql.DB, version int, tableName string) (bool, error) {
    query := `SELECT EXISTS (
              SELECT 1 FROM seeding_status 
              WHERE migration_version = $1 
              AND table_name = $2 
              AND success = true)`

    var seeded bool
    err := db.QueryRow(query, version, tableName).Scan(&seeded)
    if err != nil {
        return false, fmt.Errorf("failed to check if table %s is already seeded: %v", tableName, err)
    }

    return seeded, nil
}

func CleanupFailedSeeding(db *sql.DB, version int) error {
    query := `DELETE FROM seeding_status 
              WHERE migration_version = $1 
              AND success = false`

    _, err := db.Exec(query, version)
    if err != nil {
        return fmt.Errorf("failed to clean up failed seeding for version %d: %v", version, err)
    }

    return nil
}

func UpdateMigrationVersion(db *sql.DB, version int) error {
    query := `UPDATE seeding_status 
              SET migration_version = $1, success = true, seeded_at = CURRENT_TIMESTAMP 
              WHERE table_name IS NULL`

    _, err := db.Exec(query, version)
    if err != nil {
        return fmt.Errorf("failed to update migration version to %d: %v", version, err)
    }

    return nil
}

func tablesExists(db *sql.DB, tableNames []string) bool {
    var exists bool
    var err error

    query := `SELECT EXISTS (
              SELECT FROM information_schema.tables
              WHERE table_schema = 'public'
              AND table_name = $1)`

    for _, tableName := range tableNames {
        err = db.QueryRow(query, tableName).Scan(&exists)

        if err != nil || exists == false {
            fmt.Printf("Error checking if table %s exists: %v\n", tableName, err)
            return false
        }
    }

    return true
}

func getLastIDFromTable(db *sql.DB, tableName, idColumnName string) (int, error) {
    query := fmt.Sprintf("SELECT MAX(%s) FROM %s",
                          pq.QuoteIdentifier(idColumnName),
                          pq.QuoteIdentifier(tableName))

    var lastID sql.NullInt64

    err := db.QueryRow(query).Scan(&lastID)
    if err != nil {
        return 0, fmt.Errorf("failed to get last ID from %s: %v", tableName, err)
    }

    if !lastID.Valid {
        return 0, nil
    }

    return int(lastID.Int64), nil
}

func getTablePrimaryKeyName(db *sql.DB, tableName string) (string, error) {
    query := `
        SELECT a.attname
        FROM pg_index i
        JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
        WHERE i.indrelid = $1::regclass
        AND i.indisprimary
        LIMIT 1`

    var pkColumnName string
    err := db.QueryRow(query, tableName).Scan(&pkColumnName)
    if err != nil {
        return "", fmt.Errorf("failed to get primary key column for %s: %v", tableName, err)
    }

    return pkColumnName, nil
}

func getLastTableID(db *sql.DB, tableName string) (int, error) {
    pkColumnName, err := getTablePrimaryKeyName(db, tableName)
    if err != nil {
        return 0, err
    }

    return getLastIDFromTable(db, tableName, pkColumnName)
}

func getExistingIDs(db *sql.DB, tableName string, idColumnName ...string) ([]int, error) {
    var columnName string

    if len(idColumnName) > 0 && idColumnName[0] != "" {
        columnName = idColumnName[0]
    } else {
        var err error
        columnName, err = getTablePrimaryKeyName(db, tableName)
        if err != nil {
            return nil, fmt.Errorf("failed to get primary key for %s: %v", tableName, err)
        }
    }

    query := fmt.Sprintf("SELECT %s FROM %s",
                         pq.QuoteIdentifier(columnName),
                         pq.QuoteIdentifier(tableName))

    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to get IDs from %s: %v", tableName, err)
    }
    defer func(rows *sql.Rows) {
        err := rows.Close()
        if err != nil {
            fmt.Printf("Error closing rows: %v\n", err)
        }
    }(rows)

    var ids []int
    for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            return nil, fmt.Errorf("failed to scan ID: %v", err)
        }
        ids = append(ids, id)
    }

    if len(ids) == 0 {
        return nil, fmt.Errorf("no records found in table %s", tableName)
    }

    return ids, nil
}

func SeedTable(db *sql.DB, version int, tableName string, seedFunc func() error) error {
    alreadySeeded, err := isTableAlreadySeeded(db, version, tableName)
    if err != nil {
        return fmt.Errorf("failed to check if %s was already seeded: %v", tableName, err)
    }

    if alreadySeeded {
        fmt.Printf("Table %s already seeded in version %d, skipping\n", tableName, version)
        return nil
    }

    if err = trackSeedingStart(db, version, tableName); err != nil {
        return err
    }

    if err = seedFunc(); err != nil {
        return fmt.Errorf("failed to seed %s: %v", tableName, err)
    }

    if err = trackSeedingSuccess(db, version, tableName); err != nil {
        return err
    }

    fmt.Printf("Successfully seeded %s\n", tableName)
    return nil
}

// BatchInsertData performs batch inserts for better performance
func BatchInsertData(db *sql.DB, query string, batchSize int, dataFunc func(int) []interface{}, totalCount int) error {
    if totalCount == 0 {
        return nil
    }

    // If total count is less than batch size, insert all at once
    if totalCount <= batchSize {
        return executeBatchInsert(db, query, dataFunc, totalCount)
    }

    // Process in batches
    for i := 0; i < totalCount; i += batchSize {
        currentBatchSize := batchSize
        if i+batchSize > totalCount {
            currentBatchSize = totalCount - i
        }

        if err := executeBatchInsert(db, query, func(idx int) []interface{} {
            return dataFunc(i + idx)
        }, currentBatchSize); err != nil {
            return err
        }
    }

    return nil
}

// executeBatchInsert executes a single batch insert
func executeBatchInsert(db *sql.DB, baseQuery string, dataFunc func(int) []interface{}, batchSize int) error {
    if batchSize == 0 {
        return nil
    }

    // Get the first row to determine the number of columns
    firstRow := dataFunc(0)
    numCols := len(firstRow)

    // Build the VALUES clause dynamically
    var valuePlaceholders []string
    var args []interface{}

    for i := 0; i < batchSize; i++ {
        row := dataFunc(i)
        if len(row) != numCols {
            return fmt.Errorf("inconsistent number of columns in row %d", i)
        }

        var placeholders []string
        for j := 0; j < numCols; j++ {
            placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)+1))
            args = append(args, row[j])
        }
        valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ", ")+")")
    }

    // Combine base query with VALUES clause and ON CONFLICT
    fullQuery := baseQuery + " VALUES " + strings.Join(valuePlaceholders, ", ") + " ON CONFLICT DO NOTHING"

    _, err := db.Exec(fullQuery, args...)
    return err
}
