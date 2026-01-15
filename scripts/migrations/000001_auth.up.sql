CREATE TABLE "ceng_auth_session" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "user_id" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "expires_at" TIMESTAMP NOT NULL,
    "refresh_token" TEXT NOT NULL
);

ALTER TABLE "ceng_auth_session" ADD CONSTRAINT "idx_ceng_auth_session_refresh_token" UNIQUE ("refresh_token");