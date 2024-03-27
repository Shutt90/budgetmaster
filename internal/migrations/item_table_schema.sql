CREATE TABLE IF NOT EXISTS item (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    description VARCHAR(2550) NOT NULL DEFAULT '',
    location VARCHAR(255) NOT NULL DEFAULT'',
    cost BIGINT NULL NOT NULL DEFAULT 0,
    month ENUM('January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'Decemeber'),
    isMonthly TINYINT NOT NULL DEFAULT 0,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT NULL
);
