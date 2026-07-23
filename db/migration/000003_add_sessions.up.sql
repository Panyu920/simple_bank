CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "user_ip" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "is_blocked" bool NOT NULL DEFAULT (false),
  "refresh_token_expired_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;