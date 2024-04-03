CREATE TABLE IF NOT EXISTS item (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    description VARCHAR(2550) NOT NULL DEFAULT '',
    location VARCHAR(255) NOT NULL DEFAULT'',
    cost BIGINT NULL NOT NULL DEFAULT 0,
    month VARCHAR(255) CHECK(month IN('January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December')),
    year INT NOT NULL DEFAULT 2024,
    isRecurring TINYINT NOT NULL DEFAULT 0,
    removedReccuringAt DATETIME DEFAULT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT NULL
);
