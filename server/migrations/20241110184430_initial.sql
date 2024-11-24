-- +goose Up
CREATE TABLE IF NOT EXISTS tasks(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    priority INTEGER NOT NULL,
    data TEXT NOT NULL,
    state VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    job_id INTEGER NOT NULL,
    FOREIGN KEY(job_id) REFERENCES jobs(id)
    );

CREATE TABLE IF NOT EXISTS workers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    state VARCHAR(10) NOT NULL,
    metadata TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    task_id INTEGER,
    FOREIGN KEY(task_id) REFERENCES tasks(id)
    );

CREATE TABLE IF NOT EXISTS jobs(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL,
    priority INTEGER NOT NULL,
    state VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS workers;
DROP TABLE IF EXISTS jobs;
