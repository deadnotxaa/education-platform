CREATE TABLE IF NOT EXISTS difficulty_level (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS course_specialization (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS course (
    course_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    specialization_id INTEGER REFERENCES course_specialization(id),
    duration INTEGER,
    price INTEGER,
    difficulty_level_id INTEGER REFERENCES difficulty_level(id)
);

CREATE TABLE IF NOT EXISTS course_topic (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    technologies TEXT,
    labor_intensity_hours SMALLINT,
    projects_number SMALLINT
);

CREATE TABLE IF NOT EXISTS course_topic_association (
    course_id INTEGER REFERENCES course(course_id),
    topic_id INTEGER REFERENCES course_topic(id),
    PRIMARY KEY (course_id, topic_id)
);

CREATE TABLE IF NOT EXISTS project (
    project_id SERIAL PRIMARY KEY,
    topic_id INTEGER REFERENCES course_topic(id),
    name VARCHAR(255) NOT NULL,
    description TEXT
);
