-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
	id BIGSERIAL PRIMARY KEY,
	service_name VARCHAR(255) NOT NULL,
	price INT NOT NULL,
	user_id uuid NOT NULL,
	start_date DATE NOT NULL DEFAULT now(),
	end_date DATE
);
-- +goose Down
DROP TABLE IF EXISTS subscriptions;