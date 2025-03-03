CREATE DATABASE IF NOT EXISTS bo1;
USE bo1;

CREATE TABLE ProductSales (
  id INT AUTO_INCREMENT PRIMARY KEY,
  sale_date DATE,
  region VARCHAR(100),
  product VARCHAR(100),
  qty INT,
  cost DECIMAL(10,2),
  amt DECIMAL(10,2),
  tax DECIMAL(10,2),
  total DECIMAL(10,2)
);