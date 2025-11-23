-- +goose Up
CREATE TABLE companies (
    ID UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL,
    description VARCHAR(3000),
    employee_count INT NOT NULL,
    registered BOOLEAN NOT NULL,
    company_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    created_by UUID REFERENCES users(ID),
    updated_by UUID REFERENCES users(ID),
    CONSTRAINT company_type_check CHECK (
        company_type IN (
            'Corporation',
            'NonProfit',
            'Cooperative',
            'Sole Proprietorship'
        )
    )
);

-- +goose Down
DROP TABLE companies;
