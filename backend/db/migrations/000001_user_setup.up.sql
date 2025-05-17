CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "email" VARCHAR(256) NOT NULL UNIQUE,
    "hashed_password" VARCHAR(256) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);