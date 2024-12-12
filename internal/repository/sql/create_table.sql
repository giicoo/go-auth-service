CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    email varchar(255),
    hash_password varchar(255)
)