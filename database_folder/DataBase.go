package database_folder

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"webTest/struct_folder"
	"webTest/utilit"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
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

func (db *DB) GetMiniTypes() ([]string, error) {
	query := `
	SELECT slug FROM subcategories`
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

type Nav struct {
	Slug    string    `json:"slug"`
	Name    string    `json:"name"`
	MiniNav []MiniNav `json:"mininav"`
}

func (db *DB) GetAllCategories() ([]Nav, error) {
	query := `SELECT slug, name FROM categories`

	rows, err := db.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nav []Nav

	for rows.Next() {
		var n Nav
		if err := rows.Scan(&n.Slug, &n.Name); err != nil {
			return nil, err
		}

		miniN, err := db.GetMiniNavbar(n.Slug)
		if err != nil {
			return nil, err
		}

		n.MiniNav = miniN
		nav = append(nav, n)
	}
	return nav, nil
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

func (db *DB) GetProducts(
	productType string,
	productCategory string,
	search string,
	offset int,
	limit int,
) ([]struct_folder.MiniProducts, error) {
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
    (c.slug = ? OR ? = "")
    AND (s.slug = ? OR ? = "")
    AND (? = "" OR p.name LIKE CONCAT('%', ?, '%'))

    GROUP BY p.id
    ORDER BY p.id
    LIMIT ? OFFSET ?
    `

	rows, err := db.Db.Query(
		query,
		productType, productType,
		productCategory, productCategory,
		search, search,
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
		Value: []string{},
		Price: []int{},
	}

	for rows.Next() {
		var id int
		var val string
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

func (db *DB) getProductCategories(productSlug string, product *struct_folder.BonusInfoProduct) error {
	query := `SELECT product_subcategories.subcategory_slug 
FROM products 
JOIN product_subcategories ON products.id = product_subcategories.product_id
where products.slug = ?
`
	rows, err := db.Db.Query(query, productSlug)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cat string
		err := rows.Scan(&cat)
		if err != nil {
			return err
		}
		product.Categories = append(product.Categories, cat)
	}
	return nil
}

func (db *DB) GetBonusInfo(productSlug string) (struct_folder.BonusInfoProduct, error) {

	pi := struct_folder.BonusInfoProduct{
		Characteristic: []map[string]string{},
		Photo:          []string{},
		Variants:       struct_folder.Variant{},
		Categories:     []string{},
	}

	err := db.getPhoto(productSlug, &pi)

	if err != nil {
		return pi, err
	}

	err = db.getProductCategories(productSlug, &pi)

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

	err = db.Db.QueryRow(
		`SELECT name FROM products WHERE slug = ?`,
		productSlug,
	).Scan(&pi.Name)

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
		userPhone, err := utilit.Encrypt(info.Customer.Phone)
		if err != nil {
			return fmt.Errorf("Ошибка шифрования: %v", err)
		}
		res, errExec := tx.Exec(`INSERT INTO users(phone, name) VALUES (?, ?)`,
			userPhone, info.Customer.Name)
		if errExec != nil {
			err = errExec
			return err
		}
		lastID, _ := res.LastInsertId()
		userID = int(lastID)
	}

	// =========================
	// 2. Заказ
	// =========================
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
	// 3. Сохранение товаров в заказе (Архивный метод)
	// =========================
	for _, item := range info.Items {
		// Дефолтные заглушки на случай сбоя
		valStr := ""
		unitStr := "ml"

		// Перед вставкой подсматриваем текущие Unit и Value варианта из БД
		_ = tx.QueryRow(`
			SELECT value, unit FROM product_variants WHERE id = ?
		`, item.VariantID).Scan(&valStr, &unitStr)

		// Записываем неизменяемую текстовую копию в историю order_items
		_, err = tx.Exec(`
			INSERT INTO order_items(order_id, variant_id, product_name, variant_value, variant_unit, price_at_purchase, quantity)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, orderID, item.VariantID, item.Name, valStr, unitStr, int(item.Price), item.Qty)

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

// UpdateAdminPassword проверяет старый пароль и записывает новый для текущего пользователя
func (db *DB) UpdateAdminPassword(userName string, currentPassword string, newPassword string) (string, bool) {
	// 1. Получаем текущий хеш из базы
	querySelect := `SELECT password FROM administrator WHERE login = ?`
	var hashPassword string
	err := db.Db.QueryRow(querySelect, userName).Scan(&hashPassword)
	if err != nil {
		return "Пользователь не найден", false
	}

	// 2. Проверяем, совпадает ли введенный старый пароль
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(currentPassword))
	if err != nil {
		return "Неверный текущий пароль", false
	}

	// 3. Хешируем новый пароль
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "Ошибка генерации хеша", false
	}

	// 4. Записываем новый хеш в базу данных
	queryUpdate := `UPDATE administrator SET password = ? WHERE login = ?`
	_, err = db.Db.Exec(queryUpdate, string(newHash), userName)
	if err != nil {
		return "Ошибка обновления в БД", false
	}

	return "", true
}

func (db *DB) Login(userName string, userPassword string) bool {
	query := `SELECT password FROM administrator WHERE login = ?`
	var hashPassword string

	err := db.Db.QueryRow(query, userName).Scan(&hashPassword)
	if err != nil {
		return false // Пользователь не найден или ошибка БД
	}

	// Сравниваем введенный чистый пароль с хешем из базы данных
	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(userPassword))
	if err != nil {
		return false // Пароли не совпадают
	}

	return true
}

// Registration хеширует пароль и сохраняет его в базу данных
func (db *DB) Registration(userName string, userPassword string) bool {
	// Хешируем пароль с дефолтной стоимостью (DefaultCost = 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return false // Ошибка генерации хеша
	}

	query := `INSERT INTO administrator(login, password) VALUES(?, ?)`

	// Передаем логин и сгенерированный хеш (приводим []byte к string)
	_, err = db.Db.Exec(query, userName, string(hashedPassword))
	if err != nil {
		return false // Ошибка записи в БД (например, дубликат логина)
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
		var userPhone string
		err := rows.Scan(
			&item.OrdersInfo.Id, &item.OrdersInfo.Created_at,
			&item.OrdersInfo.CustomerName, &userPhone,
			&item.OrdersInfo.TotalPrice, &item.OrdersInfo.Status,
		)

		encUserPhone, err := utilit.Decrypt(userPhone)
		if err != nil {
			// Логируем ошибку, но НЕ останавливаем работу сервера,
			// отдавая в админку то, что есть в БД.
			log.Println("Ошибка дешифрования для заказа №", item.OrdersInfo.Id, ":", err)
			item.OrdersInfo.Phone = userPhone
		} else {
			item.OrdersInfo.Phone = encUserPhone
		}
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
	// Читаем сохраненные текстовые поля напрямую из истории без JOIN с каталогом товаров
	query := `
		SELECT 
			product_name, 
			variant_value, 
			variant_unit, 
			quantity, 
			price_at_purchase
		FROM order_items
		WHERE order_id = ?
	`
	for i := range *info {
		rows, err := db.Db.Query(query, (*info)[i].OrdersInfo.Id)
		if err != nil {
			return err
		}

		// ВАЖНО: Инициализируем слайс как пустой массив [], а не через var.
		// Благодаря этому json.Marshal вернет во фронтенд массив "ProductInfo": [], а не null.
		items := []struct_folder.ProductInfo{}

		for rows.Next() {
			var item struct_folder.ProductInfo

			// Сканируем исторические архивные данные в структуру ProductInfo
			err := rows.Scan(&item.Name, &item.Value, &item.Unit, &item.Count, &item.Price)
			if err != nil {
				rows.Close()
				return err
			}
			items = append(items, item)
		}
		rows.Close()

		// Сохраняем массив обратно в заказ
		(*info)[i].ProductInfo = items
	}
	return nil
}

func (db *DB) UpdateStatus(orderID int, orderStatus string) error {
	query := `UPDATE orders 
set status = ?
WHERE id = ?`
	_, err := db.Db.Exec(query, orderStatus, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateProduct(data struct_folder.UpdateProductData) error {
	if data.ID <= 0 {
		return fmt.Errorf("invalid product ID: %d", data.ID)
	}

	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}

	// helper для rollback
	fail := func(err error) error {
		tx.Rollback()
		return err
	}

	// ======================
	// 1. ОСНОВНЫЕ ДАННЫЕ
	// ======================
	_, err = tx.Exec(`
		UPDATE products 
		SET name = ?, description = ? 
		WHERE id = ?`,
		data.Name, data.Description, data.ID,
	)
	if err != nil {
		return fail(err)
	}

	// ======================
	// 2. ВАРИАНТЫ
	// ======================
	_, err = tx.Exec(`DELETE FROM product_variants WHERE product_id = ?`, data.ID)
	if err != nil {
		return fail(err)
	}

	// защита от кривых данных
	if len(data.Variants.Value) != len(data.Variants.Price) {
		return fail(fmt.Errorf("variants length mismatch"))
	}

	for i := range data.Variants.Value {
		_, err := tx.Exec(`
			INSERT INTO product_variants (product_id, value, unit, price)
			VALUES (?, ?, ?, ?)`,
			data.ID,
			data.Variants.Value[i],
			data.Variants.Unit,
			data.Variants.Price[i],
		)
		if err != nil {
			return fail(err)
		}
	}

	// ======================
	// 3. ХАРАКТЕРИСТИКИ
	// ======================
	_, err = tx.Exec(`DELETE FROM product_attributes WHERE product_id = ?`, data.ID)
	if err != nil {
		return fail(err)
	}

	for _, char := range data.Characteristics {
		var charID int64

		err := tx.QueryRow(`
			SELECT id FROM characteristics 
			WHERE key = ? AND value = ?`,
			char.Key, char.Value,
		).Scan(&charID)

		if err != nil {
			res, err := tx.Exec(`
				INSERT INTO characteristics (key, value) 
				VALUES (?, ?)`,
				char.Key, char.Value,
			)
			if err != nil {
				return fail(err)
			}
			charID, _ = res.LastInsertId()
		}

		_, err = tx.Exec(`
			INSERT INTO product_attributes (product_id, characteristic_id)
			VALUES (?, ?)`,
			data.ID, charID,
		)
		if err != nil {
			return fail(err)
		}
	}

	// ======================
	// 4. ПОДКАТЕГОРИИ
	// ======================
	_, err = tx.Exec(`DELETE FROM product_subcategories WHERE product_id = ?`, data.ID)
	if err != nil {
		return fail(err)
	}

	for _, subSlug := range data.Subcategories {
		if subSlug == "" {
			continue
		}

		_, err = tx.Exec(`
			INSERT INTO product_subcategories (product_id, subcategory_slug)
			VALUES (?, ?)`,
			data.ID, subSlug,
		)
		if err != nil {
			return fail(err)
		}
	}

	// ======================
	// 5. ФОТО (Сбор путей удаленных изображений)
	// ======================
	var toDelete []string

	rows, err := tx.Query(`
		SELECT url FROM product_photos 
		WHERE product_id = ?`, data.ID)

	if err != nil {
		return fail(err)
	}

	for rows.Next() {
		var dbUrl string
		rows.Scan(&dbUrl)

		stillExists := false
		for _, current := range data.ExistingPhotos {
			if dbUrl == current {
				stillExists = true
				break
			}
		}

		if !stillExists {
			toDelete = append(toDelete, dbUrl)
		}
	}
	rows.Close()

	// Пересоздаём связи фотографий в БД
	_, err = tx.Exec(`DELETE FROM product_photos WHERE product_id = ?`, data.ID)
	if err != nil {
		return fail(err)
	}

	for _, url := range data.ExistingPhotos {
		_, err = tx.Exec(`
			INSERT INTO product_photos (product_id, url)
			VALUES (?, ?)`,
			data.ID, url,
		)
		if err != nil {
			return fail(err)
		}
	}

	for _, newUrl := range data.NewPhotoPaths {
		_, err = tx.Exec(`
			INSERT INTO product_photos (product_id, url)
			VALUES (?, ?)`,
			data.ID, newUrl,
		)
		if err != nil {
			return fail(err)
		}
	}

	// ======================
	// COMMIT
	// ======================
	if err := tx.Commit(); err != nil {
		return err
	}

	// ======================
	// УДАЛЕНИЕ СТАРЫХ ФАЙЛОВ С ДИСКА (ИСПРАВЛЕНО)
	// ======================
	for _, path := range toDelete {
		// Защита от сбоя: убираем веб-префикс "/", если он сохранялся во фронтенде,
		// чтобы превратить путь в валидный локальный (например, из "/statics/..." в "statics/...")
		localPath := strings.TrimPrefix(path, "/")

		// Проверяем физическое существование файла перед удалением
		if _, err := os.Stat(localPath); err == nil {
			if err := os.Remove(localPath); err != nil {
				log.Println("ОШИБКА УДАЛЕНИЯ ФАЙЛА С ДИСКА:", err, "ПО ПУТИ:", localPath)
			}
		} else if os.IsNotExist(err) {
			log.Println("ФАЙЛ ДЛЯ УДАЛЕНИЯ НЕ НАЙДЕН НА ДИСКЕ:", localPath)
		}
	}

	return nil
}

func (db *DB) CreateProduct(data struct_folder.UpdateProductData) (int64, error) {
	tx, err := db.Db.Begin()
	if err != nil {
		return 0, err
	}

	fail := func(err error) (int64, error) {
		tx.Rollback()
		return 0, err
	}

	// ======================
	// 1. СГЕНЕРИРОВАТЬ УНИКАЛЬНЫЙ SLUG (name-1, name-2...)
	// ======================
	baseSlug := generateSlug(data.Name)
	slug, err := db.getUniqueSlug(tx, baseSlug)
	if err != nil {
		return fail(err)
	}

	// ======================
	// 2. PRODUCT (Вставка с гарантированно уникальным слагом)
	// ======================
	res, err := tx.Exec(`
		INSERT INTO products (name, description, slug)
		VALUES (?, ?, ?)`,
		data.Name, data.Description, slug,
	)
	if err != nil {
		return fail(err)
	}

	productID, _ := res.LastInsertId()

	// ======================
	// 3. ВАРИАНТЫ
	// ======================
	for i := range data.Variants.Value {
		_, err := tx.Exec(`
			INSERT INTO product_variants (product_id, value, unit, price)
			VALUES (?, ?, ?, ?)`,
			productID,
			data.Variants.Value[i],
			data.Variants.Unit,
			data.Variants.Price[i],
		)
		if err != nil {
			return fail(err)
		}
	}

	// ======================
	// 4. ХАРАКТЕРИСТИКИ
	// ======================
	for _, char := range data.Characteristics {
		var charID int64

		err := tx.QueryRow(`
			SELECT id FROM characteristics 
			WHERE key = ? AND value = ?`,
			char.Key, char.Value,
		).Scan(&charID)

		if err != nil {
			res, err := tx.Exec(`
				INSERT INTO characteristics (key, value)
				VALUES (?, ?)`,
				char.Key, char.Value,
			)
			if err != nil {
				return fail(err)
			}
			charID, _ = res.LastInsertId()
		}

		_, err = tx.Exec(`
			INSERT INTO product_attributes (product_id, characteristic_id)
			VALUES (?, ?)`,
			productID, charID,
		)
		if err != nil {
			return fail(err)
		}
	}

	// ======================
	// 5. ПОДКАТЕГОРИИ
	// ======================
	for _, sub := range data.Subcategories {
		if sub == "" {
			continue
		}
		_, err := tx.Exec(`
			INSERT INTO product_subcategories (product_id, subcategory_slug)
			VALUES (?, ?)`,
			productID, sub,
		)
		if err != nil {
			log.Println("INSERT SUBCATEGORY ERROR:", err)
			return fail(err)
		}
	}

	// ======================
	// 6. ФОТО
	// ======================
	for i, url := range data.NewPhotoPaths {
		_, err := tx.Exec(`
			INSERT INTO product_photos (product_id, url, position)
			VALUES (?, ?, ?)`,
			productID, url, i,
		)
		if err != nil {
			return fail(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return productID, nil
}

// getUniqueSlug проверяет существование слага в БД и инкрементирует суффикс до тех пор, пока не найдет свободный.
func (db *DB) getUniqueSlug(tx *sql.Tx, baseSlug string) (string, error) {
	currentSlug := baseSlug
	counter := 1

	for {
		var exists int
		// Выполняем проверку строго внутри текущей транзакции
		err := tx.QueryRow(`SELECT COUNT(1) FROM products WHERE slug = ?`, currentSlug).Scan(&exists)
		if err != nil {
			return "", err
		}

		// Если такого слага еще нет в базе — он подходит
		if exists == 0 {
			return currentSlug, nil
		}

		// Если слаг занят, добавляем/меняем числовой суффикс на следующий по порядку
		currentSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
}

func generateSlug(name string) string {
	translit := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'е': "e", 'ё': "e", 'ж': "zh", 'з': "z", 'и': "i",
		'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
		'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch",
		'ш': "sh", 'щ': "sch", 'ы': "y", 'э': "e", 'ю': "yu",
		'я': "ya",
	}

	name = strings.ToLower(name)

	var result strings.Builder

	for _, ch := range name {
		if val, ok := translit[ch]; ok {
			result.WriteString(val)
		} else if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			result.WriteRune(ch)
		} else {
			result.WriteRune('-')
		}
	}

	slug := result.String()

	// убрать двойные -
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// убрать - в начале/конце
	slug = strings.Trim(slug, "-")

	return slug
}

func (db *DB) AddSubcategory(subCat string, catSlug string) error {
	slugSubCat := generateSlug(subCat)

	query := `
	INSERT INTO subcategories(slug, parent_slug, name) VALUES(?, ?, ?)
	`

	_, err := db.Db.Exec(query, slugSubCat, catSlug, subCat)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddCategory(cat string) error {
	slugCat := generateSlug(cat)
	query := `insert into categories(slug, name) VALUES(?, ?)`

	_, err := db.Db.Exec(query, slugCat, cat)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteCategory(catSlug string, table string) error {
	query := fmt.Sprintf(`delete from %s WHERE slug = ?`, table)

	_, err := db.Db.Exec(query, catSlug)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteProduct(productID int) error {

	// ======================
	// GET PHOTOS
	// ======================
	rows, err := db.Db.Query(`
        SELECT url
        FROM product_photos
        WHERE product_id = ?
    `, productID)

	if err != nil {
		return err
	}

	defer rows.Close()

	var photos []string

	for rows.Next() {

		var path string

		if err := rows.Scan(&path); err != nil {
			return err
		}

		photos = append(photos, path)
	}

	// ======================
	// DELETE PRODUCT
	// ======================
	_, err = db.Db.Exec(`
        DELETE FROM products
        WHERE id = ?
    `, productID)

	if err != nil {
		return err
	}

	// ======================
	// DELETE FILES
	// ======================
	for _, path := range photos {

		err := os.Remove(path)

		if err != nil {

			log.Println(
				"DELETE FILE ERROR:",
				err,
			)
		}
	}

	return nil
}

func (db *DB) UpdateCategory(newCatName string, slug string, table string) error {
	query := fmt.Sprintf(`UPDATE %s SET name = ? WHERE slug = ?`, table)
	_, err := db.Db.Exec(query, newCatName, slug)

	if err != nil {
		return err
	}
	return nil
}
