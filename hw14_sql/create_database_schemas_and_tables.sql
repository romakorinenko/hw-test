CREATE DATABASE test_db OWNER postgres;
-- DROP DATABASE IF EXISTS test_db;

CREATE schema IF NOT EXISTS shop AUTHORIZATION postgres;
-- DROP SCHEMA IF EXISTS shop CASCADE;

CREATE table IF NOT EXISTS shop.users
(
    id       BIGSERIAL PRIMARY KEY,
    name     VARCHAR(255) NOT NULL DEFAULT '',
    email    VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
-- DROP TABLE shop.users;

CREATE INDEX IF NOT EXISTS users_name_idx ON shop.users USING btree (name);
-- DROP INDEX shop.users_name_idx;

CREATE TABLE IF NOT EXISTS shop.orders
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT         NOT NULL REFERENCES shop.users (id),
    order_date   DATE           NOT NULL DEFAULT CURRENT_DATE,
    total_amount NUMERIC(10, 2) NOT NULL check (total_amount > 0)
);
--DROP TABLE shop.orders;

CREATE INDEX IF NOT EXISTS orders_user_id_idx ON shop.orders USING btree (user_id);
-- DROP INDEX shop.orders_user_id_idx;

CREATE INDEX IF NOT EXISTS orders_order_date_idx ON shop.orders USING btree (order_date);
-- DROP INDEX shop.orders_order_date_idx;

CREATE INDEX IF NOT EXISTS orders_total_amount_idx ON shop.orders USING btree (total_amount);
-- DROP INDEX shop.orders_total_amount_idx;

CREATE TABLE IF NOT EXISTS shop.products
(
    id    BIGSERIAL PRIMARY KEY,
    name  varchar(255)   NOT NULL UNIQUE,
    price NUMERIC(10, 2) NOT NULL CHECK (price > 0)
);
-- DROP TABLE shop.products;

CREATE INDEX IF NOT EXISTS products_name_idx ON shop.products USING btree (name);
-- DROP INDEX shop.products_name_idx;

CREATE INDEX IF NOT EXISTS products_price_idx ON shop.products USING btree (price);
-- DROP INDEX shop.products_price_idx;

CREATE TABLE IF NOT EXISTS shop.orderProducts
(
    id         SERIAL PRIMARY KEY,
    order_id   BIGINT NOT NULL REFERENCES shop.orders (id) ON DELETE CASCADE ON UPDATE CASCADE,
    product_id BIGINT NOT NULL REFERENCES shop.products (id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- DROP TABLE shop.orders_products;

CREATE INDEX IF NOT EXISTS order_products_order_id_idx ON shop.orderProducts USING btree (order_id);
-- DROP INDEX shop.order_products_order_id_idx;

CREATE INDEX IF NOT EXISTS order_products_product_id_idx ON shop.orderProducts USING btree (product_id);
-- DROP INDEX shop.order_products_product_id_idx;
