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
        return seedCourseCalendar(db.Conn, seedCfg.SeedCount)
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
        return seedPurchase(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding purchase table completed successfully")
    }

    err = SeedTable(db.Conn, version, "career_center_student", func() error {
        return seedCareerCenterStudent(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding career_center_student table completed successfully")
    }

    err = SeedTable(db.Conn, version, "partner_company", func() error {
        return seedPartnerCompany(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding partner_company table completed successfully")
    }

    err = SeedTable(db.Conn, version, "job_application", func() error {
        return seedJobApplication(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding job_application table completed successfully")
    }

    err = SeedTable(db.Conn, version, "blog_post", func() error {
        return seedBlogPost(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding blog_post table completed successfully")
    }

    err = SeedTable(db.Conn, version, "tag", func() error {
        return seedTag(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding tag table completed successfully")
    }

    err = SeedTable(db.Conn, version, "post_tag", func() error {
        return seedPostTag(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding post_tag table completed successfully")
    }

    err = SeedTable(db.Conn, version, "course_review", func() error {
        return seedCourseReview(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding course_review table completed successfully")
    }

    err = SeedTable(db.Conn, version, "certificate", func() error {
        return seedCertificate(db.Conn, seedCfg.SeedCount)
    })
    if err != nil {
        return err
    } else {
        log.Println("Seeding certificate table completed successfully")
    }

    return nil
}

func seedCourseCalendar(db *sql.DB, seedCount int) error {
    if tablesExists(db, []string{"course_calendar"}) != true {
        return nil
    }

    courseIDs, err := getExistingIDs(db, "course", "course_id")
    if err != nil {
        return fmt.Errorf("failed to get course ids: %v", err)
    }

    query := `INSERT INTO course_calendar (course_id, start_date, end_sales_date, remaining_places)
              VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        startDate := gofakeit.Date()
        endSalesDate := gofakeit.DateRange(startDate, gofakeit.Date())
        remainingPlaces := gofakeit.Number(0, 100)

        _, err = db.Exec(query, courseID, startDate, endSalesDate, remainingPlaces)
        if err != nil {
            return fmt.Errorf("failed to insert course_calendar: %v", err)
        }
    }

    return nil
}

func seedCourseType(db *sql.DB) error {
    if tablesExists(db, []string{"course_type"}) != true {
        return nil
    }

    query := `INSERT INTO course_type (type_name, discount)
              VALUES ($1, $2) ON CONFLICT DO NOTHING`

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

func seedPurchase(db *sql.DB, seedCount int) error {
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

    query := `INSERT INTO purchase (user_id, course_id, purchase_date, course_type_id, total_price, purchase_status)
              VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        courseTypeID := courseTypeIDs[gofakeit.Number(0, len(courseTypeIDs)-1)]

        purchaseDate := gofakeit.Date()
        totalPrice := gofakeit.Number(1000, 5000)
        purchaseStatus := gofakeit.RandomString([]string{"Completed", "Pending", "Cancelled"})

        _, err = db.Exec(query, userID, courseID, purchaseDate, courseTypeID, totalPrice, purchaseStatus)
        if err != nil {
            return fmt.Errorf("failed to insert purchase: %v", err)
        }
    }

    return nil
}

func seedPartnerCompany(db *sql.DB, seedCount int) error {
    if tablesExists(db, []string{"partner_company"}) != true {
        return nil
    }

    query := `INSERT INTO partner_company (short_name, full_name, hired_graduates_count, requirements, agreement_status)
              VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        shortName := gofakeit.Company()
        fullName := gofakeit.Company() + " LLC"
        hiredGraduatesCount := gofakeit.Number(0, 100)
        requirements := gofakeit.Sentence(10)
        agreementStatus := gofakeit.Bool()

        _, err := db.Exec(query, shortName, fullName, hiredGraduatesCount, requirements, agreementStatus)
        if err != nil {
            return fmt.Errorf("failed to insert partner_company: %v", err)
        }
    }

    return nil
}

func seedCareerCenterStudent(db *sql.DB, seedCount int) error {
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

    query := `INSERT INTO career_center_student (user_id, course_id, cv_url, career_support_start, support_period)
              VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        cvURL := gofakeit.URL()
        careerSupportStart := gofakeit.Date()
        supportPeriod := gofakeit.Number(1, 12)

        _, err = db.Exec(query, userID, courseID, cvURL, careerSupportStart, supportPeriod)
        if err != nil {
            return fmt.Errorf("failed to insert career_center_student: %v", err)
        }
    }

    return nil
}

func seedJobApplication(db *sql.DB, seedCount int) error {
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

    query := `INSERT INTO job_application (student_id, company_id, application_date, status)
              VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        studentID := studentIDs[gofakeit.Number(0, len(studentIDs)-1)]
        companyID := companyIDs[gofakeit.Number(0, len(companyIDs)-1)]
        applicationDate := gofakeit.Date()
        status := gofakeit.RandomString([]string{"Applied", "Interviewed", "Rejected", "Hired"})

        _, err = db.Exec(query, studentID, companyID, applicationDate, status)
        if err != nil {
            return fmt.Errorf("failed to insert job_application: %v", err)
        }
    }

    return nil
}

func seedBlogPost(db *sql.DB, seedCount int) error {
    if tablesExists(db, []string{"blog_post", "employee"}) != true {
        return fmt.Errorf("failed to seed blog_post, some tables not found")
    }

    employeeIDs, err := getExistingIDs(db, "employee", "id")
    if err != nil {
        return fmt.Errorf("failed to get employee ids: %v", err)
    }

    query := `INSERT INTO blog_post (author_id, title, publication_date, topic, reading_time_minutes, cover_image_url, content)
              VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        authorID := employeeIDs[gofakeit.Number(0, len(employeeIDs)-1)]
        title := gofakeit.BookTitle()
        publicationDate := gofakeit.Date()
        topic := gofakeit.Sentence(3)
        readingTimeMinutes := gofakeit.Number(5, 30)
        coverImageURL := gofakeit.URL()
        content := gofakeit.Paragraph(5, 10, 100, "\n")

        _, err = db.Exec(query, authorID, title, publicationDate, topic, readingTimeMinutes, coverImageURL, content)
        if err != nil {
            return fmt.Errorf("failed to insert blog_post: %v", err)
        }
    }

    return nil
}

func seedTag(db *sql.DB, seedCount int) error {
    if tablesExists(db, []string{"tag"}) != true {
        return nil
    }

    query := `INSERT INTO tag (name) VALUES ($1) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        tagName := gofakeit.HackerAdjective() + " " + gofakeit.HackerNoun()
        _, err := db.Exec(query, tagName)
        if err != nil {
            return fmt.Errorf("failed to insert tag: %v", err)
        }
    }

    return nil
}

func seedPostTag(db *sql.DB, seedCount int) error {
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

    query := `INSERT INTO post_tag (post_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        postID := postIDs[gofakeit.Number(0, len(postIDs)-1)]
        tagID := tagIDs[gofakeit.Number(0, len(tagIDs)-1)]

        _, err = db.Exec(query, postID, tagID)
        if err != nil {
            return err
        }
    }

    return nil
}

func seedCourseReview(db *sql.DB, seedCount int) error {
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

    query := `INSERT INTO course_review (course_id, user_id, rating, comment, review_date)
              VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        rating := gofakeit.Number(1, 5)
        comment := gofakeit.Sentence(10)
        reviewDate := gofakeit.Date()

        _, err = db.Exec(query, courseID, userID, rating, comment, reviewDate)
        if err != nil {
            return fmt.Errorf("failed to insert course_review: %v", err)
        }
    }

    return nil
}

func seedCertificate(db *sql.DB, seedCount int) error {
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

    query := `INSERT INTO certificate (user_id, course_id, issue_date)
              VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

    for i := 0; i < seedCount; i++ {
        userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]
        issueDate := gofakeit.Date()

        _, err = db.Exec(query, userID, courseID, issueDate)
        if err != nil {
            return fmt.Errorf("failed to insert certificate: %v", err)
        }
    }

    return nil
}
