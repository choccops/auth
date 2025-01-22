-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
	id SERIAL PRIMARY KEY,
	username varchar(50) NOT NULL UNIQUE,
	password varchar(100) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `users`;
-- +goose StatementEnd
