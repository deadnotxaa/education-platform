package versions

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"log"
	"seeder/config"
	"seeder/internal/db_utils"
)

func SeedV2(seedCfg *config.SeedConfig, db *db_utils.DB) error {
    version := 2

    err := SeedTable(db.Conn, version, "users", func() error {
        return seedUser(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
		log.Printf("Seeding user table completed successfully")
	}

    err = SeedTable(db.Conn, version, "role", func() error {
        return seedRole(db.Conn)
    })
    if err != nil {
        return err
    } else {
		log.Printf("Seeding role table completed successfully")
	}

    err = SeedTable(db.Conn, version, "employee", func() error {
        return seedEmployee(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
		log.Printf("Seeding employee table completed successfully")
	}

    err = SeedTable(db.Conn, version, "teacher", func() error {
        return seedTeacher(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
		log.Printf("Seeding teacher table completed successfully")
	}

    err = SeedTable(db.Conn, version, "course_teacher", func() error {
        return seedCourseTeacher(db.Conn, seedCfg.SeedCount, seedCfg.InsertBatchSize)
    })
    if err != nil {
        return err
    } else {
		log.Printf("Seeding course_teacher table completed successfully")
	}

    return nil
}

func seedUser(db *sql.DB, seedCount int, insertBatchSize int) error {
	if tablesExists(db, []string{"users"}) != true {
		return nil
	}

	query := `INSERT INTO users (name, surname, birthdate, email, hashed_password, profile_picture_url, phone_number, snils_number)`

	dataFunc := func(i int) []interface{} {
		h := sha256.New()
		h.Write([]byte(gofakeit.Password(true, true, true, true, false, 12)))
		hashedPassword := fmt.Sprintf("%x", h.Sum(nil))

		return []interface{}{
			gofakeit.Name(),
			gofakeit.LastName(),
			gofakeit.Date(),
			gofakeit.Email(),
			hashedPassword,
			gofakeit.URL(),
			gofakeit.Phone(),
			gofakeit.SSN(),
		}
	}

	return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedRole(db *sql.DB) error {
	if tablesExists(db, []string{"role"}) != true {
		return nil
	}

	query := `INSERT INTO role (name) VALUES ($1)`

	roles := []string{"admin", "teacher", "student", "smm manager", "mentor", "technical support"}

	for _, role := range roles {
		_, err := db.Exec(query, role)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedEmployee(db *sql.DB, seedCount int, insertBatchSize int) error {
	if tablesExists(db, []string{"employee"}) != true {
		return nil
	}

	userIDs, err := getExistingIDs(db, "users", "account_id")
	if err != nil {
		return fmt.Errorf("failed to get users ids: %v", err)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("no user records found to associate with employees")
	}

	query := `INSERT INTO employee (user_id, role_id)`

	dataFunc := func(i int) []interface{} {
		userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		roleID := gofakeit.Number(1, 6)

		return []interface{}{
			userID,
			roleID,
		}
	}

	return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedTeacher(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"teacher"}) != true {
        return nil
    }

    employeeIDs, err := getExistingIDs(db, "employee")
    if err != nil {
        return fmt.Errorf("failed to get employee ids: %v", err)
    }

    if len(employeeIDs) == 0 {
        return fmt.Errorf("no employee records found to associate with teachers")
    }

    query := `INSERT INTO teacher (employee_id, work_place, overall_experience, specialization_experience)`

    dataFunc := func(i int) []interface{} {
        employeeID := employeeIDs[gofakeit.Number(0, len(employeeIDs)-1)]

        return []interface{}{
            employeeID,
            gofakeit.Company(),
            gofakeit.Number(1, 30),
            gofakeit.Number(1, 20),
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}

func seedCourseTeacher(db *sql.DB, seedCount int, insertBatchSize int) error {
    if tablesExists(db, []string{"course_teacher"}) != true {
        return nil
    }

    teacherIDs, err := getExistingIDs(db, "teacher")
    if err != nil {
        return fmt.Errorf("failed to get teacher IDs: %v", err)
    }

    courseIDs, err := getExistingIDs(db, "course")
    if err != nil {
        return fmt.Errorf("failed to get course IDs: %v", err)
    }

    query := `INSERT INTO course_teacher (teacher_id, course_id)`

    dataFunc := func(i int) []interface{} {
        teacherID := teacherIDs[gofakeit.Number(0, len(teacherIDs)-1)]
        courseID := courseIDs[gofakeit.Number(0, len(courseIDs)-1)]

        return []interface{}{
            teacherID,
            courseID,
        }
    }

    return BatchInsertData(db, query, insertBatchSize, dataFunc, seedCount)
}
