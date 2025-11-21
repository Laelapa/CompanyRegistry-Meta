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

-- name: GetCompanyByID :one
SELECT *
FROM companies
WHERE ID = $1;

-- name: UpdateCompany :exec
UPDATE companies
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    employee_count = COALESCE($4, employee_count),
    registered = COALESCE($5, registered),
    company_type = COALESCE($6, company_type),
    updated_at = CURRENT_TIMESTAMP,
    updated_by = $7
WHERE ID = $1;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE ID = $1;
