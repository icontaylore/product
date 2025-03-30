CREATE TABLE product(
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL ,
    description VARCHAR(255),
    price DECIMAL NOT NULL,
    stock_quantity INT NOT NULL
);

INSERT INTO product(name,description,price,stock_quantity)
VALUES('test','test',40.40,1);