ALTER TABLE "ceng_menu_category" DROP CONSTRAINT "idx_ceng_menu_category_title";
ALTER TABLE "ceng_menu_category" DROP CONSTRAINT "fk_ceng_menu_category_printer";

DROP TABLE IF EXISTS "ceng_menu_category";
