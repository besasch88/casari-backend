ALTER TABLE "ceng_course_selection" DROP CONSTRAINT "fk_ceng_course_selection_menu_option";
ALTER TABLE "ceng_course_selection" DROP CONSTRAINT "fk_ceng_course_selection_menu_item";
ALTER TABLE "ceng_course_selection" DROP CONSTRAINT "fk_ceng_course_selection_course";
DROP TABLE IF EXISTS "ceng_course_selection";

ALTER TABLE "ceng_course" DROP CONSTRAINT "fk_ceng_course_order";
DROP TABLE IF EXISTS "ceng_course";


ALTER TABLE "ceng_order" DROP CONSTRAINT "fk_ceng_order_table";
DROP TABLE IF EXISTS "ceng_order";
