CREATE TABLE "ceng_order" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "table_id" VARCHAR(36) NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_order"
ADD CONSTRAINT "fk_ceng_order_table"
FOREIGN KEY ("table_id")
REFERENCES "ceng_table" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;


CREATE TABLE "ceng_course" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "order_id" VARCHAR(36) NOT NULL,
    "position" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_course"
ADD CONSTRAINT "fk_ceng_course_order"
FOREIGN KEY ("order_id")
REFERENCES "ceng_order" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;


CREATE TABLE "ceng_course_selection" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "course_id" VARCHAR(36) NOT NULL,
    "menu_item_id" VARCHAR(36) NOT NULL,
    "menu_option_id" VARCHAR(36),
    "quantity" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "ceng_course_selection"
ADD CONSTRAINT "fk_ceng_course_selection_course"
FOREIGN KEY ("course_id")
REFERENCES "ceng_course" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;


ALTER TABLE "ceng_course_selection"
ADD CONSTRAINT "fk_ceng_course_selection_menu_item"
FOREIGN KEY ("menu_item_id")
REFERENCES "ceng_menu_item" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE "ceng_course_selection"
ADD CONSTRAINT "fk_ceng_course_selection_menu_option"
FOREIGN KEY ("menu_option_id")
REFERENCES "ceng_menu_option" ("id")
ON DELETE SET NULL
ON UPDATE CASCADE;