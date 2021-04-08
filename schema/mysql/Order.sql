create table Order 
(
    id          SMALLINT NOT NULL AUTO_INCREMENT,
    sku         VARCHAR(30), -- reference to OrderHeader.Sku
    batch_no    VARCHAR(30), -- reference to OrderHeader.BatchNo
    cust_no     VARCHAR(30),
    order_num   VARCHAR(30),
    order_qty   INT,
    price       FLOAT,
    createdAt   TIMESTAMP,
    created_by  VARCHAR(30),
    updated_at  TIMESTAMP,
    updated_by  VARCHAR(30),
    when        TIMESTAMP,
    PRIMARY KEY (id)
)
