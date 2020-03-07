CREATE TEMP TABLE orders (
  id integer
);
CREATE TEMP TABLE line_items (
  id integer,
  order_id integer,
  discount_applied integer
);
INSERT INTO orders (id) SELECT generate_series id FROM generate_series(0, 5, 1);

INSERT INTO line_items (id, order_id, discount_applied) VALUES
(0, 0, NULL),

(1, 1, 1),
(2, 1, NULL),

(3, 2, 1);

SELECT orders.id
FROM orders
LEFT OUTER JOIN line_items ON line_items.order_id = orders.id AND line_items.discount_applied IS NOT NULL
WHERE nulls.id IS NULL

 order_id |    array_agg
----------+-----------------
        0 | {NULL,NULL}
        2 | {1,1}
        1 | {1,NULL,1,NULL}


