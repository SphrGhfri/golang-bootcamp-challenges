-- Section1
CREATE INDEX idx_users_name ON users (name);
-- Section2
CREATE INDEX idx_products_covering ON products (category_id, price, id, name);