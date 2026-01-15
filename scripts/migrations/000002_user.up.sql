CREATE TABLE "ceng_user" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "username" VARCHAR(255) NOT NULL,
    "password" VARCHAR(1024) NOT NULL,
    "permissions" TEXT[] NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_user" ADD CONSTRAINT "idx_ceng_user_username" UNIQUE ("username");