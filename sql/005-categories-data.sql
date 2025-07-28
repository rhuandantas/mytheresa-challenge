-- 2. Seed categories
INSERT INTO categories (code, name) VALUES
('clothing', 'Clothing'),
('shoes', 'Shoes'),
('accessories', 'Accessories')
    ON CONFLICT (code) DO NOTHING;

-- 3. Add category_id to products
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

-- 4. Update products with correct category_id
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'clothing')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'shoes')
WHERE code IN ('PROD002', 'PROD006');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'accessories')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');

-- 5. Make category_id NOT NULL if all products are assigned
ALTER TABLE products
    ALTER COLUMN category_id SET NOT NULL;

-- 6. Add indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_code ON categories(code);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);