package versions

import (
    "database/sql"
    "fmt"
    "github.com/brianvoe/gofakeit/v7"
    "log"
    "seeder/config"
    "seeder/internal/db_utils"
)

func SeedV1(seedCfg *config.SeedConfig, db *db_utils.DB) error {
    version := 1

    err := SeedTable(db.Conn, version, "difficulty_level", func() error {
        return seedDifficultyLevel(db.Conn)
    })
    if err != nil {
        return err
    } else {
        log.Printf("Seeding difficulty_level table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course_specialization", func() error {
        return seedCourseSpecialization(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Printf("Seeding course_specialization table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course", func() error {
        return seedCourse(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Printf("Seeding course table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course_topic", func() error {
        return seedCourseTopic(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Printf("Seeding course_topic table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course_topic_association", func() error {
        return seedCourseTopicAssociation(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Printf("Seeding course_topic_association table completed successfully")
    }

    err = SeedTable(db.Conn, version, "project", func() error {
        return seedProject(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Printf("Seeding project table completed successfully")
    }

    return nil
}

func seedDifficultyLevel(db *sql.DB) error {
    if tablesExists(db, []string{"difficulty_level"}) != true {
        return nil
    }

    query := `INSERT INTO difficulty_level (name, description)
            VALUES ($1, $2) ON CONFLICT DO NOTHING`

    levelName := []string{
        "beginner",
        "intermediate",
        "advanced",
    }

    levelDescription := []string{
        "for beginners",
        "for intermediate",
        "for advanced",
    }

    var err error
    for i := range len(levelName) {
        _, err = db.Exec(query, levelName[i], levelDescription[i])
    }

    if err != nil {
        return fmt.Errorf("failed to seed (init) difficulty_level: %v", err)
    }

    return nil
}

func seedCourseSpecialization(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course_specialization"}) != true {
        return nil
    }

    query := `INSERT INTO course_specialization (name, description)`

    dataFunc := func(i int) []interface{} {
        return []interface{}{
            gofakeit.JobTitle(),
            gofakeit.Sentence(10),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCourse(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course", "course_specialization", "difficulty_level"}) != true {
        return fmt.Errorf("failed to seed course, some tables not found")
    }

    // Get actual difficulty level IDs that exist in the database
    difficultyLevelIDs, err := getExistingIDs(db, "difficulty_level")
    if err != nil {
        return fmt.Errorf("failed to get difficulty_level ids: %v", err)
    }

    // Get actual specialization IDs that exist in the database
    specializationIDs, err := getExistingIDs(db, "course_specialization")
    if err != nil {
        return fmt.Errorf("failed to get course_specialization ids: %v", err)
    }

    query := `INSERT INTO course (name, description, specialization_id, duration, price, difficulty_level_id)`

    dataFunc := func(i int) []interface{} {
        // Select random IDs from existing IDs
        specializationID := specializationIDs[gofakeit.Number(0, len(specializationIDs)-1)]
        difficultyLevelID := difficultyLevelIDs[gofakeit.Number(0, len(difficultyLevelIDs)-1)]

        return []interface{}{
            gofakeit.BookTitle(),
            gofakeit.Sentence(25),
            specializationID,
            gofakeit.Number(100, 500),
            gofakeit.Number(10000, 200000),
            difficultyLevelID,
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCourseTopic(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course_topic"}) != true {
        return nil
    }

    query := `INSERT INTO course_topic (name, description, technologies, labor_intensity_hours, projects_number)`

    dataFunc := func(i int) []interface{} {
        return []interface{}{
            gofakeit.BookTitle(),
            gofakeit.Sentence(15),
            gofakeit.Sentence(10),
            gofakeit.Number(10, 200),
            gofakeit.Number(1, 10),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCourseTopicAssociation(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course_topic_association", "course", "course_topic"}) != true {
        return fmt.Errorf("failed to seed course_topic_association, some tables not found")
    }

    courseCount, err := getLastTableID(db, "course")
    if err != nil {
        return fmt.Errorf("failed to get course count: %v", err)
    }

    topicCount, err := getLastTableID(db, "course_topic")
    if err != nil {
        return fmt.Errorf("failed to get course_topic count: %v", err)
    }

    query := `INSERT INTO course_topic_association (course_id, topic_id)`

    dataFunc := func(i int) []interface{} {
        return []interface{}{
            gofakeit.Number(1, courseCount),
            gofakeit.Number(1, topicCount),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedProject(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"project", "course_topic"}) != true {
        return fmt.Errorf("failed to seed project, some tables not found")
    }

    topicCount, err := getLastTableID(db, "course_topic")
    if err != nil {
        return fmt.Errorf("failed to get course_topic count: %v", err)
    }

    query := `INSERT INTO project (topic_id, name, description)`

    dataFunc := func(i int) []interface{} {
        return []interface{}{
            gofakeit.Number(1, topicCount),
            gofakeit.BookTitle(),
            gofakeit.Sentence(20),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}
