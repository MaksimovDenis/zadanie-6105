-- name: GetTenders :many
SELECT id, name, description, status, service_type, version, created_at
FROM tender
WHERE status = 'Published'
  AND ($1::text[] IS NULL OR array_length($1::text[], 1) = 0 OR service_type = ANY($1::text[]))
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: CreateTender :one
INSERT INTO tender (
    name,
    description,
    service_type,
    status,
    organization_id,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, (
        SELECT id FROM employee WHERE username = $6
    )
) RETURNING id, name, description, service_type, status, version, created_at;

-- name: CreateTenderHistory :exec
INSERT INTO tender_history (
    tender_id, 
    name, 
    description, 
    service_type, 
    status, 
    version, 
    created_at
)
SELECT id, name, description, service_type, status, version, created_at
FROM tender
WHERE tender.id = $1;

-- name: GetUserTenders :many
SELECT t.id, t.name, t.description, t.status, t.service_type, t.version, t.created_at
FROM tender AS t
INNER JOIN employee AS e ON t.created_by = e.id
WHERE e.username = $1
ORDER BY t.name ASC
LIMIT $2 OFFSET $3;

-- name: GetTenderStatus :one
SELECT t.status
FROM tender AS t
LEFT JOIN employee AS e ON t.created_by = e.id
WHERE (t.id = $1 AND e.username = $2) OR (t.id = $1 AND t.status = 'Published')
ORDER BY t.name ASC
LIMIT 1;

-- name: UpdateTenderStatus :one
UPDATE tender
SET status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE tender.id = $2
    AND EXISTS (
        SELECT 1
        FROM employee AS e
        JOIN organization_responsible AS ors ON e.id = ors.user_id
        WHERE e.username = $3
            AND ors.organization_id = (SELECT tender.organization_id FROM tender WHERE tender.id = $2)
    )
RETURNING tender.id, name, description, status, service_type, version, created_at;

-- name: EditTender :one
UPDATE tender t
SET name = COALESCE($1, t.name),
    description = COALESCE($2, t.description),
    service_type = COALESCE($3, t.service_type),
    version = version + 1
FROM employee e
JOIN organization_responsible ors ON e.id = ors.user_id
WHERE t.id = $4
    AND e.username = $5
    AND ors.organization_id = (SELECT organization_id FROM tender WHERE id = t.id)
RETURNING t.id, t.name, t.description, t.service_type, t.status, t.version, t.created_at;

-- name: RollbackTender :one
WITH
user_check AS (
    SELECT 1
    FROM employee e
    JOIN organization_responsible as ors ON e.id = ors.user_id
    JOIN tender t ON ors.organization_id = t.organization_id
    WHERE e.username = $3
      AND t.id = $1
),
current_tender AS (
    SELECT id, name, description, service_type, status, version, created_at, updated_at
    FROM tender
    WHERE id = $1
      AND version = (SELECT MAX(version) FROM tender WHERE id = $1)
),
previous_version AS (
    SELECT th.tender_id, th.name, th.description, th.service_type, th.status, th.version, th.created_at, updated_at
    FROM tender_history th
    WHERE th.tender_id = $1
      AND th.version = $2
),
updated_tender AS (
    UPDATE tender
    SET
        name = previous_version.name,
        description = previous_version.description,
        service_type = previous_version.service_type,
        status = previous_version.status,
        version = current_tender.version + 1,
        updated_at = CURRENT_TIMESTAMP
    FROM previous_version, current_tender
    WHERE tender.id = $1
      AND current_tender.version = (SELECT MAX(version) FROM tender WHERE id = $1)
      AND current_tender.id = previous_version.tender_id
      AND EXISTS (SELECT 1 FROM user_check)
    RETURNING tender.*
),
insert_history AS (
    INSERT INTO tender_history (tender_id, name, description, service_type, status, version, created_at, updated_at)
    SELECT id, name, description, service_type, status, version, created_at, updated_at
    FROM current_tender
    RETURNING *
)
SELECT * FROM updated_tender;