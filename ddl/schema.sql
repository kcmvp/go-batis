CREATE TABLE Company
(
    company_id      INTEGER,
    company_name    VARCHAR(20),
    company_address VARCHAR(100),
    ss_no           int
);

CREATE TABLE Department
(
    dept_id   INTEGER,
    dept_name VARCHAR(30),
    status    BOOLEAN
);

CREATE TABLE Employee
(
    emp_id     INTEGER,
    first_name VARCHAR(30),
    last_name  VARCHAR(30),
    birthday   DATE,
    salary     DECIMAL(10, 2),
    gender     TINYINT,
    status     BOOLEAN
);
