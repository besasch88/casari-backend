ALTER TABLE "ceng_menu_item" DROP CONSTRAINT "idx_ceng_menu_item_title";
ALTER TABLE "ceng_menu_item" DROP CONSTRAINT "fk_ceng_menu_item_menu_category";

DROP TABLE IF EXISTS "ceng_menu_item";
