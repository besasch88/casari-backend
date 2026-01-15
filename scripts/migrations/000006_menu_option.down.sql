ALTER TABLE "ceng_menu_option" DROP CONSTRAINT "idx_ceng_menu_option_title";
ALTER TABLE "ceng_menu_option" DROP CONSTRAINT "fk_ceng_menu_option_menu_item";

DROP TABLE IF EXISTS "ceng_menu_option";
