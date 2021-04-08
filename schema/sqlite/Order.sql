create table Order 
(
    id          INTEGER NOT NULL AUTOINCREMENT,
    sku         VARCHAR(30), -- reference to OrderHeader.Sku
    batch_no    VARCHAR(30), -- reference to OrderHeader.BatchNo
    cust_no     VARCHAR(30),
    order_num   VARCHAR(30),
    order_qty   INTEGER,
    price       BIGINT,
    createdAt   VARCHAR,
    created_by  VARCHAR(30),
    updated_at  VARCHAR,
    updated_by  VARCHAR(30),
    when        VARCHAR,
    PRIMARY KEY (id)
)
