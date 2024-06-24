CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY,
    firstName VARCHAR(255) NOT NULL DEFAULT '',
    surname VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT'',
    password VARCHAR(255) NOT NULL DEFAULT '',
    roles text[] DEFAULT [],
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT NULL,
    UNIQUE(email)
);
