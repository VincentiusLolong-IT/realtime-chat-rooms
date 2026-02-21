-- name: GetRoom :one
SELECT * FROM rooms WHERE id = ? LIMIT 1;

-- name: ListRooms :many
SELECT * FROM rooms ORDER BY created_at DESC;

-- name: ListRoomsByUser :many
SELECT * FROM rooms WHERE created_by = ? ORDER BY created_at DESC;

-- name: CreateRoom :execresult
INSERT INTO rooms (name, created_by) VALUES (?, ?);

-- name: UpdateRoom :exec
UPDATE rooms SET name = ? WHERE id = ?;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = ?;

-- name: GetRoomsByUserID :many
SELECT * FROM rooms
WHERE created_by = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountRoomsByUserID :one
SELECT COUNT(*) FROM rooms
WHERE created_by = ?;