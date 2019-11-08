
DROP TABLE IF EXISTS "tbl_category";
DROP SEQUENCE IF EXISTS tbl_category_id_seq;
CREATE SEQUENCE tbl_category_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 4 CACHE 1;

CREATE TABLE "public"."tbl_category" (
    "id" integer DEFAULT nextval('tbl_category_id_seq') NOT NULL,
    "title" character varying(255) NOT NULL,
    "description" text NOT NULL,
    CONSTRAINT "tbl_category_id" PRIMARY KEY ("id")
) WITH (oids = false);




DROP TABLE IF EXISTS "tbl_user";
DROP SEQUENCE IF EXISTS tbl_user_id_seq;
CREATE SEQUENCE tbl_user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."tbl_user" (
    "id" integer DEFAULT nextval('tbl_user_id_seq') NOT NULL,
    "first_name" character varying(255) NOT NULL,
    "last_name" character varying(255) NOT NULL,
    "user_name" character varying(255) NOT NULL,
    "email" character varying(255) NOT NULL,
    "email_confirmed" boolean DEFAULT false NOT NULL,
    "avatar" character varying(255),
    "country_id" integer,
    "city_id" integer,
    "nationality_id" integer,
    "gender" character varying(255),
    "birth_date" date,
    "password" character varying(255) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "tbl_user_email" UNIQUE ("email"),
    CONSTRAINT "tbl_user_id" PRIMARY KEY ("id"),
    CONSTRAINT "tbl_user_user_name" UNIQUE ("user_name")
) WITH (oids = false);


DROP TABLE IF EXISTS "tbl_item";
DROP SEQUENCE IF EXISTS tbl_item_id_seq;
CREATE SEQUENCE tbl_item_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 4 CACHE 1;

CREATE TABLE "public"."tbl_item" (
    "id" integer DEFAULT nextval('tbl_item_id_seq') NOT NULL,
    "title" character varying(500) NOT NULL,
    "description" text NOT NULL,
    "price" numeric(12,2) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    "user_id" integer NOT NULL,
    "category_id" integer NOT NULL,
    CONSTRAINT "tbl_item_id" PRIMARY KEY ("id"),
    CONSTRAINT "tbl_item_category_id_fkey" FOREIGN KEY (category_id) REFERENCES tbl_category(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE,
    CONSTRAINT "tbl_item_user_id_fkey" FOREIGN KEY (user_id) REFERENCES tbl_user(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE
) WITH (oids = false);




DROP TABLE IF EXISTS "tbl_item_images";
DROP SEQUENCE IF EXISTS tbl_item_images_id_seq;
CREATE SEQUENCE tbl_item_images_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 4 CACHE 1;

CREATE TABLE "public"."tbl_item_images" (
    "id" integer DEFAULT nextval('tbl_item_images_id_seq') NOT NULL,
    "item_id" integer NOT NULL,
    "hash" character varying(255) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "tbl_item_images_id" PRIMARY KEY ("id"),
    CONSTRAINT "tbl_item_images_item_id_fkey1" FOREIGN KEY (item_id) REFERENCES tbl_item(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE
) WITH (oids = false);

CREATE INDEX "tbl_item_images_item_id" ON "public"."tbl_item_images" USING btree ("item_id");


DROP TABLE IF EXISTS "tbl_item_orders";
DROP SEQUENCE IF EXISTS tbl_item_orders_id_seq;
CREATE SEQUENCE tbl_item_orders_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 4 CACHE 1;

CREATE TABLE "public"."tbl_item_orders" (
    "id" integer DEFAULT nextval('tbl_item_orders_id_seq') NOT NULL,
    "item_id" integer NOT NULL,
    "user_id" integer NOT NULL,
    "finished" boolean NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "tbl_item_orders_id" PRIMARY KEY ("id"),
    CONSTRAINT "tbl_item_orders_item_id_user_id_finished" UNIQUE ("item_id", "user_id", "finished"),
    CONSTRAINT "tbl_item_orders_item_id_fkey" FOREIGN KEY (item_id) REFERENCES tbl_item(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE,
    CONSTRAINT "tbl_item_orders_user_id_fkey" FOREIGN KEY (user_id) REFERENCES tbl_user(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE
) WITH (oids = false);

CREATE INDEX "tbl_item_orders_item_id" ON "public"."tbl_item_orders" USING btree ("item_id");

CREATE INDEX "tbl_item_orders_user_id" ON "public"."tbl_item_orders" USING btree ("user_id");


DROP TABLE IF EXISTS "tbl_item_rating";
CREATE TABLE "public"."tbl_item_rating" (
    "item_id" integer NOT NULL,
    "rating" integer NOT NULL,
    "user_id" integer NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    CONSTRAINT "tbl_item_rating_item_id_user_id" PRIMARY KEY ("item_id", "user_id"),
    CONSTRAINT "tbl_item_rating_item_id_fkey" FOREIGN KEY (item_id) REFERENCES tbl_item(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE,
    CONSTRAINT "tbl_item_rating_user_id_fkey" FOREIGN KEY (user_id) REFERENCES tbl_user(id) ON UPDATE RESTRICT ON DELETE RESTRICT NOT DEFERRABLE
) WITH (oids = false);



CREATE INDEX "tbl_user_city_id" ON "public"."tbl_user" USING btree ("city_id");

CREATE INDEX "tbl_user_country_id" ON "public"."tbl_user" USING btree ("country_id");

CREATE INDEX "tbl_user_nationality_id" ON "public"."tbl_user" USING btree ("nationality_id");


CREATE INDEX "tbl_item_category_id" ON "public"."tbl_item" USING btree ("category_id");

CREATE INDEX "tbl_item_user_id" ON "public"."tbl_item" USING btree ("user_id");