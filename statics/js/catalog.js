import { openProductModal } from "./product.js";

let previousPage = "/";

// -------------------------
// INIT
// -------------------------
export async function initCatalog() {
    const category = window.location.pathname.split("/")[1];

    const miniContainer = document.getElementById("mini-navbar");

    const resCats = await fetch(`/api/mini_categories?category=${category}`);
    const categories = await resCats.json();

    const allButton = `
        <a href="/${category}" data-slug="">
            Все
        </a>
    `;

    miniContainer.innerHTML =
        allButton +
        categories.map(item =>
            `<a href="/${category}/${item.slug}" data-slug="${item.slug}">
                ${item.name}
            </a>`
        ).join("");

    await loadProducts(category, "");

    // фильтр
    miniContainer.addEventListener("click", (e) => {
        const link = e.target.closest("a");
        if (!link) return;

        e.preventDefault();

        const sub = link.dataset.slug;

        if (!sub) {
            history.pushState({}, "", `/${category}`);
            loadProducts(category, "");
            return;
        }

        history.pushState({}, "", `/${category}/${sub}`);
        loadProducts(category, sub);
    });

    // клик по карточке
    document.addEventListener("click", async (e) => {
        const card = e.target.closest(".product_card");
        if (!card) return;

        e.preventDefault();

        const slug = card.dataset.slug;
        if (!slug) return;

        previousPage = window.location.pathname;

        history.pushState({}, "", `/product/${slug}`);

        const res = await fetch(`/api/product/${slug}`);
        const product = await res.json();

        openProductModal(product, previousPage);
    });
}

// -------------------------
// LOAD PRODUCTS
// -------------------------
async function loadProducts(type, subcategory) {
    const res = await fetch(
        `/api/get_products?type=${type}&category=${subcategory}&limit=20&offset=0`
    );

    const data = await res.json();

    const container = document.getElementById("products-container");

    container.innerHTML = data.map(p => `
        <div class="product_card" data-slug="${p.Url}">
            <a href="/product/${p.Url}">
                <img src="/${p.MainPhoto}" class="product_image">
                <p class="product_name_p">${p.Name}</p>
                <div class="product_real_price">${p.Price} ₽</div>
                <div class="product_price">Подробнее</div>
            </a>
        </div>
    `).join("");
}