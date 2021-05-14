create table Employee
(
    id          INTEGER NOT NULL AUTOINCREMENT,
    first_name  VARCHAR(20),
    last_name   VARCHAR(20),
    birthday    VARCHAR,
    salary      BIGINT,
    gender      INTEGER,
    status      INTEGER,
    created_at  VARCHAR,
    created_by  VARCHAR(20),
    updated_at  VARCHAR,
    updated_by  VARCHAR(20),
    PRIMARY KEY (id)
)

