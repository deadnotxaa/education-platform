package versions

import (
    "database/sql"
    "fmt"
    "github.com/brianvoe/gofakeit/v7"
    "log"
    "seeder/config"
    "seeder/internal/db_utils"
)

func SeedV3(seedCfg *config.SeedConfig, db *db_utils.DB) error {
    version := 3

    err := SeedTable(db.Conn, version, "course_calendar", func() error {
        return seedCourseCalendar(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding course_calendar table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course_type", func() error {
        return seedCourseType(db.Conn)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding course_type table completed successfully")
    }

    err = SeedTable(db.Conn, version, "purchase", func() error {
        return seedPurchase(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding purchase table completed successfully")
    }

    err = SeedTable(db.Conn, version, "career_center_student", func() error {
        return seedCareerCenterStudent(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding career_center_student table completed successfully")
    }

    err = SeedTable(db.Conn, version, "partner_company", func() error {
        return seedPartnerCompany(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding partner_company table completed successfully")
    }

    err = SeedTable(db.Conn, version, "job_application", func() error {
        return seedJobApplication(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding job_application table completed successfully")
    }

    err = SeedTable(db.Conn, version, "blog_post", func() error {
        return seedBlogPost(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding blog_post table completed successfully")
    }

    err = SeedTable(db.Conn, version, "tag", func() error {
        return seedTag(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding tag table completed successfully")
    }

    err = SeedTable(db.Conn, version, "post_tag", func() error {
        return seedPostTag(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding post_tag table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course_review", func() error {
        return seedCourseReview(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding course_review table completed successfully")
    }

    err = SeedTable(db.Conn, version, "certificate", func() error {
        return seedCertificate(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding certificate table completed successfully")
    }

    return nil
}

func seedCourseCalendar(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course_calendar"}) != true {
        return nil
    }

    courseIDs, err := getExistingIDs(db, "course", "course_id")
    if err != nil {
        return fmt.Errorf("failed to get course ids: %v", err)
    }

    query := `INSERT INTO course_calendar (course_id, start_date, end_sales_date, remaining_places)`

    dataFunc := func(i int) []interface{} {
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        startDate := gofakeit.Date()
        endSalesDate := gofakeit.DateRange(startDate, gofakeit.Date())

        return []interface{}{
            courseID,
            startDate,
            endSalesDate,
            gofakeit.Number(0, 100),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCourseType(db *sql.DB) error {
    if tablesExists(db, []string{"course_type"}) != true {
        return nil
    }

    query := `INSERT INTO course_type (type_name, discount)
              VALUES ($1, $2)`

    types := []struct {
        name     string
        discount int
    }{
        {"Standard", 0},
        {"Students", 10},
        {"Disabled", 20},
        {"Veteran", 70},
    }

    for _, t := range types {
        _, err := db.Exec(query, t.name, t.discount)
        if err != nil {
            return fmt.Errorf("failed to insert course_type: %v", err)
        }
    }

    return nil
}

func seedPurchase(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"purchase"}) != true {
        return nil
    }

    userIDs, err := getExistingIDs(db, "users", "account_id")
    if err != nil {
        return fmt.Errorf("failed to get user ids: %v", err)
    }

    courseIDs, err := getExistingIDs(db, "course", "course_id")
    if err != nil {
        return fmt.Errorf("failed to get course ids: %v", err)
    }

    courseTypeIDs, err := getExistingIDs(db, "course_type", "id")
    if err != nil {
        return fmt.Errorf("failed to get course type ids: %v", err)
    }

    query := `INSERT INTO purchase (user_id, course_id, purchase_date, course_type_id, total_price, purchase_status)`

    dataFunc := func(i int) []interface{} {
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        courseTypeID := courseTypeIDs[gofakeit.Number(0, len(courseTypeIDs)-1)]

        return []interface{}{
            userID,
            courseID,
            gofakeit.Date(),
            courseTypeID,
            gofakeit.Number(1000, 5000),
            gofakeit.RandomString([]string{"Completed", "Pending", "Cancelled"}),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedPartnerCompany(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"partner_company"}) != true {
        return nil
    }

    query := `INSERT INTO partner_company (short_name, full_name, hired_graduates_count, requirements, agreement_status)`

    dataFunc := func(i int) []interface{} {
        return []interface{}{
            gofakeit.Company(),
            gofakeit.Company() + " LLC",
            gofakeit.Number(0, 100),
            gofakeit.Sentence(10),
            gofakeit.Bool(),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCareerCenterStudent(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"career_center_student"}) != true {
        return nil
    }

    userIDs, err := getExistingIDs(db, "users", "account_id")
    if err != nil {
        return fmt.Errorf("failed to get user ids: %v", err)
    }

    courseIDs, err := getExistingIDs(db, "course", "course_id")
    if err != nil {
        return fmt.Errorf("failed to get course ids: %v", err)
    }

    query := `INSERT INTO career_center_student (user_id, course_id, cv_url, career_support_start, support_period)`

    dataFunc := func(i int) []interface{} {
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]

        return []interface{}{
            userID,
            courseID,
            gofakeit.URL(),
            gofakeit.Date(),
            gofakeit.Number(1, 12),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedJobApplication(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"job_application"}) != true {
        return nil
    }

    studentIDs, err := getExistingIDs(db, "career_center_student", "id")
    if err != nil {
        return fmt.Errorf("failed to get student ids: %v", err)
    }

    companyIDs, err := getExistingIDs(db, "partner_company", "company_id")
    if err != nil {
        return fmt.Errorf("failed to get company ids: %v", err)
    }

    query := `INSERT INTO job_application (student_id, company_id, application_date, status)`

    dataFunc := func(i int) []interface{} {
        studentID := studentIDs[gofakeit.Number(0, len(studentIDs)-1)]
        companyID := companyIDs[gofakeit.Number(0, len(companyIDs)-1)]

        return []interface{}{
            studentID,
            companyID,
            gofakeit.Date(),
            gofakeit.RandomString([]string{"Applied", "Interviewed", "Rejected", "Hired"}),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedBlogPost(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"blog_post", "employee"}) != true {
        return fmt.Errorf("failed to seed blog_post, some tables not found")
    }

    employeeIDs, err := getExistingIDs(db, "employee", "id")
    if err != nil {
        return fmt.Errorf("failed to get employee ids: %v", err)
    }

    query := `INSERT INTO blog_post (author_id, title, publication_date, topic, reading_time_minutes, cover_image_url, content)`

    dataFunc := func(i int) []interface{} {
        authorID := employeeIDs[gofakeit.Number(0, len(employeeIDs)-1)]

        return []interface{}{
            authorID,
            gofakeit.BookTitle(),
            gofakeit.Date(),
            gofakeit.Sentence(3),
            gofakeit.Number(5, 30),
            gofakeit.URL(),
            gofakeit.Paragraph(5, 10, 100, "\n"),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedTag(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"tag"}) != true {
        return nil
    }

    query := `INSERT INTO tag (name)`

    dataFunc := func(i int) []interface{} {
        return []interface{}{
            gofakeit.HackerAdjective() + " " + gofakeit.HackerNoun(),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedPostTag(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"post_tag", "blog_post", "tag"}) != true {
        return fmt.Errorf("failed to seed post_tag, some tables not found")
    }

    postIDs, err := getExistingIDs(db, "blog_post")
    if err != nil {
        return fmt.Errorf("failed to get blog_post IDs: %v", err)
    }

    tagIDs, err := getExistingIDs(db, "tag")
    if err != nil {
        return fmt.Errorf("failed to get tag IDs: %v", err)
    }

    query := `INSERT INTO post_tag (post_id, tag_id)`

    dataFunc := func(i int) []interface{} {
        postID := postIDs[gofakeit.Number(0, len(postIDs)-1)]
        tagID := tagIDs[gofakeit.Number(0, len(tagIDs)-1)]

        return []interface{}{
            postID,
            tagID,
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCourseReview(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course_review", "course", "users"}) != true {
        return fmt.Errorf("failed to seed course_review, some tables not found")
    }

    userIDs, err := getExistingIDs(db, "users", "account_id")
    if err != nil {
        return fmt.Errorf("failed to get user ids: %v", err)
    }

    courseIDs, err := getExistingIDs(db, "course", "course_id")
    if err != nil {
        return fmt.Errorf("failed to get course ids: %v", err)
    }

    query := `INSERT INTO course_review (course_id, user_id, rating, comment, review_date)`

    dataFunc := func(i int) []interface{} {
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]

        return []interface{}{
            courseID,
            userID,
            gofakeit.Number(1, 5),
            gofakeit.Sentence(10),
            gofakeit.Date(),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCertificate(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"certificate", "users", "course"}) != true {
        return nil
    }

    userIDs, err := getExistingIDs(db, "users", "account_id")
    if err != nil {
        return fmt.Errorf("failed to get user ids: %v", err)
    }

    courseIDs, err := getExistingIDs(db, "course", "course_id")
    if err != nil {
        return fmt.Errorf("failed to get course ids: %v", err)
    }

    query := `INSERT INTO certificate (user_id, course_id, issue_date)`

    dataFunc := func(i int) []interface{} {
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]

        return []interface{}{
            userID,
            courseID,
            gofakeit.Date(),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}
