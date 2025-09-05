-- Indexes for query1
CREATE INDEX idx_course_course_id ON course(course_id);
CREATE INDEX idx_purchase_course_id ON purchase(course_id);
CREATE INDEX idx_purchase_user_id ON purchase(user_id);
CREATE INDEX idx_purchase_date ON purchase(purchase_date);
CREATE INDEX idx_user_account_id ON users(account_id);
CREATE INDEX idx_employee_user_id ON employee(user_id);

-- Indexes for query2
CREATE INDEX idx_employee_role_id ON employee(role_id);
CREATE INDEX idx_role_name ON role(name);
CREATE INDEX idx_purchase_date_desc ON purchase(purchase_date DESC);

-- Indexes for query3
CREATE INDEX idx_course_specialization_id ON course(specialization_id);
CREATE INDEX idx_course_specialization_name ON course_specialization(name);

-- Composite indexes
CREATE INDEX idx_purchase_course_user ON purchase(course_id, user_id, purchase_date);
CREATE INDEX idx_employee_user_role ON employee(user_id, role_id);

-- Functional index
CREATE INDEX idx_course_specialization_programming
ON course_specialization((name = 'Programming'))
WHERE name = 'Programming';