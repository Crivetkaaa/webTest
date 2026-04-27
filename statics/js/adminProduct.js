let offset = 0
let category = ""
let subcategory = ""


let limit = 20
async function getCategories() {
    const response = await fetch("/api/categories")
    const data = await response.json()
    return data
}
async function getProducts() {
    const bigCat = document.getElementById("category")
    const miniCat = document.getElementById("subcategory")


    const response = await fetch(`/api/get_products?type=${bigCat.value}&category=${miniCat.value}&offset=${offset}&limit=${limit}`)
    const data = await response.json()

    if (!data || data.length < limit) {
        const loadMore = document.getElementById("getProduct")
        loadMore.style.display = "none"
    }
    offset += data.length
    return data
}

function addOption(cat) {
    const bigCat = document.getElementById("category")
    const miniCat = document.getElementById("subcategory")

    cat.forEach(element => {
        const option = document.createElement("option")

        option.textContent = element.name
        option.value = element.slug
        bigCat.appendChild(option)
        element.mininav.forEach(el => {
            const miniOption = document.createElement("option")

            miniOption.textContent = el.name
            miniOption.value = el.slug
            miniCat.appendChild(miniOption)
        })

    });
}

function drawProduct(products, append) {
    const productCards = document.getElementById("productCards")
    if (!append) {
        productCards.innerHTML = ""
    }
    products.forEach(p => {
        const productCard = document.createElement("div")
        productCard.dataset.id = p.ID
        productCard.dataset.slug = p.Url
        const image = document.createElement("img")
        image.src = `/${p.MainPhoto}`
        const name = document.createElement("div")
        name.textContent = p.Name
        name.id = "productName"
        const price = document.createElement("div")
        price.textContent = p.Price + " P"

        productCard.appendChild(image)
        productCard.appendChild(name)
        productCard.appendChild(price)
        productCard.className = "product"

        productCards.appendChild(productCard)
    })

}

async function init() {
    const categories = await getCategories()
    const products = await getProducts()

    addOption(categories)
    drawProduct(products)
}

async function refreshProducts(append = false) {
    if (!append) {
        offset = 0
    }
    const products = await getProducts()
    drawProduct(products, append)

}

init()

function title(text) {
    const el = document.createElement("h3")
    el.textContent = text
    el.style.fontWeight = "bold"
    el.style.marginBottom = "8px"
    return el
}

function drawModal(data, product, allSubcategories) {
    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")
    content.innerHTML = "" // Очистка перед отрисовкой

    // ======================
    // 🏷 НАЗВАНИЕ
    // ======================
    const nameValue = product.querySelector(".name")?.textContent || product.children?.[1]?.textContent || ""
    const nameInput = document.createElement("input")
    nameInput.value = nameValue
    nameInput.className = "input"
    
    const nameBlock = document.createElement("div")
    nameBlock.className = "block"
    nameBlock.appendChild(title("Название"))
    nameBlock.appendChild(nameInput)
    content.appendChild(nameBlock)

    // ======================
    // 📸 ФОТО
    // ======================
    const photosBlock = document.createElement("div")
    photosBlock.className = "block"
    const photosContainer = document.createElement("div")
    photosContainer.className = "photos"

    const createPhotoElem = (src, file = null) => {
        const wrapper = document.createElement("div")
        wrapper.className = "photo-wrapper"
        
        // Сохраняем файл прямо в объекте элемента, если он есть
        if (file) wrapper._file = file
        // Сохраняем старый путь, если это не новый файл
        else wrapper.dataset.path = src

        const img = document.createElement("img")
        img.src = file ? URL.createObjectURL(file) : `/${src}`
        img.style.width = "100px"

        const removeBtn = document.createElement("button")
        removeBtn.innerHTML = "×"
        removeBtn.className = "remove-photo-btn"
        removeBtn.onclick = () => wrapper.remove()

        wrapper.appendChild(img)
        wrapper.appendChild(removeBtn)
        return wrapper
    }

    data.Photo.forEach(src => photosContainer.appendChild(createPhotoElem(src)))

    const addPhotoInput = document.createElement("input")
    addPhotoInput.type = "file"
    addPhotoInput.multiple = true
    addPhotoInput.accept = "image/*"
    addPhotoInput.style.display = "none"

    const addPhotoBtn = document.createElement("button")
    addPhotoBtn.textContent = "+ Добавить фото"
    addPhotoBtn.onclick = () => addPhotoInput.click()

    addPhotoInput.onchange = (e) => {
        Array.from(e.target.files).forEach(file => {
            photosContainer.appendChild(createPhotoElem(null, file))
        })
        addPhotoInput.value = ""
    }

    photosBlock.appendChild(title("Фото"))
    photosBlock.appendChild(photosContainer)
    photosBlock.appendChild(addPhotoBtn)
    content.appendChild(photosBlock)

    // ======================
    // 📝 ОПИСАНИЕ
    // ======================
    const desc = document.createElement("textarea")
    desc.value = data.Decscription || ""
    desc.className = "input"
    content.appendChild(title("Описание"))
    content.appendChild(desc)

    // ======================
    // 📦 ВАРИАНТЫ
    // ======================
    const variantsBlock = document.createElement("div")
    variantsBlock.className = "block"
    const variantsContainer = document.createElement("div")

    const createVariantRow = (val = "", pr = "") => {
        const row = document.createElement("div")
        row.className = "variant-row"
        row.innerHTML = `
            <input class="v-val" value="${val}" placeholder="${data.Variants.Unit}">
            <input class="v-price" value="${pr}" placeholder="Цена">
            <button class="remove-btn">❌</button>
        `
        row.querySelector(".remove-btn").onclick = () => row.remove()
        return row
    }

    data.Variants.Value.forEach((v, i) => {
        variantsContainer.appendChild(createVariantRow(v, data.Variants.Price[i]))
    })

    const addVariantBtn = document.createElement("button")
    addVariantBtn.textContent = "+ вариант"
    addVariantBtn.onclick = () => variantsContainer.appendChild(createVariantRow())

    variantsBlock.appendChild(title("Варианты"))
    variantsBlock.appendChild(variantsContainer)
    variantsBlock.appendChild(addVariantBtn)
    content.appendChild(variantsBlock)

    // ======================
    // 🏷 ХАРАКТЕРИСТИКИ
    // ======================
    const attrBlock = document.createElement("div")
    attrBlock.className = "block"
    const attrContainer = document.createElement("div")

    const createAttrRow = (k = "", v = "") => {
        const row = document.createElement("div")
        row.className = "attr-row"
        row.innerHTML = `
            <input class="a-key" value="${k}" placeholder="Ключ">
            <input class="a-val" value="${v}" placeholder="Значение">
            <button class="remove-btn">❌</button>
        `
        row.querySelector(".remove-btn").onclick = () => row.remove()
        return row
    }

    data.Characteristic.forEach(c => attrContainer.appendChild(createAttrRow(c.key, c.value)))

    const addAttrBtn = document.createElement("button")
    addAttrBtn.textContent = "+ характеристика"
    addAttrBtn.onclick = () => attrContainer.appendChild(createAttrRow())

    attrBlock.appendChild(title("Характеристики"))
    attrBlock.appendChild(attrContainer)
    attrBlock.appendChild(addAttrBtn)
    content.appendChild(attrBlock)

    // 📂 КАТЕГОРИИ (НОВОЕ)
    // ======================
    const catBlock = document.createElement("div")
    catBlock.className = "block"
    catBlock.appendChild(title("Подкатегории"))

    const catContainer = document.createElement("div")
    catContainer.className = "categories-grid" // Сетка для чекбоксов
    
    allSubcategories.forEach(sub => {
        const label = document.createElement("label")
        label.className = "cat-label"
        
        const checkbox = document.createElement("input")
        checkbox.type = "checkbox"
        checkbox.value = sub.slug // Используем slug как ID
        checkbox.className = "cat-checkbox"
        
        // Если слаг подкатегории есть в массиве data.Categories, отмечаем его
        if (data.Categories && data.Categories.includes(sub.slug)) {
            checkbox.checked = true
        }

        label.appendChild(checkbox)
        label.append(` ${sub.name}`)
        catContainer.appendChild(label)
    })
    
    catBlock.appendChild(catContainer)
    content.appendChild(catBlock)

    // ======================
    // 💾 СОХРАНЕНИЕ (ОБНОВЛЕНО)
    // ======================
    const save = document.createElement("button")
    save.textContent = "Сохранить"
    save.className = "save-btn"

    save.onclick = async () => {
    const productId = product.dataset.id;
    const newName = nameInput.value.trim();

    // 1. Валидация
    if (!productId || productId === "0") {
        alert("Ошибка: ID товара не найден!");
        return;
    }
    if (!newName) {
        alert("Название товара не может быть пустым!");
        return;
    }

    const formData = new FormData();
    
    // 2. Основные данные
    formData.append("id", productId);
    formData.append("name", newName);
    formData.append("description", desc.value); // Должно быть определено выше в drawModal

    // 3. Сбор фото
    const existing = []
    photosContainer.querySelectorAll(".photo-wrapper").forEach(el => {
        if (el._file) formData.append("newPhotos", el._file);
        if (el.dataset.path) existing.push(el.dataset.path);
    });
    formData.append("existingPhotos", JSON.stringify(existing));

    // 4. Сбор Вариантов
    const vars = { Unit: data.Variants.Unit, Value: [], Price: [] };
    variantsContainer.querySelectorAll(".variant-row").forEach(row => {
        const v = row.querySelector(".v-val").value.trim();
        const p = row.querySelector(".v-price").value.trim();
        if (v && p) { 
            vars.Value.push(v); 
            vars.Price.push(p); 
        }
    });
    formData.append("variants", JSON.stringify(vars));

    // 5. Сбор Характеристик
    const attrs = [];
    attrContainer.querySelectorAll(".attr-row").forEach(row => {
        const k = row.querySelector(".a-key").value.trim();
        const v = row.querySelector(".a-val").value.trim();
        if (k) attrs.push({ key: k, value: v });
    });
    formData.append("characteristics", JSON.stringify(attrs));

    // 6. Сбор Категорий
    const selectedCats = [];
    // Используем тот контейнер, где отрисовывали чекбоксы
    catContainer.querySelectorAll(".cat-checkbox:checked").forEach(cb => {
        selectedCats.push(cb.value);
    });
    formData.append("subcategories", JSON.stringify(selectedCats));

    // 7. Отправка
    try {
        const res = await fetch("/admin/update_product", { method: "POST", body: formData });
        if (res.ok) { 
            alert("Сохранено!");
            modal.classList.remove("active"); 
            location.reload(); 
        } else {
            const errData = await res.json();
            alert("Ошибка сервера: " + (errData.error || "Неизвестная ошибка"));
        }
    } catch (e) {
        console.error(e);
        alert("Ошибка сети или сервера");
    }
}


    // Собираем всё в конце
    // (Убедитесь, что appendChild идут в нужном вам порядке)
    content.appendChild(save)

}


async function getProductsInfo(product) {
    // 1. Получаем инфо о товаре
    const product_info = await fetch(`/api/product_info/${product.dataset.slug}`)
    const data = await product_info.json()

    // 2. Получаем все категории (через твою функцию или fetch)
    const cat = await getCategories() 
    
    // 3. Собираем все подкатегории (mininav) в один плоский массив
    let allSubcategories = []
    cat.forEach(c => {
        if (c.mininav) {
            c.mininav.forEach(m => {
                allSubcategories.push({
                    slug: m.slug,
                    name: `${c.name} > ${m.name}` // Для наглядности добавим имя родителя
                })
            })
        }
    })

    const modal = document.getElementById("modal")
    // 4. Передаем плоский список в drawModal
    await drawModal(data, product, allSubcategories)
    modal.classList.add("active")
}


function addProduct() {
    refreshProducts(true)
}

document.getElementById("productCards").addEventListener("click", (e) => {
    const card = e.target.closest(".product")
    if (!card) return
    getProductsInfo(card)
    console.log("клик по товару:", card.dataset.id)
})

document.addEventListener("keydown", (e) => {
    const modal = document.getElementById("modal")
    if (e.key === "Escape") {
        modal.classList.remove("active")
    }
})

document.addEventListener("click", (e) => {
    const modal = document.getElementById("modal")
    if (e.target === modal) {
        modal.classList.remove("active")
    }
})