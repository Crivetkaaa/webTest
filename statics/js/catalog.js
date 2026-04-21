import { openProductModal } from "./product.js";

let previousPage = "/";
let currentOffset = 0;
const limit = 20;
let isLoading = false;
let allLoaded = false;

// Определяем текущую основную категорию из URL сразу (например, 'clothes' или 'parfume')
const currentMainCategory = window.location.pathname.split("/")[1] || 'parfume';

// -------------------------
// INIT
// -------------------------
export async function initCatalog() {
    const miniContainer = document.getElementById("mini-navbar");
    if (!miniContainer) return;

    const resCats = await fetch(`/api/mini_categories?category=${currentMainCategory}`);
    const categories = await resCats.json();

    const allButton = `
        <a href="/${currentMainCategory}" data-slug="">
            Все
        </a>
    `;

    miniContainer.innerHTML = `
    <div class="category">
        ${allButton}
        ${categories.map(item => `
            <a href="/${currentMainCategory}/${item.slug}" data-slug="${item.slug}">
                ${item.name}
            </a>
        `).join("")}
    </div>
`;
    // фильтр по мини-категориям
    miniContainer.addEventListener("click", (e) => {
        const link = e.target.closest("a");
        if (!link) return;

        e.preventDefault();

        const sub = link.dataset.slug;

        // ВАЖНО: Сбрасываем состояние пагинации перед загрузкой новой категории
        currentOffset = 0;
        allLoaded = false;
        
        const loadMoreBtn = document.getElementById("load-more");
        if (loadMoreBtn) loadMoreBtn.style.display = "block";

        if (!sub) {
            history.pushState({}, "", `/${currentMainCategory}`);
            loadProducts(currentMainCategory, "");
        } else {
            history.pushState({}, "", `/${currentMainCategory}/${sub}`);
            loadProducts(currentMainCategory, sub);
        }
    });

    // клик по карточке (модальное окно)
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
export async function loadProducts(type, subcategory, append = false) {
    if (isLoading || (allLoaded && append)) return;
    isLoading = true;

    try {
        const res = await fetch(
            `/api/get_products?type=${type}&category=${subcategory}&limit=${limit}&offset=${currentOffset}`
        );

        const data = await res.json();
        const container = document.getElementById("products-container");
        const loadMoreBtn = document.getElementById("load-more");

        // Если данных нет совсем
        if (!data || data.length === 0) {
            allLoaded = true;
            if (loadMoreBtn) loadMoreBtn.style.display = "none";
            if (!append) container.innerHTML = "<p class='text-center w-full py-10'>Товары не найдены</p>";
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

        // Если пришло меньше лимита — товары в этой категории закончились
        if (data.length < limit) {
            allLoaded = true;
            if (loadMoreBtn) loadMoreBtn.style.display = "none";
        } else {
            if (loadMoreBtn) loadMoreBtn.style.display = "block";
        }

    } catch (err) {
        console.error("Ошибка загрузки:", err);
    } finally {
        isLoading = false;
    }
}

// -------------------------
// DOM CONTENT LOADED
// -------------------------
document.addEventListener("DOMContentLoaded", () => {
    // Определяем подкатегорию из URL при загрузке (если есть)
    const pathParts = window.location.pathname.split("/");
    const subCatFromUrl = pathParts[2] || "";

    // Инициализируем меню и загружаем первую партию
    initCatalog();
    loadProducts(currentMainCategory, subCatFromUrl, false);

    const loadMoreBtn = document.getElementById("load-more");
    if (loadMoreBtn) {
        
        const fetchNext = () => {
            if (!isLoading && !allLoaded) {
                currentOffset += limit;
                // Всегда берем актуальную подкатегорию из URL
                const currentPathParts = window.location.pathname.split("/");
                const currentSubCat = currentPathParts[2] || "";
                loadProducts(currentMainCategory, currentSubCat, true);
            }
        };

        // Клик по кнопке
        loadMoreBtn.addEventListener("click", (e) => {
            e.preventDefault();
            fetchNext();
        });

        // Скролл (Intersection Observer)
        const observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && !isLoading && !allLoaded) {
                // Загружаем следующую партию только если на странице уже что-то есть
                if (document.querySelectorAll('.product_card').length > 0) {
                   fetchNext();
                }
            }
        }, { threshold: 0.1 });

        observer.observe(loadMoreBtn);
    }
});
