CREATE TABLE "ceng_course" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "table_id" VARCHAR(36) NOT NULL,
    "user_id" VARCHAR(36) NOT NULL,
    "close" BOOLEAN NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_course"
ADD CONSTRAINT "fk_ceng_course_table"
FOREIGN KEY ("table_id")
REFERENCES "ceng_table" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;