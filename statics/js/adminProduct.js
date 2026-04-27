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


function drawModal(data, product) {
    console.log(data)
    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")


    // ======================
    // 🏷 NAME (из карточки)
    // ======================
    const nameValue =
        product.querySelector(".name")?.textContent ||
        product.children?.[1]?.textContent || ""
    console.log("dsfsdfsdf", product.dataset.id)
    const nameInput = document.createElement("input")
    nameInput.value = nameValue
    nameInput.className = "input"
    nameInput.placeholder = "Название товара"
    content.innerHTML = ""
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

    const photos = document.createElement("div")
    photos.className = "photos"

    data.Photo.forEach(src => {
        const img = document.createElement("img")
        img.src = `/${src}`
        photos.appendChild(img)
    })

    photosBlock.appendChild(title("Фото"))
    photosBlock.appendChild(photos)

    // ======================
    // 📝 ОПИСАНИЕ
    // ======================
    const desc = document.createElement("textarea")
    desc.value = data.Decscription || ""
    desc.className = "input"

    // ======================
    // 📦 ВАРИАНТЫ
    // ======================
    const variantsBlock = document.createElement("div")
    variantsBlock.className = "block"

    const variantsContainer = document.createElement("div")

    for (let i = 0; i < data.Variants.Value.length; i++) {
        const row = document.createElement("div")
        row.className = "variant-row"

        const value = document.createElement("input")
        value.value = data.Variants.Value[i]
        value.placeholder = data.Variants.Unit

        const price = document.createElement("input")
        price.value = data.Variants.Price[i]
        price.placeholder = "Цена"

        const remove = document.createElement("button")
        remove.textContent = "❌"
        remove.onclick = () => row.remove()

        row.appendChild(value)
        row.appendChild(price)
        row.appendChild(remove)

        variantsContainer.appendChild(row)
    }

    const addVariantBtn = document.createElement("button")
    addVariantBtn.textContent = "+ вариант"
    addVariantBtn.onclick = () => {
        const row = document.createElement("div")
        row.className = "variant-row"

        row.innerHTML = `
            <input placeholder="${data.Variants.Unit}">
            <input placeholder="Цена">
        `

        variantsContainer.appendChild(row)
    }


    variantsBlock.appendChild(title("Варианты"))
    variantsBlock.appendChild(variantsContainer)
    variantsBlock.appendChild(addVariantBtn)

    // ======================
    // 🏷 ХАРАКТЕРИСТИКИ
    // ======================
    const attrBlock = document.createElement("div")
    attrBlock.className = "block"

    const attrContainer = document.createElement("div")

    data.Characteristic.forEach(c => {
        const row = document.createElement("div")
        row.className = "attr-row"

        const key = document.createElement("input")
        key.value = c.key

        const value = document.createElement("input")
        value.value = c.value

        const remove = document.createElement("button")
        remove.textContent = "❌"
        remove.onclick = () => row.remove()

        row.appendChild(key)
        row.appendChild(value)
        row.appendChild(remove)

        attrContainer.appendChild(row)
    })

    const addAttrBtn = document.createElement("button")
    addAttrBtn.textContent = "+ характеристика"
    addAttrBtn.onclick = () => {
        const row = document.createElement("div")
        row.className = "attr-row"

        row.innerHTML = `
            <input placeholder="ключ">
            <input placeholder="значение">
        `

        attrContainer.appendChild(row)
    }

    attrBlock.appendChild(title("Характеристики"))
    attrBlock.appendChild(attrContainer)
    attrBlock.appendChild(addAttrBtn)

    // ======================
    // 💾 КНОПКА
    // ======================
    const save = document.createElement("button")
    save.textContent = "Сохранить"
    save.className = "save-btn"

    // ======================
    // СБОРКА
    // ======================
    // content.append(modalName)
    content.appendChild(photosBlock)
    content.appendChild(title("Описание"))
    content.appendChild(desc)
    content.appendChild(variantsBlock)
    content.appendChild(attrBlock)
    content.appendChild(save)
}

async function getProductsInfo(product) {
    console.log(product)
    const product_info = await fetch(`/api/product_info/${product.dataset.slug}`)
    const data = await product_info.json()

    const modal = document.getElementById("modal")
    await drawModal(data, product)
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