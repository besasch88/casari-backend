CREATE TABLE "ceng_menu_option" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "menu_item_id" VARCHAR(36) NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "position" BIGINT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "price" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_menu_option"
ADD CONSTRAINT "fk_ceng_menu_option_menu_item"
FOREIGN KEY ("menu_item_id")
REFERENCES "ceng_menu_item" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE "ceng_menu_option" ADD CONSTRAINT "idx_ceng_menu_option_title" UNIQUE ("title");