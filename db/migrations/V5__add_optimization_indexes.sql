-- Migration V5__add_optimization_indexes.sql

-- Indexes for query1
-- This query joins course, purchase, user, and employee tables
CREATE INDEX idx_course_course_id ON course(course_id);
CREATE INDEX idx_purchase_course_id ON purchase(course_id);
CREATE INDEX idx_purchase_user_id ON purchase(user_id);
CREATE INDEX idx_purchase_date ON purchase(purchase_date);
CREATE INDEX idx_user_account_id ON "user"(account_id);
CREATE INDEX idx_employee_user_id ON employee(user_id);

-- Indexes for query2
-- This query additionally uses the role table
CREATE INDEX idx_employee_role_id ON employee(role_id);
CREATE INDEX idx_role_name ON role(name);
CREATE INDEX idx_purchase_date_desc ON purchase(purchase_date DESC);

-- Indexes for query3
-- This query additionally uses the course_specialization table
CREATE INDEX idx_course_specialization_id ON course(specialization_id);
CREATE INDEX idx_course_specialization_name ON course_specialization(name);

-- Composite indexes for commonly joined columns
CREATE INDEX idx_purchase_course_user ON purchase(course_id, user_id, purchase_date);
CREATE INDEX idx_employee_user_role ON employee(user_id, role_id);

-- Create a functional index for course_specialization filtering
CREATE INDEX idx_course_specialization_programming
ON course_specialization((name = 'Programming'))
WHERE name = 'Programming';