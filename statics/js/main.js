import { initCatalog } from "./catalog.js";
import { initProduct } from "./product.js";

const active = window.__ACTIVE_PRODUCT__;
const isProductPage = window.location.pathname.startsWith("/product/");

// 🔥 1. СНАЧАЛА открываем модалку (если есть SSR)
if (isProductPage && active) {
    initProduct(true); // передаём флаг SSR
}

// 🔥 2. потом уже грузим каталог
document.addEventListener("DOMContentLoaded", async () => {
    await initCatalog();

    // 🔥 если не SSR — обычная инициализация
    if (!(isProductPage && active)) {
        initProduct(false);
    }
});