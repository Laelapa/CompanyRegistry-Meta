-- name: CreateCompany :one
INSERT INTO companies (
    name,
    description,
    employee_count,
    registered,
    company_type,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetCompanyByName :one
SELECT *
FROM companies
WHERE name = $1;

-- name: UpdateCompany :one
UPDATE companies
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    employee_count = COALESCE(sqlc.narg('employee_count'), employee_count),
    registered = COALESCE(sqlc.narg('registered'), registered),
    company_type = COALESCE(sqlc.narg('company_type'), company_type),
    updated_at = CURRENT_TIMESTAMP,
    updated_by = sqlc.arg('updated_by')
WHERE ID = sqlc.arg('id')
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE ID = $1;
