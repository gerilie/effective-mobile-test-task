-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	service_name VARCHAR(255) NOT NULL,
	price INT NOT NULL,
	user_id uuid NOT NULL,
	start_date DATE NOT NULL DEFAULT now(),
	end_date DATE,
	CONSTRAINT idx_subscriptions_unique UNIQUE (user_id, service_name)
);
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
-- +goose Down
DROP TABLE IF EXISTS subscriptions;