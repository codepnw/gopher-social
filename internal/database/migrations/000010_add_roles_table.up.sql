CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, description, level)
VALUES 
    ('user', 'a user can create posts and comments', 1),
    ('staff', 'a staff can update other users posts', 2),
    ('admin', 'an admin can update and delete other users posts', 3);