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
let currentOffset = 0;
const limit = 20;
let isLoading = false;
let allLoaded = false; // Флаг, чтобы не спамить запросами, если товары кончились

export async function loadProducts(type, subcategory, append = false) {
    if (isLoading || allLoaded) return;
    isLoading = true;

    try {
        const res = await fetch(
            `/api/get_products?type=${type}&category=${subcategory}&limit=${limit}&offset=${currentOffset}`
        );

        const data = await res.json();
        const container = document.getElementById("products-container");
        const buttonWrapper = document.querySelector(".button-container");

        // Если данных нет — скрываем кнопку и выходим
        if (!data || data.length === 0) {
            allLoaded = true;
            if (buttonWrapper) buttonWrapper.style.display = "none";
            return;
        }

        // Рендерим карточки
        const html = data.map(p => `
            <div class="product_card" data-slug="${p.Url}">
                <a href="/product/${p.Url}">
                    <img src="/${p.MainPhoto}" class="product_image">
                    <p class="product_name_p">${p.Name}</p>
                    <div class="product_real_price">${p.Price} ₽</div>
                    <div class="product_price">Подробнее</div>
                </a>
            </div>
        `).join("");

        if (append) {
            container.insertAdjacentHTML('beforeend', html);
        } else {
            container.innerHTML = html;
        }

        // Если пришло меньше лимита — товары закончились
        if (data.length < limit) {
            allLoaded = true;
            if (buttonWrapper) buttonWrapper.style.display = "none";
        }

    } catch (err) {
        console.error("Ошибка загрузки:", err);
    } finally {
        isLoading = false;
    }
}

// Логика инициализации
document.addEventListener("DOMContentLoaded", () => {
    // 1. Сразу загружаем первую партию (offset 0)
    loadProducts('parfume', '', false);

    const loadMoreBtn = document.getElementById("load-more");
    if (loadMoreBtn) {
        
        const fetchNext = () => {
            if (!isLoading && !allLoaded) {
                currentOffset += limit; // Прибавляем ТОЛЬКО при загрузке следующей части
                loadProducts('parfume', '', true);
            }
        };

        // Клик
        loadMoreBtn.addEventListener("click", (e) => {
            e.preventDefault();
            fetchNext();
        });

        // Скролл (Intersection Observer)
        const observer = new IntersectionObserver((entries) => {
            // Сработает только если кнопка видна И первая партия уже загружена
            if (entries[0].isIntersecting && currentOffset >= 0 && !isLoading) {
                // Небольшая задержка, чтобы не было ложных срабатываний при старте
                if (document.querySelectorAll('.product_card').length > 0) {
                   fetchNext();
                }
            }
        }, { threshold: 0.1 });

        observer.observe(loadMoreBtn);
    }
});
