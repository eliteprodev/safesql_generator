CREATE TABLE products (
    id          BIGINT  NOT NULL AUTO_INCREMENT,
    name        TEXT    NOT NULL,
    price       INT     NOT NULL,
    PRIMARY KEY (id)
);

-- name: GetBetweenPrices :many
SELECT  *
FROM    products
WHERE   price BETWEEN ? AND ?;
