document.getElementById("addCategoryButton")
    .addEventListener("click", async () => {

    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")

    const categories = await getCategories()

    renderCategoriesPage(categories)

    modal.classList.add("active")
})


// ======================================
// РЕНДЕР СПИСКА КАТЕГОРИЙ
// ======================================
function renderCategoriesPage(categories) {

    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")

    content.innerHTML = `
        <div class="category-modal">

            <div class="category-header">

                <h2>Категории</h2>

                <button
                    id="createCategoryBtn"
                    class="save-btn"
                >
                    + Новая категория
                </button>

            </div>

            <div
                id="categoriesList"
                class="categories-list"
            ></div>

        </div>
    `

    const list = document.getElementById("categoriesList")

    categories.forEach(cat => {

        const item = document.createElement("div")

        item.className = "category-item"

        item.innerHTML = `
            <div class="category-info">

                <div class="category-name">
                    📁 ${cat.name}
                </div>

                <div class="category-slug">
                    ${cat.slug}
                </div>

            </div>

            <div class="category-actions">

                <button
                    class="open-category-btn"
                >
                    Открыть
                </button>

                <button
                    class="edit-category-btn"
                >
                    ✏️
                </button>

            </div>
        `

        // ======================
        // OPEN CATEGORY
        // ======================
        item.querySelector(".open-category-btn")
            .addEventListener("click", () => {

            renderSubcategoriesPage(cat)
        })

        // ======================
        // EDIT CATEGORY
        // ======================
        item.querySelector(".edit-category-btn")
            .addEventListener("click", () => {

            openCategoryEditor(cat)
        })

        list.appendChild(item)
    })

    // ======================
    // CREATE CATEGORY
    // ======================
    document.getElementById("createCategoryBtn")
        .addEventListener("click", () => {

        openCategoryEditor(null)
    })
}


// ======================================
// СТРАНИЦА ПОДКАТЕГОРИЙ
// ======================================
function renderSubcategoriesPage(category) {

    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")

    const subcategories = category.mininav || []

    content.innerHTML = `
        <div class="category-modal">

            <div class="category-header">

                <div>

                    <button
                        id="backToCategoriesBtn"
                        class="back-btn"
                    >
                        ← Назад
                    </button>

                    <h2>
                        ${category.name}
                    </h2>

                </div>

            </div>

            <div
                id="subcategoriesList"
                class="categories-list"
            ></div>

        </div>
    `

    // ======================
    // BACK
    // ======================
    document.getElementById("backToCategoriesBtn")
        .addEventListener("click", async () => {

        const categories = await getCategories()

        renderCategoriesPage(categories)
    })

    

    const list =
        document.getElementById("subcategoriesList")

    // ======================
    // EMPTY
    // ======================
    if (subcategories.length === 0) {

        list.innerHTML = `
            <div class="empty-block">
                Подкатегорий пока нет
            </div>
        `

        return
    }

    // ======================
    // RENDER SUBCATEGORIES
    // ======================
    subcategories.forEach(sub => {

        const item = document.createElement("div")

        item.className = "category-item"

        item.innerHTML = `
            <div class="category-info">

                <div class="category-name">
                    📄 ${sub.name}
                </div>

                <div class="category-slug">
                    ${sub.slug}
                </div>

            </div>

            <div class="category-actions">

                <button
                    class="edit-category-btn"
                >
                    ✏️
                </button>

            </div>
        `

        item.querySelector(".edit-category-btn")
            .addEventListener("click", () => {

            openCategoryEditor({
                ...sub,
                parent_slug: category.slug
            })
        })

        list.appendChild(item)
    })
}

async function openCategoryEditor(category = null) {

    const categories = await getCategories()

    const modal = document.getElementById("modal")
    const content = modal.querySelector(".modal-content")

    const isEdit = !!category
    const isSubcategory = !!category?.parent_slug

    content.innerHTML = `
        <div class="category-editor">

            <h2>
                ${isEdit
                    ? "Редактирование"
                    : "Создание"}
            </h2>

            <div class="block">

                <div class="label">
                    Тип
                </div>

                <select
                    id="categoryType"
                    class="input"
                >
                    <option value="category">
                        Категория
                    </option>

                    <option value="subcategory">
                        Подкатегория
                    </option>
                </select>

            </div>

            <div
                class="block"
                id="parentCategoryBlock"
                style="display:none;"
            >

                <div class="label">
                    Родительская категория
                </div>

                <select
                    id="parentCategorySelect"
                    class="input"
                >

                    <option value="">
                        Выберите категорию
                    </option>

                    ${
                        categories.map(c => `
                            <option value="${c.slug}">
                                ${c.name}
                            </option>
                        `).join("")
                    }

                </select>

            </div>

            <div class="block">

                <div class="label">
                    Название
                </div>

                <input
                    id="categoryNameInput"
                    class="input"
                    value="${category?.name || ""}"
                    placeholder="Название"
                >

            </div>

            <div class="editor-actions">

                <button
                    id="saveCategoryBtn"
                    class="save-btn"
                >
                    ${isEdit
                        ? "Сохранить"
                        : "Создать"}
                </button>

                ${
                    isEdit
                    ? `
                    <button
                        id="deleteCategoryBtn"
                        class="delete-btn"
                    >
                        Удалить
                    </button>
                    `
                    : ""
                }

            </div>

        </div>
    `

    const typeSelect =
        document.getElementById("categoryType")

    const parentBlock =
        document.getElementById("parentCategoryBlock")

    const parentSelect =
        document.getElementById("parentCategorySelect")

    function updateTypeUI() {

        if (typeSelect.value === "subcategory") {

            parentBlock.style.display = "block"

        } else {

            parentBlock.style.display = "none"
        }
    }

    typeSelect.addEventListener(
        "change",
        updateTypeUI
    )

    // ======================
    // EDIT MODE
    // ======================
    if (isSubcategory) {

        typeSelect.value = "subcategory"

        parentSelect.value =
            category.parent_slug || ""

    } else {

        typeSelect.value = "category"
    }

    updateTypeUI()

    // ======================
    // SAVE
    // ======================
    document.getElementById("saveCategoryBtn")
        .addEventListener("click", async () => {

        const name = document
            .getElementById("categoryNameInput")
            .value
            .trim()

        if (!name) {
            alert("Введите название")
            return
        }

        const payload = {
            name,
            type: typeSelect.value
        }

        if (
            typeSelect.value === "subcategory"
        ) {

            if (!parentSelect.value) {

                alert("Выберите категорию")
                return
            }

            payload.parent_slug =
                parentSelect.value
        }

        if (isEdit) {
            payload.slug = category.slug
        }

        const url = isEdit
            ? "/admin/update_category"
            : "/admin/create_category"

        const res = await fetch(url, {

            method: "POST",

            headers: {
                "Content-Type": "application/json"
            },

            body: JSON.stringify(payload)
        })

        if (res.ok) {

            alert(
                isEdit
                    ? "Обновлено"
                    : "Создано"
            )

            location.reload()

        } else {

            alert("Ошибка")
        }
    })

    // ======================
// DELETE
// ======================
if (isEdit) {

    document.getElementById("deleteCategoryBtn")
        .addEventListener("click", async () => {

        const ok = confirm(
            `Удалить "${category.name}"?`
        )

        if (!ok) return

        const payload = {
            slug: category.slug,

            type: isSubcategory
                ? "subcategory"
                : "category"
        }

        const res = await fetch(
            `/admin/delete_category`,
            {
                method: "DELETE",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify(payload)
            }
        )

        if (res.ok) {

            alert("Удалено")

            location.reload()

        } else {

            const err = await res.json()

            alert(
                err.error || "Ошибка удаления"
            )
        }
    })
}
}