PRAGMA foreign_keys = ON;

-- 1. Категории
CREATE TABLE IF NOT EXISTS categories (
    slug TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- 2. Подкатегории
CREATE TABLE IF NOT EXISTS subcategories (
    slug TEXT PRIMARY KEY,
    parent_slug TEXT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (parent_slug) REFERENCES categories(slug) ON DELETE CASCADE
);

-- 3. Характеристики
CREATE TABLE IF NOT EXISTS characteristics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL,
    value TEXT NOT NULL
);

-- 4. Продукты
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    slug TEXT UNIQUE NOT NULL
    -- category_slug УБРАН (теперь через subcategories)
);

-- 5. Связь товаров и подкатегорий (НОВОЕ)
CREATE TABLE IF NOT EXISTS product_subcategories (
    product_id INTEGER NOT NULL,
    subcategory_slug TEXT NOT NULL,
    PRIMARY KEY (product_id, subcategory_slug),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (subcategory_slug) REFERENCES subcategories(slug) ON DELETE CASCADE
);

-- 6. Варианты товара
CREATE TABLE IF NOT EXISTS product_variants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    unit TEXT DEFAULT 'ml',
    price INTEGER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- 7. Фото
CREATE TABLE IF NOT EXISTS product_photos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    position INTEGER DEFAULT 0,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- 8. Атрибуты товаров
CREATE TABLE IF NOT EXISTS product_attributes (
    product_id INTEGER NOT NULL,
    characteristic_id INTEGER NOT NULL,
    PRIMARY KEY (product_id, characteristic_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (characteristic_id) REFERENCES characteristics(id) ON DELETE CASCADE
);

-- 10. Заказы (Исправлено: убран variant_id)
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    comment TEXT,
    total_price INTEGER NOT NULL, -- Общая сумма заказа
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 11. Состав заказа (НОВАЯ ТАБЛИЦА)
CREATE TABLE order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    variant_id INTEGER NOT NULL,        -- Оставляем ID для истории (БЕЗ FOREIGN KEY)
    product_name TEXT NOT NULL,         -- Сюда жестко запишется имя товара
    variant_value TEXT NOT NULL,        -- Сюда запишется объем/размер (например, "50")
    variant_unit TEXT DEFAULT 'ml',     -- Единица измерения ("ml", "кг")
    price_at_purchase INTEGER NOT NULL, -- Точная цена на момент покупки
    quantity INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);
-- 12. Пользователь
CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	phone TEXT UNIQUE,
	name TEXT
	);

CREATE TABLE IF NOT EXISTS administrator (
    login TEXT,
    password Text
);