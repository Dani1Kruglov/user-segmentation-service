CREATE TABLE users_segments (
    user_id INT,
    segment_id INT,
    created_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    duration VARCHAR(255) DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (segment_id) REFERENCES segments(id)
);