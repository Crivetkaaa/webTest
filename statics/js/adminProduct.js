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

    const response = await fetch(
        `/api/get_products?type=${bigCat.value}&category=${miniCat.value}&offset=${offset}&limit=${limit}`
    )

    let data = await response.json()

    // защита от null
    if (!Array.isArray(data)) {
        data = []
    }

    const loadMore = document.getElementById("getProduct")

    if (data.length < limit) {
        loadMore.style.display = "none"
    } else {
        loadMore.style.display = "block"
    }

    offset += data.length

    return data
}

function addOption(cat) {
    const bigCat = document.getElementById("category")
    const miniCat = document.getElementById("subcategory")

    bigCat.innerHTML = `<option value="">Все категории</option>`
    miniCat.innerHTML = `<option value="">Все подкатегории</option>`

    cat.forEach(element => {

        // ======================
        // CATEGORY
        // ======================
        const option = document.createElement("option")

        option.textContent = element.name
        option.value = element.slug

        bigCat.appendChild(option)

        // ======================
        // SUBCATEGORIES
        // ======================
        const subcategories = element.mininav || []

        subcategories.forEach(el => {

            const miniOption = document.createElement("option")

            miniOption.textContent =
                `${element.name} > ${el.name}`

            miniOption.value = el.slug

            miniCat.appendChild(miniOption)
        })
    })
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
        image.className = "product-img"

        const name = document.createElement("div")
        name.textContent = p.Name
        name.className = "product-title"

        const price = document.createElement("div")
        price.textContent = p.Price + " P"
        price.className = "product-price"

        productCard.appendChild(image)
        productCard.appendChild(name)
        productCard.appendChild(price)

        productCard.className = "product-card"

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

function drawModal(product, allSubcategories, isCreate = false) {
    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")
    content.innerHTML = ""

    const closeBtn = document.createElement("button")

    closeBtn.innerHTML = "×"
    closeBtn.className = "modal-close-btn"

    closeBtn.onclick = () => {
        modal.classList.remove("active")
    }

    content.appendChild(closeBtn)
    // ======================
    // 🏷 НАЗВАНИЕ
    // ======================
    const nameInput = document.createElement("input")
    nameInput.value = product.name
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

        if (file) wrapper._file = file
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

    product.photos.forEach(src => {
        photosContainer.appendChild(createPhotoElem(src))
    })

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
    desc.value = product.description
    desc.className = "input"

    const descBlock = document.createElement("div")
    descBlock.className = "block"
    descBlock.appendChild(title("Описание"))
    descBlock.appendChild(desc)
    content.appendChild(descBlock)

    // ======================
    // 📦 ВАРИАНТЫ + UNIT
    // ======================
    const variantsBlock = document.createElement("div")
    variantsBlock.className = "block"

    const unitInput = document.createElement("input")
    unitInput.className = "input"
    unitInput.placeholder = "Единица (ml, size, kg...)"
    unitInput.value = product.variants.unit || ""

    const variantsContainer = document.createElement("div")

    const createVariantRow = (val = "", pr = "") => {
        const row = document.createElement("div")
        row.className = "variant-row"

        row.innerHTML = `
            <input class="v-val" value="${val}" placeholder="Значение">
            <input class="v-price" value="${pr}" placeholder="Цена">
            <button class="remove-btn">❌</button>
        `

        row.querySelector(".remove-btn").onclick = () => row.remove()
        return row
    }

    product.variants.value.forEach((v, i) => {
        variantsContainer.appendChild(
            createVariantRow(v, product.variants.price[i])
        )
    })

    const addVariantBtn = document.createElement("button")
    addVariantBtn.textContent = "+ вариант"
    addVariantBtn.onclick = () =>
        variantsContainer.appendChild(createVariantRow())

    variantsBlock.appendChild(title("Варианты"))
    variantsBlock.appendChild(unitInput)
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

    product.characteristics.forEach(c => {
        attrContainer.appendChild(createAttrRow(c.key, c.value))
    })

    const addAttrBtn = document.createElement("button")
    addAttrBtn.textContent = "+ характеристика"
    addAttrBtn.onclick = () =>
        attrContainer.appendChild(createAttrRow())

    attrBlock.appendChild(title("Характеристики"))
    attrBlock.appendChild(attrContainer)
    attrBlock.appendChild(addAttrBtn)
    content.appendChild(attrBlock)

    // ======================
    // 📂 КАТЕГОРИИ
    // ======================
    const catBlock = document.createElement("div")
    catBlock.className = "block"

    const catContainer = document.createElement("div")
    catContainer.className = "categories-grid"

    allSubcategories.forEach(sub => {
        const label = document.createElement("label")
        label.className = "cat-label"

        const checkbox = document.createElement("input")
        checkbox.type = "checkbox"
        checkbox.value = sub.slug
        checkbox.className = "cat-checkbox"

        if (product.categories.includes(sub.slug)) {
            checkbox.checked = true
        }

        label.appendChild(checkbox)
        label.append(` ${sub.name}`)
        catContainer.appendChild(label)
    })

    catBlock.appendChild(title("Подкатегории"))
    catBlock.appendChild(catContainer)
    content.appendChild(catBlock)

    // ======================
    // 💾 ACTIONS (SAVE + DELETE)
    // ======================
    const actions = document.createElement("div")
    actions.className = "modal-actions"

const saveBtn = document.createElement("button")
saveBtn.textContent = isCreate ? "Создать" : "Сохранить"
saveBtn.className = "save-btn"
saveBtn.type = "button"

const deleteBtn = document.createElement("button")
deleteBtn.textContent = "Удалить товар"
deleteBtn.className = "delete-btn"
deleteBtn.type = "button"

    // ---------------------
    // SAVE
    // ---------------------
saveBtn.onclick = async (e) => {

    e.preventDefault()

    // ======================
    // ПРОВЕРКИ
    // ======================

    if (!nameInput.value.trim()) {
        alert("Введите название")
        return
    }

    if (!desc.value.trim()) {
        alert("Введите описание")
        return
    }

    // фото
    const photos =
        photosContainer.querySelectorAll(".photo-wrapper")

    if (photos.length === 0) {
        alert("Добавьте фото")
        return
    }

    // варианты
    const variantRows =
        variantsContainer.querySelectorAll(".variant-row")

    if (variantRows.length === 0) {
        alert("Добавьте вариант")
        return
    }

    let hasVariant = false

    for (const row of variantRows) {

        const v =
            row.querySelector(".v-val").value.trim()

        const p =
            row.querySelector(".v-price").value.trim()

        // одно поле заполнено — второе обязательно
        if ((v && !p) || (!v && p)) {
            alert("Заполните вариант полностью")
            return
        }

        if (v && p) {
            hasVariant = true
        }
    }

    if (!hasVariant) {
        alert("Добавьте хотя бы один вариант")
        return
    }

    // подкатегории
    const checked =
        catContainer.querySelectorAll(".cat-checkbox:checked")

    if (checked.length === 0) {
        alert("Выберите подкатегорию")
        return
    }

    // ======================
    // FORM DATA
    // ======================

    const formData = new FormData()

    formData.append("name", nameInput.value)
    formData.append("description", desc.value)
    formData.append("variants_unit", unitInput.value)

    const existing = []

    photosContainer.querySelectorAll(".photo-wrapper").forEach(el => {

        if (el._file) {
            formData.append("newPhotos", el._file)
        }

        if (el.dataset.path) {
            existing.push(el.dataset.path)
        }
    })

    formData.append(
        "existingPhotos",
        JSON.stringify(existing)
    )

    // структура ОСТАЛАСЬ прежней
    const vars = {
        Unit: unitInput.value,
        Value: [],
        Price: []
    }

    variantsContainer.querySelectorAll(".variant-row").forEach(row => {

        const v =
            row.querySelector(".v-val").value.trim()

        const p =
            row.querySelector(".v-price").value.trim()

        if (v && p) {
            vars.Value.push(v)
            vars.Price.push(p)
        }
    })

    formData.append(
        "variants",
        JSON.stringify(vars)
    )

    const attrs = []

    attrContainer.querySelectorAll(".attr-row").forEach(row => {

        const k =
            row.querySelector(".a-key").value.trim()

        const v =
            row.querySelector(".a-val").value.trim()

        if (k) {
            attrs.push({
                key: k,
                value: v
            })
        }
    })

    formData.append(
        "characteristics",
        JSON.stringify(attrs)
    )

    const selected = []

    checked.forEach(cb => {
        selected.push(cb.value)
    })

    formData.append(
        "subcategories",
        JSON.stringify(selected)
    )

    if (!isCreate) {
        formData.append("id", product.id)
    }

    const url =
        isCreate
            ? "/admin/create_product"
            : "/admin/update_product"

    const res = await fetch(url, {
    method: "POST",
    body: formData
})


// ======================
// ПРОВЕРКА РАЗМЕРА
// ======================

let totalSize = 0

photosContainer.querySelectorAll(".photo-wrapper").forEach(el => {

    if (el._file) {
        totalSize += el._file.size
    }
})

const maxSize = 100 * 1024 * 1024

if (totalSize > maxSize) {

    const mb =
        (totalSize / 1024 / 1024).toFixed(2)

    alert(`Размер фото слишком большой: ${mb} MB`)
    return
}
const text = await res.text()

if (res.ok) {
    alert(isCreate ? "Создано!" : "Сохранено!")
    location.reload()
} else {
    console.error(text)
    alert(text || "Ошибка сохранения")
}
}

    // ---------------------
    // DELETE
    // ---------------------
    deleteBtn.onclick = async () => {
        const ok = confirm("Вы действительно хотите удалить товар?")
        if (!ok) return

        const res = await fetch(`/admin/delete_product/${product.id}`, {
            method: "DELETE"
        })

        if (res.ok) {
            alert("Товар удалён")
            modal.classList.remove("active")
            location.reload()
        } else {
            alert("Ошибка удаления")
        }
    }

    actions.appendChild(saveBtn)
    if (!isCreate) actions.appendChild(deleteBtn)

    content.appendChild(actions)
}
async function getProductsInfo(card) {
    const res = await fetch(`/api/product_info/${card.dataset.slug}`)
    const data = await res.json()

    const categories = await getCategories()

const allSubcategories = categories.flatMap(c =>
    (c.mininav || []).map(m => ({
        slug: m.slug,
        name: `${c.name} > ${m.name}`
    }))
)

    const product = normalizeProduct(card, data)

    drawModal(product, allSubcategories)

    document.getElementById("modal").classList.add("active")
}
function normalizeProduct(card, data) {
    return {
        id: card.dataset.id,
        slug: card.dataset.slug,

        name: data.Name || "",
        description: data.Decscription || "",

        photos: data.Photo || [],

        variants: {
            unit: data.Variants?.Unit || "",
            value: data.Variants?.Value || [],
            price: data.Variants?.Price || []
        },

        characteristics: data.Characteristic || [],
        categories: data.Categories || []
    }
}

function addProduct() {
    refreshProducts(true)
}

document.getElementById("productCards").addEventListener("click", (e) => {
    const card = e.target.closest(".product-card")
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


function createEmptyProduct() {
    return {
        id: null,
        slug: "",
        name: "",
        description: "",
        photos: [],
        variants: {
            unit: "",
            value: [],
            price: []
        },
        characteristics: [],
        categories: []
    }
}

document.getElementById("addButton").addEventListener("click", async () => {
    const categories = await getCategories()

    const allSubcategories = categories.flatMap(c =>
        (c.mininav || []).map(m => ({
            slug: m.slug,
            name: `${c.name} > ${m.name}`
        }))
    )

    const emptyProduct = createEmptyProduct()

    drawModal(emptyProduct, allSubcategories, true)

    document.getElementById("modal").classList.add("active")
})