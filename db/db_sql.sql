CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    avatar TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sms_code (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phone TEXT NOT NULL,
    code TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_profile (
    user_id INTEGER PRIMARY KEY REFERENCES user(id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    middle_name TEXT,
    status TEXT,
    private_profile BOOLEAN DEFAULT FALSE,
    birth_date DATE,
    gender TEXT CHECK(gender IN ('male', 'female')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_relation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    related_user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    relation_type TEXT CHECK(relation_type IN ('friend', 'relation')) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sms_code_created_at
    ON sms_code (created_at);
