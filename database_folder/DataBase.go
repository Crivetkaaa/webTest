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

func (db *DB) Login(userName string, userPassword string) bool {
	query := `SELECT password FROM administrator WHERE login = ?`
	var password string
	err := db.Db.QueryRow(query, userName).Scan(&password)
	if err != nil {
		return false
	}
	// TODO hash
	fmt.Println(password)
	if password != userPassword {
		return false
	}
	return true
}

// Обновляем сигнатуру функции: добавляем limit, offset и status
func (db *DB) getOrders(info *[]struct_folder.AdminInfo, limit, offset int, status string) error {
	// Базовый запрос
	query := `
    SELECT o.id, o.created_at, u.name, u.phone, o.total_price, o.status
    FROM orders o
    JOIN users u ON o.user_id = u.id `

	var args []interface{}

	// Если статус не "all", добавляем фильтрацию в SQL
	if status != "all" && status != "" {
		query += " WHERE o.status = ? "
		args = append(args, status)
	}

	// Добавляем сортировку и пагинацию
	query += " ORDER BY o.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item struct_folder.AdminInfo
		err := rows.Scan(
			&item.OrdersInfo.Id, &item.OrdersInfo.Created_at,
			&item.OrdersInfo.CustomerName, &item.OrdersInfo.Phone,
			&item.OrdersInfo.TotalPrice, &item.OrdersInfo.Status,
		)
		if err != nil {
			return err
		}
		*info = append(*info, item)
	}
	return nil
}

// Обновляем главную функцию
func (db *DB) GetOrdersAdmin(limit, offset int, status string) ([]struct_folder.AdminInfo, error) {
	var adminInfo []struct_folder.AdminInfo

	err := db.getOrders(&adminInfo, limit, offset, status)
	if err != nil {
		return nil, err
	}

	err = db.getOrdersInfo(&adminInfo)
	if err != nil {
		return nil, err
	}

	return adminInfo, nil
}

func (db *DB) getOrdersInfo(info *[]struct_folder.AdminInfo) error {
	query :=
		`
	SELECT 
    p.name, 
    pv.value, 
    pv.unit, 
    oi.quantity, 
    oi.price_at_purchase
FROM order_items oi
JOIN product_variants pv ON oi.variant_id = pv.id
JOIN products p ON pv.product_id = p.id
WHERE oi.order_id = ?
	`
	for i := range *info {
		rows, err := db.Db.Query(query, (*info)[i].OrdersInfo.Id)
		if err != nil {
			return err
		}

		var items []struct_folder.ProductInfo
		for rows.Next() {
			var item struct_folder.ProductInfo

			err := rows.Scan(&item.Name, &item.Value, &item.Unit, &item.Count, &item.Price)
			if err != nil {
				return err
			}
			items = append(items, item)
		}
		rows.Close()
		(*info)[i].ProductInfo = items
	}
	return nil
}

func (db *DB) UpdateStatus(orderID int, orderStatus string) error {
	fmt.Println(orderStatus)
	query := `UPDATE orders 
set status = ?
WHERE id = ?`
	_, err := db.Db.Exec(query, orderStatus, orderID)
	if err != nil {
		return err
	}
	return nil
}
