CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login TEXT UNIQUE NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    avatar TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for user table
CREATE INDEX IF NOT EXISTS idx_user_phone ON user(phone);
CREATE INDEX IF NOT EXISTS idx_user_login ON user(login);
CREATE INDEX IF NOT EXISTS idx_user_created_at ON user(created_at);

CREATE TABLE IF NOT EXISTS sms_code (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phone TEXT NOT NULL,
    code TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for sms_code table
CREATE INDEX IF NOT EXISTS idx_sms_code_phone ON sms_code(phone);
CREATE INDEX IF NOT EXISTS idx_sms_code_created_at ON sms_code(created_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sms_code_phone_unique ON sms_code(phone);

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

-- Indexes for user_profile table
CREATE INDEX IF NOT EXISTS idx_user_profile_private ON user_profile(private_profile);
CREATE INDEX IF NOT EXISTS idx_user_profile_status ON user_profile(status);
CREATE INDEX IF NOT EXISTS idx_user_profile_updated_at ON user_profile(updated_at);

CREATE TABLE IF NOT EXISTS user_relation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    related_user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    relation_type TEXT CHECK(relation_type IN ('friend', 'relation')) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for user_relation table
CREATE INDEX IF NOT EXISTS idx_user_relation_user_id ON user_relation(user_id);
CREATE INDEX IF NOT EXISTS idx_user_relation_related_user_id ON user_relation(related_user_id);
CREATE INDEX IF NOT EXISTS idx_user_relation_type ON user_relation(relation_type);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_relation_unique ON user_relation(user_id, related_user_id);

CREATE TABLE IF NOT EXISTS habbit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    complition_in_period_count INTEGER NOT NULL,
    complition_period TEXT CHECK(complition_period IN ('hour', 'day', 'week','month', 'none')) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for habbit table
CREATE INDEX IF NOT EXISTS idx_habbit_author_id ON habbit(author_id);
CREATE INDEX IF NOT EXISTS idx_habbit_created_at ON habbit(created_at);
CREATE INDEX IF NOT EXISTS idx_habbit_period ON habbit(complition_period);

CREATE TABLE IF NOT EXISTS habbit_complition(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habbit_id INTEGER NOT NULL REFERENCES habbit(id) ON DELETE CASCADE,
    complition TEXT CHECK(complition IN ('complete','not_complete')) NOT NULL,
    complited_at DATETIME NOT NULL,
    note TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for habbit_complition table
CREATE INDEX IF NOT EXISTS idx_habbit_complition_habbit_id ON habbit_complition(habbit_id);
CREATE INDEX IF NOT EXISTS idx_habbit_complition_complited_at ON habbit_complition(complited_at);
CREATE INDEX IF NOT EXISTS idx_habbit_complition_status ON habbit_complition(complition);

CREATE TABLE IF NOT EXISTS habbit_relation(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    habbit_id INTEGER NOT NULL REFERENCES habbit(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE
);

-- Indexes for habbit_relation table
CREATE INDEX IF NOT EXISTS idx_habbit_relation_habbit_id ON habbit_relation(habbit_id);
CREATE INDEX IF NOT EXISTS idx_habbit_relation_user_id ON habbit_relation(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_habbit_relation_unique ON habbit_relation(habbit_id, user_id);

CREATE TABLE IF NOT EXISTS notifications(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    message TEXT,
    habbit_id INTEGER REFERENCES habbit(id) ON DELETE SET NULL,
    user_id INTEGER REFERENCES user(id) ON DELETE SET NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for notifications table
CREATE INDEX IF NOT EXISTS idx_notifications_habbit_id ON notifications(habbit_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

CREATE TABLE IF NOT EXISTS notifications_views(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    habbit_id INTEGER NOT NULL REFERENCES habbit(id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for notifications_views table
CREATE INDEX IF NOT EXISTS idx_notifications_views_user_id ON notifications_views(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_views_habbit_id ON notifications_views(habbit_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_notifications_views_unique ON notifications_views(user_id, habbit_id);
