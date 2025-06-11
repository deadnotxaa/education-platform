package seed

import (
	"log"
	"seeder/config"
	"seeder/internal/db_utils"
	"seeder/internal/versions"
	"strconv"
)

func RunSeeding(cfg *config.Config) {
    dbInstance, err := db_utils.NewDB(cfg)
    if err != nil {
        log.Fatalf("Could not connect to DB: %v", err)
    }

    defer dbInstance.CloseConnection()

    migrationVersion, err := checkMigrationVersion(dbInstance)
    if err != nil {
        log.Fatalf("Could not check seeding version: %v", err)
    }

    log.Printf("Current migration version: %v", migrationVersion)

    migrationMap := map[int]func(version int, seedConfig *config.SeedConfig, db *db_utils.DB) error{
        1: func(version int, seedConfig *config.SeedConfig, db *db_utils.DB) error {
            return versions.SeedV1(seedConfig, db)
        },
        2: func(version int, seedConfig *config.SeedConfig, db *db_utils.DB) error {
            return versions.SeedV2(seedConfig, db)
        },
        3: func(version int, seedConfig *config.SeedConfig, db *db_utils.DB) error {
            return versions.SeedV3(seedConfig, db)
        },
    }

    var maxMigrationVersion int

    if cfg.SeedConfig.MigrationVersion == "latest" {
        maxMigrationVersion = len(migrationMap)
    } else {
        maxMigrationVersion, err = strconv.Atoi(cfg.SeedConfig.MigrationVersion)
    }

    log.Printf("[pre] Migration version set to: %d", maxMigrationVersion)

    for version := 1; version <= maxMigrationVersion; version++ {
        if version <= maxMigrationVersion && version > migrationVersion {
            log.Printf("Seeding migration version: %v", version)

            err = versions.CleanupFailedSeeding(dbInstance.Conn, version)
            if err != nil {
                log.Fatalf("Error cleaning up failed seeding attempts: %v", err)
            }

            err = migrationMap[version](version, &cfg.SeedConfig, dbInstance)

            if err != nil {
                log.Fatalf("Error seeding version %d: %v", version, err)
            } else {
                log.Printf("Successfully seeded version %d", version)

                err = versions.UpdateMigrationVersion(dbInstance.Conn, version)
                if err != nil {
                    log.Fatalf("Error setting migration version to %d: %v", version, err)
                }
            }
        }
    }

    log.Printf("Successfully seeded all versions up to %d", maxMigrationVersion)
}

func checkMigrationVersion(dbInstance *db_utils.DB) (int, error) {
    var version int

    err := dbInstance.Conn.QueryRow("SELECT migration_version FROM seeding_status WHERE table_name IS NULL").Scan(&version)
    if err != nil {
        return 0, err
    }

    return version, nil
}

func SetMigrationVersion(dbInstance *db_utils.DB, version int) error {
	_, err := dbInstance.Conn.Exec("UPDATE seeding_status SET migration_version = $1", version)
	if err != nil {
		return err
	}

	log.Printf("[post] Migration version set to %d", version)
	return nil
}
