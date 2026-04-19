import sqlite3
import random
import string

DB_PATH = "database_folder/db.db"

PHOTO_PATH = "statics/img/product_img/12.png"

PRODUCT_TYPES = ["cloth", "parfume"]

CATEGORIES = ["new", "popular", "luxury", "flower", "bags", "man", "women"]

WEB_TYPES = ["cloth", "parfume"]


def rand_name():
    words = ["Crystal", "Golden", "Silk", "Urban", "Midnight", "Rose", "Amber", "Velvet", "Night", "Blue"]
    return f"{random.choice(words)}-{random.choice(words)}-{random.randint(1, 999)}"


def connect():
    conn = sqlite3.connect(DB_PATH)
    conn.execute("PRAGMA foreign_keys = ON")
    return conn


def clear_db(cur):
    tables = [
        "product_categories",
        "product_characteristics",
        "product_variants",
        "product_photos",
        "products",
        "categories",
        "mini_web_categories",
        "web_categories",
        "web_type"
    ]

    for t in tables:
        cur.execute(f"DELETE FROM {t}")


def seed_base(cur):
    # web_type
    for wt in WEB_TYPES:
        cur.execute("INSERT INTO web_type(categories) VALUES (?)", (wt,))

    # web_categories
    cur.execute("INSERT INTO web_categories(slug, cat_name) VALUES (?, ?)", ("parfume", "Парфюм"))
    cur.execute("INSERT INTO web_categories(slug, cat_name) VALUES (?, ?)", ("cloth", "Одежда"))

    # mini categories
    cur.execute("INSERT INTO mini_web_categories(web_categories_id, slug, label) VALUES (1, 'flower', 'Цветочные')")
    cur.execute("INSERT INTO mini_web_categories(web_categories_id, slug, label) VALUES (1, 'luxury', 'Люкс')")
    cur.execute("INSERT INTO mini_web_categories(web_categories_id, slug, label) VALUES (2, 'bags', 'Сумки')")
    cur.execute("INSERT INTO mini_web_categories(web_categories_id, slug, label) VALUES (2, 'new', 'Новинки')")

    # categories
    for c in CATEGORIES:
        cur.execute("INSERT INTO categories(name) VALUES (?)", (c,))


def seed_products(cur):
    for i in range(1, 101):

        ptype = random.choice(PRODUCT_TYPES)
        name = rand_name()
        url = f"product/{name}"

        cur.execute("""
            INSERT INTO products(name, url, description, product_type)
            VALUES (?, ?, ?, ?)
        """, (
            name,
            url,
            f"Описание товара {name}",
            ptype
        ))

        product_id = cur.lastrowid

        # photo
        cur.execute("""
            INSERT INTO product_photos(product_id, photo_url, position)
            VALUES (?, ?, 0)
        """, (product_id, PHOTO_PATH))

        # variants (2-3 варианта)
        for v in [30, 50, 100]:
            cur.execute("""
                INSERT INTO product_variants(product_id, value, price, unit)
                VALUES (?, ?, ?, ?)
            """, (
                product_id,
                v,
                random.randint(3000, 15000),
                "ml" if ptype == "parfume" else "size"
            ))

        # characteristics
        cur.execute("""
            INSERT INTO product_characteristics(product_id, key, value)
            VALUES (?, 'Производитель', 'AutoGen Brand')
        """, (product_id,))

        cur.execute("""
            INSERT INTO product_characteristics(product_id, key, value)
            VALUES (?, 'Страна', 'Italy')
        """, (product_id,))

        # categories (рандом 1-3)
        for cat in random.sample(CATEGORIES, k=random.randint(1, 3)):
            cur.execute("""
                INSERT INTO product_categories(product_id, category_id)
                VALUES (?, (SELECT id FROM categories WHERE name=?))
            """, (product_id, cat))


def main():
    conn = connect()
    cur = conn.cursor()

    clear_db(cur)
    seed_base(cur)
    seed_products(cur)

    conn.commit()
    conn.close()

    print("DB seeded successfully (100 products)")


if __name__ == "__main__":
    main()