CREATE TABLE IF NOT EXISTS course_calendar (
    id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES course(course_id),
    start_date DATE,
    end_sales_date DATE,
    remaining_places SMALLINT
);

CREATE TABLE IF NOT EXISTS course_type (
    id SMALLSERIAL PRIMARY KEY,
    type_name VARCHAR(255) NOT NULL,
    discount SMALLINT
);

CREATE TABLE IF NOT EXISTS purchase (
    purchase_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(account_id),
    course_id INTEGER REFERENCES course(course_id),
    purchase_date TIMESTAMP,
    course_type_id SMALLINT REFERENCES course_type(id),
    total_price INTEGER,
    purchase_status VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS career_center_student (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(account_id),
    course_id INTEGER REFERENCES course(course_id),
    cv_url VARCHAR(255),
    career_support_start DATE,
    support_period INTEGER
);

CREATE TABLE IF NOT EXISTS partner_company (
    company_id SERIAL PRIMARY KEY,
    short_name VARCHAR(255),
    full_name VARCHAR(255),
    hired_graduates_count INTEGER,
    requirements TEXT,
    agreement_status BOOLEAN
);

CREATE TABLE IF NOT EXISTS job_application (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES career_center_student(id),
    company_id INTEGER REFERENCES partner_company(company_id),
    application_date DATE,
    status VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS blog_post (
    post_id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES employee(id),
    title VARCHAR(255),
    publication_date TIMESTAMP,
    topic VARCHAR(255),
    reading_time_minutes INTEGER,
    cover_image_url VARCHAR(255),
    content TEXT
);

CREATE TABLE IF NOT EXISTS tag (
    tag_id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS post_tag (
    post_id INTEGER REFERENCES blog_post(post_id),
    tag_id INTEGER REFERENCES tag(tag_id),
    PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE IF NOT EXISTS course_review (
    review_id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES course(course_id),
    user_id INTEGER REFERENCES users(account_id),
    rating SMALLINT,
    comment TEXT,
    review_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS certificate (
    certificate_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(account_id),
    course_id INTEGER REFERENCES course(course_id),
    issue_date DATE
);
