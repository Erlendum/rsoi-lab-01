-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS persons(
    id serial primary key,
    name text not null,
    age int not null,
    address text,
    work text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS persons;
-- +goose StatementEnd
