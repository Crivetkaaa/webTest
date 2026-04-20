package database_folder

import (
	"database/sql"
	"fmt"
	"webTest/struct_folder"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Db *sql.DB
}

func (db *DB) GetType() (string, error) {
	query := "SELECT slug FROM categories LIMIT 1"
	row, err := db.Db.Query(query)
	if err != nil {
		return "", err
	}
	defer row.Close()
	var cat string
	for row.Next() {
		err := row.Scan(&cat)
		if err != nil {
			return "", nil
		}
	}
	return cat, nil
}

func (db *DB) GetTypes() ([]string, error) {
	query := `
	SELECT slug FROM categories`
	rows, err := db.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []string

	for rows.Next() {
		var cat string
		err := rows.Scan(&cat)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

type MiniCategory struct {
	Name string
	Slug string
}

type Category struct {
	ID    int
	Name  string
	Slug  string
	Items []MiniCategory
}

type NavItem struct {
	Slug string
	Name string
}

func (db *DB) GetBigNavbar() ([]NavItem, error) {
	query := `select slug, name from categories`
	rows, err := db.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navbar []NavItem

	for rows.Next() {
		var item NavItem
		if err := rows.Scan(&item.Slug, &item.Name); err != nil {
			return nil, err
		}
		navbar = append(navbar, item)
	}

	return navbar, nil
}

type MiniNav struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (db *DB) GetMiniNavbar(slug string) ([]MiniNav, error) {
	query := `select slug, name from subcategories where parent_slug=?`
	rows, err := db.Db.Query(query, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MiniNav

	for rows.Next() {
		var item MiniNav
		if err := rows.Scan(&item.Slug, &item.Name); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (db *DB) getSubcat(products *[]struct_folder.MiniProducts) error {
	query := `SELECT subcategory_slug FROM product_subcategories WHERE product_id = ?`

	for i := range *products {
		// Работаем напрямую с элементом слайса по индексу
		p := &(*products)[i]

		rows, err := db.Db.Query(query, p.ID)
		if err != nil {
			return err
		}

		// Важно: defer внутри цикла — плохая практика (может забить память)
		// Закрываем rows вручную в конце итерации
		for rows.Next() {
			var cat string
			if err := rows.Scan(&cat); err == nil {
				p.Categories = append(p.Categories, cat)
			}
		}
		rows.Close()
	}
	return nil
}

func (db *DB) GetProducts(productType string, productCategory string, offset int, limit int) ([]struct_folder.MiniProducts, error) {
	query := `
    SELECT 
		p.id,
        p.name, 
        p.slug, 
        MIN(pp.url) as url, 
        MIN(pv.price) as price
    FROM products p

    -- связь с подкатегориями
    LEFT JOIN product_subcategories ps ON ps.product_id = p.id
    LEFT JOIN subcategories s ON s.slug = ps.subcategory_slug
    LEFT JOIN categories c ON c.slug = s.parent_slug

    -- варианты и фото
    JOIN product_variants pv ON pv.product_id = p.id
    JOIN product_photos pp ON pp.product_id = p.id

    WHERE 
        (c.slug = ? OR ? = "")          -- фильтр по категории
        AND (s.slug = ? OR ? = "")      -- фильтр по подкатегории

    GROUP BY p.id
    ORDER BY p.id
    LIMIT ? OFFSET ?
    `

	rows, err := db.Db.Query(
		query,
		productType, productType, // category
		productCategory, productCategory, // subcategory
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []struct_folder.MiniProducts
	for rows.Next() {
		var product struct_folder.MiniProducts
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Url,
			&product.MainPhoto,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	db.getSubcat(&products)
	return products, nil
}

func (db *DB) getPhoto(productSlug string, pi *struct_folder.BonusInfoProduct) error {
	query := `SELECT pp.url
        FROM product_photos pp
        JOIN products p ON pp.product_id = p.id
        WHERE p.slug = ?
        ORDER BY pp.position ASC`
	rows, err := db.Db.Query(query, productSlug)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			return err
		}
		pi.Photo = append(pi.Photo, url)
	}
	return nil
}

func (db *DB) getVariants(productSlug string, pi *struct_folder.BonusInfoProduct) error {
	query := `SELECT pv.id, pv.value, pv.price, pv.unit
        FROM product_variants pv
        JOIN products p ON pv.product_id = p.id
        WHERE p.slug = ?
		ORDER BY pv.price ASC`

	rows, err := db.Db.Query(query, productSlug)
	if err != nil {
		return err
	}
	defer rows.Close()

	pi.Variants = struct_folder.Variant{
		Id:    []int{},
		Value: []int{},
		Price: []int{},
	}

	for rows.Next() {
		var id int
		var val int
		var pr int
		var unit string

		err := rows.Scan(&id, &val, &pr, &unit)
		if err != nil {
			return err
		}

		pi.Variants.Unit = unit
		pi.Variants.Id = append(pi.Variants.Id, id)
		pi.Variants.Price = append(pi.Variants.Price, pr)
		pi.Variants.Value = append(pi.Variants.Value, val)
	}
	fmt.Println(pi.Variants.Id)
	return nil
}

func (db *DB) getCharacteristic(productSlug string, pi *struct_folder.BonusInfoProduct) error {
	query := `SELECT 
    c.key, 
    c.value
FROM characteristics c
JOIN product_attributes pa ON pa.characteristic_id = c.id
JOIN products p ON pa.product_id = p.id
WHERE p.slug = ?;
`
	rows, err := db.Db.Query(query, productSlug)
	if err != nil {
		return err
	}
	defer rows.Close()

	var char []map[string]string
	for rows.Next() {
		var lefr string
		var right string

		err := rows.Scan(&lefr, &right)
		if err != nil {
			return err
		}
		char = append(char, map[string]string{
			"key":   lefr,
			"value": right,
		})
	}
	pi.Characteristic = char

	return nil
}

func (db *DB) getDescription(productSlug string, pi *struct_folder.BonusInfoProduct) error {
	query := `select description from products where slug=?`

	err := db.Db.QueryRow(query, productSlug).Scan(&pi.Decscription)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getBaseProductInfo(productSlug string, pr *struct_folder.Product) error {
	query := `SELECT products.name FROM products WHERE products.slug = ?`
	err := db.Db.QueryRow(query, productSlug).Scan(&pr.Name)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetBonusInfo(productSlug string) (struct_folder.BonusInfoProduct, error) {
	pi := struct_folder.BonusInfoProduct{
		Characteristic: []map[string]string{},
		Photo:          []string{},
		Variants:       struct_folder.Variant{},
	}
	err := db.getPhoto(productSlug, &pi)
	if err != nil {
		return pi, err
	}

	err = db.getVariants(productSlug, &pi)
	if err != nil {
		return pi, err
	}
	err = db.getCharacteristic(productSlug, &pi)
	if err != nil {
		return pi, err
	}
	err = db.getDescription(productSlug, &pi)
	if err != nil {
		return pi, err
	}

	return pi, nil
}

func (db *DB) GetProduct(productSlug string) (struct_folder.Product, error) {
	products := struct_folder.Product{}
	miniProd := struct_folder.BonusInfoProduct{
		Characteristic: []map[string]string{},
		Photo:          []string{},
		Variants:       struct_folder.Variant{},
	}

	err := db.getPhoto(productSlug, &miniProd)
	if err != nil {
		return products, err
	}
	products.Photo = miniProd.Photo

	err = db.getVariants(productSlug, &miniProd)
	if err != nil {
		return products, err
	}
	products.Variants = miniProd.Variants
	err = db.getCharacteristic(productSlug, &miniProd)
	if err != nil {
		return products, err
	}
	products.Characteristic = miniProd.Characteristic
	err = db.getDescription(productSlug, &miniProd)
	if err != nil {
		return products, err
	}
	products.Description = miniProd.Decscription
	err = db.getBaseProductInfo(productSlug, &products)
	if err != nil {
		return products, err
	}
	products.Url = productSlug

	return products, nil

}

func (db *DB) InsertOrder(info struct_folder.OrderData) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}

	// Откат при любой ошибке
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// =========================
	// 1. Пользователь
	// =========================
	var userID int
	err = tx.QueryRow(`SELECT id FROM users WHERE phone = ?`, info.Customer.Phone).Scan(&userID)

	if err != nil {
		res, errExec := tx.Exec(`INSERT INTO users(phone, name) VALUES (?, ?)`,
			info.Customer.Phone, info.Customer.Name)
		if errExec != nil {
			err = errExec
			return err
		}
		lastID, _ := res.LastInsertId()
		userID = int(lastID)
	}

	// =========================
	// 2. Заказ (Сначала создаем запись в orders!)
	// =========================
	// ВАЖНО: Мы вставляем в таблицу orders, а не order_items
	res, err := tx.Exec(`
		INSERT INTO orders(user_id, comment, total_price) 
		VALUES (?, ?, ?)
	`, userID, info.Customer.Comment, info.Total)

	if err != nil {
		return err
	}

	orderID64, _ := res.LastInsertId()
	orderID := int(orderID64)

	// =========================
	// 3. Товары в заказе
	// =========================
	for _, item := range info.Items {
		// ОШИБКА БЫЛА ТУТ: Заменяем item.ID (0) на item.VariantID (3)
		_, err = tx.Exec(`
			INSERT INTO order_items(order_id, variant_id, quantity, price_at_purchase)
			VALUES (?, ?, ?, ?)
		`, orderID, item.VariantID, item.Qty, item.Price)

		if err != nil {
			return err
		}
	}

	// =========================
	// 4. Коммит
	// =========================
	err = tx.Commit()
	return err
}
