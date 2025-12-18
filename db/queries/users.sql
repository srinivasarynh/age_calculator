INSERT INTO users (name, dob)
VALUES ($1, $2)
RETURNING id, name, dob, created_at, updated_at;

SELECT id, name, dob, created_at, updated_at
FROM users
WHERE id = $1;

SELECT id, name, dob, created_at, updated_at
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

UPDATE users
SET name = $1, dob = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3
RETURNING id, name, dob, created_at, updated_at;

DELETE FROM users
WHERE id = $1;

SELECT COUNT(*) FROM users;
