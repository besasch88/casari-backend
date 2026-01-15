ALTER TABLE "ceng_auth_session" DROP CONSTRAINT "idx_ceng_auth_session_refresh_token";

DROP TABLE IF EXISTS "ceng_auth_session";