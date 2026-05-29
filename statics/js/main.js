import { initCatalog, loadProducts } from "./catalog.js";
import { initProduct, openProductModal } from "./product.js";
import { initCart, updateCartUI } from "./cart.js";
import { setSearch } from "./catalog.js";

document.addEventListener("DOMContentLoaded", () => {
    const input = document.getElementById("search-input");
    const btn = document.getElementById("search-trigger");

    if (!input || !btn) return;

let debounce;

const runSearch = () => {

    clearTimeout(debounce);

    const value = input.value.trim();

    // если пусто — возвращаем каталог
    if (value.length === 0) {
        setSearch("");
        return;
    }

    // меньше 3 букв — ничего не делаем
    if (value.length < 3) {
        return;
    }

    // debounce 1 секунда
    debounce = setTimeout(() => {
        setSearch(value);
    }, 1000);
};

    input.addEventListener("input", runSearch);

    btn.addEventListener("click", () => {
        input.classList.toggle("active");
        input.focus();
    });

    // Enter тоже работает
    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            setSearch(input.value);
        }
    });
});
const active = window.__ACTIVE_PRODUCT__;
const isProductPage = window.location.pathname.startsWith("/product/");

document.addEventListener("DOMContentLoaded", async () => {

    // Каталог всегда
    await initCatalog();

    const pathParts = window.location.pathname.split("/");
    const subCatFromUrl = pathParts[2] || "";

    // первая загрузка товаров
    const type = window.location.pathname.split("/")[1] || "parfume";
    await loadProducts(type, subCatFromUrl, false);

    // корзина
    initCart();
    updateCartUI();

    // продукт (если не SSR)
    if (!(isProductPage && active)) {
        initProduct(false);
    }
});

// SSR продукт
if (isProductPage && active) {
    initProduct(true);
}

// реакция на роут
window.addEventListener("route:changed", () => {
    const type = window.location.pathname.split("/")[1] || "parfume";
    loadProducts(type, "", false);
});

window.addEventListener("popstate", async () => {
    const path = window.location.pathname;
    const modal = document.getElementById("product-modal");

    if (path.startsWith("/product/")) {
        const slug = path.split("/product/")[1];

        const res = await fetch(`/api/product/${slug}`);
        const product = await res.json();

        openProductModal(product);
        return;
    }

    modal.style.display = "none";

    const parts = path.split("/");
    const type = parts[1] || "parfume";
    const sub = parts[2] || "";

    await loadProducts(type, sub, false);
});