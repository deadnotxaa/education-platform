CREATE TABLE IF NOT EXISTS users (
    account_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    surname VARCHAR(255),
    birthdate DATE,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    profile_picture_url VARCHAR(255),
    phone_number VARCHAR(20),
    snils_number VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS role (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(account_id),
    role_id INTEGER REFERENCES role(id)
);

CREATE TABLE IF NOT EXISTS teacher (
    employee_id INTEGER PRIMARY KEY REFERENCES employee(id),
    work_place VARCHAR(255),
    overall_experience INTEGER,
    specialization_experience INTEGER
);

CREATE TABLE IF NOT EXISTS course_teacher (
    teacher_id INTEGER REFERENCES teacher(employee_id),
    course_id INTEGER REFERENCES course(course_id),
    PRIMARY KEY (teacher_id, course_id)
);
