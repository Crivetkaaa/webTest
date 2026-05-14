import { openProductModal } from "./product.js";
import { setLastListUrl } from "./router.js";

// -------------------------
// STATE
// -------------------------
let searchQuery = "";
let currentSubCategory = "";

let currentOffset = 0;
const limit = 20;

let isLoading = false;
let allLoaded = false;

const currentMainCategory =
    window.location.pathname.split("/")[1] || "parfume";

// -------------------------
// SEARCH API
// -------------------------
export function setSearch(query) {
    const container = document.getElementById("products-container");
if (container) container.innerHTML = "";
    searchQuery = query.trim();

    currentOffset = 0;
    allLoaded = false;

    currentSubCategory = window.location.pathname.split("/")[2] || "";

    loadProducts(currentMainCategory, currentSubCategory, false);
}

// -------------------------
// INIT CATALOG
// -------------------------
export async function initCatalog() {
    const miniContainer = document.getElementById("mini-navbar");
    if (!miniContainer) return;

    const res = await fetch(`/api/mini_categories?category=${currentMainCategory}`);
    const categories = await res.json();

    if (categories != null) {
        if (categories.length > 1) {
            miniContainer.innerHTML = `
                <div class="category">
                    <a href="/${currentMainCategory}" data-slug="">Все</a>
                    ${categories.map(item => `
                        <a href="/${currentMainCategory}/${item.slug}" data-slug="${item.slug}">
                            ${item.name}
                        </a>
                    `).join("")}
                </div>
            `;
        }
    }

    // -------------------------
    // CATEGORY CLICK
    // -------------------------
    miniContainer.addEventListener("click", (e) => {
        const link = e.target.closest("a");
        if (!link) return;

        e.preventDefault();

        const sub = link.dataset.slug || "";

        currentSubCategory = sub;
        searchQuery = "";

        currentOffset = 0;
        allLoaded = false;

        history.pushState({}, "", sub
            ? `/${currentMainCategory}/${sub}`
            : `/${currentMainCategory}`
        );

        loadProducts(currentMainCategory, sub, false);
    });

    // -------------------------
    // PRODUCT CLICK
    // -------------------------
    document.addEventListener("click", async (e) => {
        const card = e.target.closest(".product_card");
        if (!card) return;

        e.preventDefault();

        const slug = card.dataset.slug;

        setLastListUrl(window.location.pathname);

        history.pushState(
            { type: "product", slug },
            "",
            `/product/${slug}`
        );

        const res = await fetch(`/api/product/${slug}`);
        const product = await res.json();

        openProductModal(product);
    });
}

// -------------------------
// LOAD PRODUCTS
// -------------------------
export async function loadProducts(type, subcategory, append = false) {
    if (isLoading || (allLoaded && append)) return;

    isLoading = true;

    const container = document.getElementById("products-container");
    const btn = document.getElementById("load-more");
    if (!container) {
    isLoading = false;
    return;
}

    try {
        const res = await fetch(
            `/api/get_products?type=${type}` +
            `&category=${subcategory}` +
            `&limit=${limit}` +
            `&offset=${currentOffset}` +
            `&search=${encodeURIComponent(searchQuery)}`
        );

        const data = await res.json();

        // -------------------------
        // EMPTY RESULT
        // -------------------------
        if (!data || !data.length) {
            allLoaded = true;
            if (btn) btn.style.display = "none";
            if (!append) container.innerHTML = "Товары не найдены";
            return;
        }

        // -------------------------
        // RESET ON NEW SEARCH / CATEGORY
        // -------------------------
        if (!append) {
            container.innerHTML = "";
        }

        // -------------------------
        // RENDER
        // -------------------------
        const html = data.map(p => `
            <div class="product_card" data-slug="${p.Url}">
                <a href="/product/${p.Url}">
                    <img src="/${p.MainPhoto}" class="product_image">

                    <p class="product_name_p">
                        ${p.Name}
                    </p>

                    <div class="product_real_price">
                        ${p.Price} ₽
                    </div>

                    <div class="product_price">
                        Подробнее
                    </div>
                </a>
            </div>
        `).join("");

        if (append) {
            container.insertAdjacentHTML("beforeend", html);
        } else {
            container.innerHTML = html;
        }

        // -------------------------
        // PAGINATION STATE
        // -------------------------
        if (data.length < limit) {
            allLoaded = true;
            if (btn) btn.style.display = "none";
        } else {
            if (btn) btn.style.display = "block";
        }

    } finally {
        isLoading = false;
    }
}

// -------------------------
// LOAD MORE SUPPORT
// -------------------------
export function loadNextPage() {
    if (isLoading || allLoaded) return;

    currentOffset += limit;

    loadProducts(
        currentMainCategory,
        currentSubCategory,
        true
    );
}