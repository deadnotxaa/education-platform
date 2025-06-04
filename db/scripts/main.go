package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "github.com/brianvoe/gofakeit/v7"
    _ "github.com/lib/pq"
)

type DBConfig struct {
    dbHost     string
    dbPort     string
    dbUser     string
    dbPassword string
    dbName     string
}

type SeedConfig struct {
    seedCount int
}

func main() {
    // Get environment variables
    dbConfig := DBConfig{
        dbHost:     os.Getenv("DB_HOST"),
        dbPort:     os.Getenv("DB_PORT"),
        dbUser:     os.Getenv("DB_USER"),
        dbPassword: os.Getenv("DB_PASSWORD"),
        dbName:     os.Getenv("DB_NAME"),
    }

    seedCnt, err := strconv.Atoi(os.Getenv("SEED_COUNT"))
    currentSeedConfig := SeedConfig{
        seedCount: seedCnt,
    }

    if err != nil {
        currentSeedConfig.seedCount = 10
    }

    // Construct connection string
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbConfig.dbHost, dbConfig.dbPort, dbConfig.dbUser, dbConfig.dbPassword, dbConfig.dbName)

    // Healthcheck: Wait for Postgres to be ready
    db, err := waitForDB(connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    defer func(db *sql.DB) {
        err = db.Close()
        if err != nil {
            log.Printf("Error closing database connection: %v", err)
        }
    }(db)

    if isAlreadySeeded(db) {
        log.Println("Database already seeded, skipping seeding")
        return
    }

    // Seed data
    if err = seedData(db, currentSeedConfig.seedCount); err != nil {
        log.Fatalf("Seeding failed: %v", err)
    }

    if err := markSeeded(db); err != nil {
        log.Fatalf("Failed to mark seeding as complete: %v", err)
    }
    log.Println("Seeding completed successfully")
}

func isAlreadySeeded(db *sql.DB) bool {
    if !tableExists(db, "seeding_status") {
        return false
    }
    var seeded bool
    err := db.QueryRow(`SELECT seeded FROM seeding_status WHERE id = 1`).Scan(&seeded)
    if err != nil {
        log.Printf("Error checking seeding status: %v", err)
        return false
    }
    return seeded
}

func markSeeded(db *sql.DB) error {
    _, err := db.Exec(`UPDATE seeding_status SET seeded = TRUE, seeded_at = CURRENT_TIMESTAMP WHERE id = 1`)
    return err
}

func waitForDB(connStr string) (*sql.DB, error) {
    var db *sql.DB
    var err error

    for i := 0; i < 60; i++ {
        db, err = sql.Open("postgres", connStr)

        if err != nil {
            log.Printf("Failed to open DB: %v, retrying...", err)
            time.Sleep(1 * time.Second)
            continue
        }
        err = db.Ping()

        if err == nil {
            return db, nil
        }
        log.Printf("DB not ready: %v, retrying...", err)

        err = db.Close()
        if err != nil {
            return nil, err
        }

        time.Sleep(1 * time.Second)
    }

    return nil, fmt.Errorf("database not ready after 60 attempts")
}

func tableExists(db *sql.DB, tableName string) bool {
    var exists bool

    query := `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = $1
	)`

    err := db.QueryRow(query, tableName).Scan(&exists)

    return err == nil && exists
}

func seedDifficultyLevel(db *sql.DB) error {
    if tableExists(db, "difficulty_level") != true {
        return nil
    }

    query := `INSERT INTO difficulty_level (name, description)
            VALUES ($1, $2) ON CONFLICT DO NOTHING`

    levelName := []string{"beginner", "intermediate", "advanced"}
    levelDescription := []string{"for beginners", "for intermediate", "for advanced"}

    var err error
    for i := range len(levelName) {
        _, err = db.Exec(query, levelName[i], levelDescription[i])
    }

    if err != nil {
        return fmt.Errorf("failed to seed (init) difficulty_level: %v", err)
    }

    return nil
}

func seedCourseSpecialization(db *sql.DB, seedCount int) error {
    if tableExists(db, "course_specialization") != true {
        return nil
    }

    for i := 0; i < seedCount; i++ {
        _, err := db.Exec(`INSERT INTO course_specialization (name, description) 
				VALUES ($1, $2) ON CONFLICT DO NOTHING`,
            gofakeit.JobTitle(), gofakeit.Sentence(10))
        if err != nil {
            return fmt.Errorf("failed to seed course_specialization: %v", err)
        }
    }

    return nil
}

func seedCourse(db *sql.DB, seedCount int) error {
    if tableExists(db, "course") != true {
        return nil
    }

    // Ensure difficulty_level and course_specialization are seeded
    if err := seedDifficultyLevel(db); err != nil {
        return err
    }
    if err := seedCourseSpecialization(db, seedCount); err != nil {
        return err
    }

    for i := 0; i < seedCount; i++ {
        _, err := db.Exec(`INSERT INTO course (name, description, specialization_id, duration, price, difficulty_level_id) 
            VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`,
            gofakeit.ProductName(), gofakeit.Sentence(20), gofakeit.Number(1, 10), gofakeit.Number(10, 100), gofakeit.Number(100, 1000), gofakeit.Number(1, 3))
        if err != nil {
            return fmt.Errorf("failed to seed course: %v", err)
        }
    }

    return nil
}

func seedData(db *sql.DB, seedCount int) error {
    err := gofakeit.Seed(0)
    if err != nil {
        return err
    }

    err = seedDifficultyLevel(db)
    if err != nil {
        return err
    }

    err = seedCourseSpecialization(db, seedCount)
    if err != nil {
        return err
    }

    // Seed course (requires difficulty_level and course_specialization)
    if tableExists(db, "course") {
        var difficultyIDs, specializationIDs []int
        if tableExists(db, "difficulty_level") {
            rows, err := db.Query(`SELECT id FROM difficulty_level`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        difficultyIDs = append(difficultyIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course_specialization") {
            rows, err := db.Query(`SELECT id FROM course_specialization`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        specializationIDs = append(specializationIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount; i++ {
            var difficultyID, specializationID *int
            if len(difficultyIDs) > 0 {
                id := difficultyIDs[i%len(difficultyIDs)]
                difficultyID = &id
            }
            if len(specializationIDs) > 0 {
                id := specializationIDs[i%len(specializationIDs)]
                specializationID = &id
            }
            _, err := db.Exec(`INSERT INTO course (name, description, specialization_id, duration, price, difficulty_level_id) 
				VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`,
                gofakeit.ProductName(), gofakeit.Sentence(20), specializationID, gofakeit.Number(10, 100), gofakeit.Number(100, 1000), difficultyID)
            if err != nil {
                return fmt.Errorf("failed to seed course: %v", err)
            }
        }
    }

    // Seed course_topic
    if tableExists(db, "course_topic") {
        for i := 0; i < seedCount; i++ {
            _, err := db.Exec(`INSERT INTO course_topic (name, description, technologies, labor_intensity_hours, projects_number) 
				VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
                gofakeit.Word(), gofakeit.Sentence(10), gofakeit.Sentence(5), gofakeit.Number(5, 50), gofakeit.Number(1, 5))
            if err != nil {
                return fmt.Errorf("failed to seed course_topic: %v", err)
            }
        }
    }

    // Seed course_topic_association (requires course and course_topic)
    if tableExists(db, "course_topic_association") {
        var courseIDs, topicIDs []int
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course_topic") {
            rows, err := db.Query(`SELECT id FROM course_topic`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        topicIDs = append(topicIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(courseIDs) > 0 && len(topicIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO course_topic_association (course_id, topic_id) 
				VALUES ($1, $2) ON CONFLICT DO NOTHING`,
                courseIDs[i%len(courseIDs)], topicIDs[i%len(topicIDs)])
            if err != nil {
                return fmt.Errorf("failed to seed course_topic_association: %v", err)
            }
        }
    }

    // Seed project (requires course_topic)
    if tableExists(db, "project") {
        var topicIDs []int
        if tableExists(db, "course_topic") {
            rows, err := db.Query(`SELECT id FROM course_topic`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        topicIDs = append(topicIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount; i++ {
            var topicID *int
            if len(topicIDs) > 0 {
                id := topicIDs[i%len(topicIDs)]
                topicID = &id
            }
            _, err := db.Exec(`INSERT INTO project (topic_id, name, description) 
				VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
                topicID, gofakeit.ProductName(), gofakeit.Sentence(15))
            if err != nil {
                return fmt.Errorf("failed to seed project: %v", err)
            }
        }
    }

    // Seed user
    if tableExists(db, "user") {
        for i := 0; i < seedCount; i++ {
            _, err := db.Exec(`INSERT INTO "user" (name, surname, birthdate, email, hashed_password, profile_picture_url, phone_number, snils_number) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING`,
                gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Date(), gofakeit.Email(), gofakeit.Password(true, true, true, true, false, 12), gofakeit.URL(), gofakeit.Phone(), gofakeit.Numerify("###-###-### ##"))
            if err != nil {
                return fmt.Errorf("failed to seed user: %v", err)
            }
        }
    }

    // Seed role
    if tableExists(db, "role") {
        roles := []string{"teacher", "mentor", "reviewer", "support", "smm_manager"}
        for _, role := range roles {
            _, err := db.Exec(`INSERT INTO role (name) VALUES ($1) ON CONFLICT DO NOTHING`, role)
            if err != nil {
                return fmt.Errorf("failed to seed role: %v", err)
            }
        }
    }

    // Seed employee (requires user and role)
    if tableExists(db, "employee") {
        var userIDs, roleIDs []int
        if tableExists(db, "user") {
            rows, err := db.Query(`SELECT account_id FROM "user"`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        userIDs = append(userIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "role") {
            rows, err := db.Query(`SELECT id FROM role`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        roleIDs = append(roleIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(userIDs) > 0 && len(roleIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO employee (user_id, role_id) 
				VALUES ($1, $2) ON CONFLICT DO NOTHING`,
                userIDs[i%len(userIDs)], roleIDs[i%len(roleIDs)])
            if err != nil {
                return fmt.Errorf("failed to seed employee: %v", err)
            }
        }
    }

    // Seed teacher (requires employee)
    if tableExists(db, "teacher") {
        var employeeIDs []int
        if tableExists(db, "employee") {
            rows, err := db.Query(`SELECT id FROM employee`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        employeeIDs = append(employeeIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(employeeIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO teacher (employee_id, work_place, overall_experience, specialization_experience) 
				VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
                employeeIDs[i%len(employeeIDs)], gofakeit.Company(), gofakeit.Number(1, 20), gofakeit.Number(1, 10))
            if err != nil {
                return fmt.Errorf("failed to seed teacher: %v", err)
            }
        }
    }

    // Seed course_teacher (requires teacher and course)
    if tableExists(db, "course_teacher") {
        var teacherIDs, courseIDs []int
        if tableExists(db, "teacher") {
            rows, err := db.Query(`SELECT employee_id FROM teacher`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        teacherIDs = append(teacherIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(teacherIDs) > 0 && len(courseIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO course_teacher (teacher_id, course_id) 
				VALUES ($1, $2) ON CONFLICT DO NOTHING`,
                teacherIDs[i%len(teacherIDs)], courseIDs[i%len(courseIDs)])
            if err != nil {
                return fmt.Errorf("failed to seed course_teacher: %v", err)
            }
        }
    }

    // Seed course_calendar (requires course)
    if tableExists(db, "course_calendar") {
        var courseIDs []int
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(courseIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO course_calendar (course_id, start_date, end_sales_date, remaining_places) 
				VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
                courseIDs[i%len(courseIDs)], gofakeit.FutureDate(), gofakeit.Date(), gofakeit.Number(0, 50))
            if err != nil {
                return fmt.Errorf("failed to seed course_calendar: %v", err)
            }
        }
    }

    // Seed course_type
    if tableExists(db, "course_type") {
        courseTypes := []string{"basic", "career_track", "discounted"}
        for _, ct := range courseTypes {
            _, err := db.Exec(`INSERT INTO course_type (type_name, discount) 
				VALUES ($1, $2) ON CONFLICT DO NOTHING`,
                ct, gofakeit.Number(0, 50))
            if err != nil {
                return fmt.Errorf("failed to seed course_type: %v", err)
            }
        }
    }

    // Seed purchase (requires user, course, course_type)
    if tableExists(db, "purchase") {
        var userIDs, courseIDs, courseTypeIDs []int
        if tableExists(db, "user") {
            rows, err := db.Query(`SELECT account_id FROM "user"`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        userIDs = append(userIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course_type") {
            rows, err := db.Query(`SELECT id FROM course_type`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseTypeIDs = append(courseTypeIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(userIDs) > 0 && len(courseIDs) > 0 && len(courseTypeIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO purchase (user_id, course_id, purchase_date, course_type_id, total_price, purchase_status) 
				VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`,
                userIDs[i%len(userIDs)], courseIDs[i%len(courseIDs)], gofakeit.Date(), courseTypeIDs[i%len(courseTypeIDs)], gofakeit.Number(100, 1000), gofakeit.RandomString([]string{"paid", "unpaid"}))
            if err != nil {
                return fmt.Errorf("failed to seed purchase: %v", err)
            }
        }
    }

    // Seed career_center_student (requires user, course)
    if tableExists(db, "career_center_student") {
        var userIDs, courseIDs []int
        if tableExists(db, "user") {
            rows, err := db.Query(`SELECT account_id FROM "user"`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        userIDs = append(userIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(userIDs) > 0 && len(courseIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO career_center_student (user_id, course_id, cv_url, career_support_start, support_period) 
				VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
                userIDs[i%len(userIDs)], courseIDs[i%len(courseIDs)], gofakeit.URL(), gofakeit.Date(), gofakeit.Number(1, 12))
            if err != nil {
                return fmt.Errorf("failed to seed career_center_student: %v", err)
            }
        }
    }

    // Seed partner_company
    if tableExists(db, "partner_company") {
        for i := 0; i < seedCount; i++ {
            _, err := db.Exec(`INSERT INTO partner_company (short_name, full_name, hired_graduates_count, requirements, agreement_status) 
				VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
                gofakeit.Company(), gofakeit.Company(), gofakeit.Number(0, 100), gofakeit.Sentence(20), gofakeit.Bool())
            if err != nil {
                return fmt.Errorf("failed to seed partner_company: %v", err)
            }
        }
    }

    // Seed job_application (requires career_center_student, partner_company)
    if tableExists(db, "job_application") {
        var studentIDs, companyIDs []int
        if tableExists(db, "career_center_student") {
            rows, err := db.Query(`SELECT id FROM career_center_student`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        studentIDs = append(studentIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "partner_company") {
            rows, err := db.Query(`SELECT company_id FROM partner_company`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        companyIDs = append(companyIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(studentIDs) > 0 && len(companyIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO job_application (student_id, company_id, application_date, status) 
				VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
                studentIDs[i%len(studentIDs)], companyIDs[i%len(companyIDs)], gofakeit.Date(), gofakeit.RandomString([]string{"pending", "accepted", "rejected"}))
            if err != nil {
                return fmt.Errorf("failed to seed job_application: %v", err)
            }
        }
    }

    // Seed blog_post (requires employee)
    if tableExists(db, "blog_post") {
        var employeeIDs []int
        if tableExists(db, "employee") {
            rows, err := db.Query(`SELECT id FROM employee`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        employeeIDs = append(employeeIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(employeeIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO blog_post (author_id, title, publication_date, topic, reading_time_minutes, cover_image_url, content) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`,
                employeeIDs[i%len(employeeIDs)], gofakeit.Sentence(5), gofakeit.Date(), gofakeit.Word(), gofakeit.Number(5, 30), gofakeit.URL(), gofakeit.Sentence(50))
            if err != nil {
                return fmt.Errorf("failed to seed blog_post: %v", err)
            }
        }
    }

    // Seed tag
    if tableExists(db, "tag") {
        for i := 0; i < seedCount; i++ {
            _, err := db.Exec(`INSERT INTO tag (name) VALUES ($1) ON CONFLICT DO NOTHING`,
                gofakeit.Word())
            if err != nil {
                return fmt.Errorf("failed to seed tag: %v", err)
            }
        }
    }

    // Seed post_tag (requires blog_post, tag)
    if tableExists(db, "post_tag") {
        var postIDs, tagIDs []int
        if tableExists(db, "blog_post") {
            rows, err := db.Query(`SELECT post_id FROM blog_post`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        postIDs = append(postIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "tag") {
            rows, err := db.Query(`SELECT tag_id FROM tag`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        tagIDs = append(tagIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(postIDs) > 0 && len(tagIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO post_tag (post_id, tag_id) 
				VALUES ($1, $2) ON CONFLICT DO NOTHING`,
                postIDs[i%len(postIDs)], tagIDs[i%len(tagIDs)])
            if err != nil {
                return fmt.Errorf("failed to seed post_tag: %v", err)
            }
        }
    }

    // Seed course_review (requires course, user)
    if tableExists(db, "course_review") {
        var courseIDs, userIDs []int
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "user") {
            rows, err := db.Query(`SELECT account_id FROM "user"`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        userIDs = append(userIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(courseIDs) > 0 && len(userIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO course_review (course_id, user_id, rating, comment, review_date) 
				VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
                courseIDs[i%len(courseIDs)], userIDs[i%len(userIDs)], gofakeit.Number(1, 5), gofakeit.Sentence(20), gofakeit.Date())
            if err != nil {
                return fmt.Errorf("failed to seed course_review: %v", err)
            }
        }
    }

    // Seed certificate (requires user, course)
    if tableExists(db, "certificate") {
        var userIDs, courseIDs []int
        if tableExists(db, "user") {
            rows, err := db.Query(`SELECT account_id FROM "user"`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        userIDs = append(userIDs, id)
                    }
                }
                rows.Close()
            }
        }
        if tableExists(db, "course") {
            rows, err := db.Query(`SELECT course_id FROM course`)
            if err == nil {
                for rows.Next() {
                    var id int
                    if err := rows.Scan(&id); err == nil {
                        courseIDs = append(courseIDs, id)
                    }
                }
                rows.Close()
            }
        }
        for i := 0; i < seedCount && len(userIDs) > 0 && len(courseIDs) > 0; i++ {
            _, err := db.Exec(`INSERT INTO certificate (user_id, course_id, issue_date) 
				VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
                userIDs[i%len(userIDs)], courseIDs[i%len(courseIDs)], gofakeit.Date())
            if err != nil {
                return fmt.Errorf("failed to seed certificate: %v", err)
            }
        }
    }

    return nil
}
