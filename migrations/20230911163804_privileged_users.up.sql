CREATE TABLE IF NOT EXISTS privileged_users (
    user_id INTEGER,
    privilege_id INTEGER,

    FOREIGN KEY (privilege_id) REFERENCES privileges(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, privilege_id)
);