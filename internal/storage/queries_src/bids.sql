-- name: CreateBid :one
WITH user_check AS (
    SELECT 1
    FROM employee e
    JOIN organization_responsible ors ON e.id = ors.user_id
    WHERE e.username = $7
      AND ors.organization_id = $5
),
new_bid AS (
    INSERT INTO bids (
        name,
        description,
        status,
        tender_id,
        organization_id,
        user_id,
        author_type
    ) VALUES (
        $1, $2, $3, $4, $5, (
            SELECT id FROM employee WHERE username = $7
        ), $6
    ) 
    RETURNING id, name, description, status, organization_id, user_id, author_type, tender_id, version, created_at
)
SELECT * FROM new_bid
WHERE EXISTS (SELECT 1 FROM user_check);

-- name: CreateBidsHistory :exec
INSERT INTO bids_history (
    bid_id,
    name,
    description,
    status,
    tender_id,
    organization_id,
    user_id,
    author_type,
    version,
    created_at
)
SELECT 
  id, 
  name, 
  description, 
  status, 
  tender_id, 
  organization_id, 
  user_id, 
  author_type,
  version, 
  created_at
FROM bids
WHERE bids.id = $1;

-- name: GetUserBids :many
SELECT 
  b.id, 
  b.name, 
  b.description, 
  b.status, 
  b.tender_id, 
  b.organization_id, 
  b.user_id, 
  b.author_type,
  b.tender_id, 
  b.version, 
  b.created_at
FROM bids AS b
INNER JOIN employee AS e ON b.user_id = e.id
WHERE e.username = $1
ORDER BY b.name ASC
LIMIT $2 OFFSET $3;

-- name: GetBidsForTender :many
SELECT 
  b.id, 
  b.name, 
  b.description, 
  b.status, 
  b.author_type, 
  b.user_id, 
  b.tender_id,
  b.version, 
  b.created_at
FROM bids AS b
WHERE b.tender_id = $1
  AND EXISTS (
      SELECT 1
      FROM employee e
      JOIN organization_responsible ors ON e.id = ors.user_id
      WHERE e.username = $2
        AND (b.user_id = e.id OR b.organization_id = ors.organization_id)
  )
ORDER BY b.name ASC
LIMIT $3 OFFSET $4;

-- name: GetBidStatus :one
SELECT b.status
FROM bids b
JOIN employee e ON b.user_id = e.id
WHERE b.id = $1
  AND e.username = $2
UNION
SELECT b.status
FROM bids b
JOIN organization_responsible ors ON b.organization_id = ors.organization_id
JOIN employee e ON ors.user_id = e.id
WHERE b.id = $1
  AND e.username = $2;

-- name: UpdateBidStatus :one
UPDATE bids
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE bids.id = $1
  AND (
    EXISTS (
        SELECT 1
        FROM employee e
        WHERE e.username = $3
          AND e.id = (SELECT user_id FROM bids WHERE id = $1)
    )
    OR
    EXISTS (
        SELECT 1
        FROM organization_responsible ors
        JOIN employee e ON ors.user_id = e.id
        WHERE e.username = $3
          AND ors.organization_id = (SELECT organization_id FROM bids WHERE id = $1)
    )
)
RETURNING bids.id, name, description, status, user_id, author_type, tender_id, version, created_at, updated_at;


-- name: EditBid :one
UPDATE bids AS b
SET
    name = COALESCE($2, b.name),
    description = COALESCE($3, b.description),
    updated_at = CURRENT_TIMESTAMP,
    version = version + 1
WHERE b.id = $1
  AND (
    (user_id = (SELECT e.id FROM employee e WHERE e.username = $4))
    OR
    (organization_id = (SELECT ors.organization_id FROM organization_responsible ors JOIN employee e ON ors.user_id = e.id WHERE e.username = $4))
  )
RETURNING b.id, b.name, b.description, b.status, b.user_id, b.author_type, b.tender_id, b.version, b.created_at;

-- name: SubmitBidDecision :one
WITH user_check AS (
    SELECT EXISTS (
        SELECT 1
        FROM employee e
        JOIN organization_responsible ors ON e.id = ors.user_id
        JOIN tender t ON t.id = (SELECT tender_id FROM bids WHERE bids.id = $1)
        WHERE e.username = $3
          AND ors.organization_id = t.organization_id
    ) AS has_permission
),
update_bid AS (
    UPDATE bids
    SET
        status = CASE
            WHEN $2::text = 'Rejected' THEN 'Rejected'
            WHEN $2::text = 'Approved'
                 AND (SELECT COUNT(*) FROM bids_decisions WHERE bid_id = $1 AND decision = 'Approved') >= LEAST(3, (SELECT COUNT(DISTINCT user_id) 
                 FROM organization_responsible WHERE organization_id = (SELECT organization_id FROM bids WHERE id = $1))) THEN 'Approved'
            ELSE status
        END,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = $1
      AND EXISTS (SELECT 1 FROM user_check WHERE has_permission)
    RETURNING id, name, description, status, organization_id, user_id, author_type, tender_id, version, created_at
),
tender_update AS (
    UPDATE tender
    SET status = 'Closed'
    WHERE id = (SELECT tender_id FROM bids WHERE id = $1)
      AND EXISTS (SELECT 1 FROM update_bid WHERE status = 'Approved')
    RETURNING id, status
),
insert_history AS (
    INSERT INTO bids_history (bid_id, name, description, status, version, created_at, updated_at)
    SELECT id, name, description, status, version, created_at, updated_at
    FROM update_bid
    RETURNING id, name, description, status, version, created_at, updated_at
)
SELECT * FROM update_bid;


-- name: SubmitBidFeedback :one
WITH user_check AS (
    SELECT EXISTS (
        SELECT 1
        FROM employee e
        JOIN organization_responsible ors ON e.id = ors.user_id
        JOIN tender t ON t.id = (SELECT tender_id FROM bids WHERE id = $1)
        WHERE e.username = $3
          AND ors.organization_id = t.organization_id
    ) AS has_permission
),
insert_feedback AS (
    INSERT INTO bid_feedbacks (bid_id, user_id, feedback, created_at)
    SELECT $1, e.id, $2, CURRENT_TIMESTAMP
    FROM employee e
    WHERE e.username = $3
      AND (SELECT has_permission FROM user_check) = true
    RETURNING id AS feedback_id
)
SELECT b.id, b.name, b.description, b.status, b.user_id, b.author_type, b.version, b.tender_id, b.created_at, b.updated_at
FROM bids b
WHERE b.id = $1;

-- name: RollbackBid :one
WITH 
user_check AS (
    SELECT 1
    FROM employee e
    JOIN organization_responsible ors ON e.id = ors.user_id
    JOIN bids b ON ors.organization_id = b.organization_id
    WHERE e.username = $3
      AND b.id = $1
),
current_bid AS (
    SELECT b.id, b.tender_id, b.organization_id, b.user_id, b.author_type, b.name, b.description, b.status, b.version, b.created_at, b.updated_at
    FROM bids AS b
    WHERE b.id = $1
      AND b.version = (SELECT MAX(version) FROM bids WHERE id = $1)
),
previous_version AS (
    SELECT bh.bid_id, bh.tender_id, bh.organization_id, bh.user_id, bh.author_type, bh.name, bh.description, bh.status, bh.version, bh.created_at, bh.updated_at
    FROM bids_history bh
    WHERE bh.bid_id = $1
      AND bh.version = $2
),
updated_bid AS (
    UPDATE bids
    SET
        name = previous_version.name,
        description = previous_version.description,
        status = previous_version.status,
        version = current_bid.version + 1,
        updated_at = CURRENT_TIMESTAMP
    FROM previous_version, current_bid
    WHERE bids.id = $1
      AND current_bid.version = (SELECT MAX(version) FROM bids WHERE id = $1)
      AND current_bid.id = bids.id
      AND EXISTS (SELECT 1 FROM user_check)
    RETURNING bids.*
),
insert_history AS (
    INSERT INTO bids_history (bid_id, tender_id, organization_id, user_id, author_type, name, description, status, version, created_at, updated_at)
    SELECT id, tender_id, organization_id, user_id, author_type, name, description, status, version, created_at, updated_at
    FROM current_bid
    RETURNING *
)
SELECT * FROM updated_bid;



-- name: GetBidReviews :many
SELECT r.id, r.bid_id, r.user_id, r.feedback, r.created_at
FROM bid_feedbacks r
JOIN bids b ON r.bid_id = b.id
WHERE b.tender_id = $1
  AND b.user_id = (
    SELECT id
    FROM employee
    WHERE employee.username = $2
  )
  AND EXISTS (
    SELECT 1
    FROM organization_responsible ors
    JOIN employee e ON ors.user_id = e.id
    WHERE e.username = $3
      AND ors.organization_id = b.organization_id
)
ORDER BY r.created_at DESC
LIMIT $4 OFFSET $5;