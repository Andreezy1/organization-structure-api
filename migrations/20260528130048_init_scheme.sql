-- +goose Up
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_parent_department
        FOREIGN KEY(parent_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_department_name_per_parent
ON departments(name, parent_id);

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_employee_department
        FOREIGN KEY(department_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE employees;
DROP TABLE departments;
