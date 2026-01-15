CREATE TABLE "ceng_menu_category" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "position" BIGINT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "printer_id" VARCHAR(36),
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_menu_category"
ADD CONSTRAINT "fk_ceng_menu_category_printer"
FOREIGN KEY ("printer_id")
REFERENCES "ceng_printer" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE "ceng_menu_category" ADD CONSTRAINT "idx_ceng_menu_category_title" UNIQUE ("title");
