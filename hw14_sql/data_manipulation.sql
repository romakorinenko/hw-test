-- Запросы на вставку, редактирование и удаление пользователей
INSERT INTO shop.users (name, email, password)
VALUES ('petya', 'petya@gmail.com', 'pass'),
       ('vasya', 'vasya@gmail.com', 'password'),
       ('kolya', 'kolya@gmail.com', 'passs'),
       ('grisha', 'grisha@gmail.com', 'PassworD'),
       ('kostya', 'kostya@gmail.com', 'qwerty123');

UPDATE shop.users
SET name = 'nikolay'
WHERE id = 3;

DELETE
FROM shop.users
WHERE id = 4;

-- Запросы на вставку, редактирование и удаление продуктов
INSERT INTO shop.products (name, price)
VALUES ('lemon', 10.00),
       ('banana', 20.00),
       ('hummer', 30.00),
       ('coke', 40.00),
       ('cherry', 60.00);

UPDATE shop.products
SET price = 50.00
WHERE id = 5;

DELETE
FROM shop.products
WHERE name = 'banana';

-- Запросы на сохранение и удаление заказов
INSERT INTO shop.orders (user_id, order_date, total_amount)
VALUES (1, now(), 30);
INSERT INTO shop.orders (user_id, order_date, total_amount)
VALUES (2, now(), 60);
INSERT INTO shop.orders (user_id, order_date, total_amount)
VALUES (5, now(), 20);

INSERT INTO shop.orderProducts (order_id, product_id)
VALUES (2, 1);
INSERT INTO shop.orderProducts (order_id, product_id)
VALUES (1, 1);
INSERT INTO shop.orderProducts (order_id, product_id)
VALUES (1, 5);
INSERT INTO shop.orderProducts (order_id, product_id)
VALUES (3, 2);

DELETE
FROM shop.orders
WHERE id = 3;

-- Запрос на выборку пользователей
SELECT *
FROM shop.users
WHERE name LIKE '%ya';

-- Запрос на выборку товаров
SELECT *
FROM shop.products
WHERE price > 20;

-- Запрос на выборку заказов по пользователю
SELECT *
FROM shop.orders
WHERE user_id = 1

SELECT *
FROM shop.orders o
         LEFT JOIN shop.users u
                   ON o.user_id = u.id
WHERE u.name = 'vasya';

-- Запрос на выборку статистики по пользователю (общая сумма заказов/средняя цена товара
SELECT u.name, SUM(o.total_amount) AS total_amount, AVG(p.price) AS avg_price
from shop.orders o
         LEFT JOIN shop.orderproducts op ON o.id = op.order_id
         LEFT JOIN shop.products p ON op.product_id = p.id
         LEFT JOIN shop.users u ON o.user_id = u.id
group by u.name
having u.name = 'petya'
