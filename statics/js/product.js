import { resetCatalog } from "./catalog.js";
let currentPhotoIndex = 0;

import { getLastListUrl } from "./router.js";
import { loadProducts } from "./catalog.js";

let currentPhoto = null;
let previousPage = "/";
let isAnimating = false; // Блокировка переключений во время анимации

// Сохраняем ссылку на обработчик клавиш, чтобы удалять его при закрытии
let globalKeyHandler = null; 

export function initProduct(isSSR = false) {
    const active = window.__ACTIVE_PRODUCT__;

    if (isSSR && active) {
        openProductModal(active);
        history.replaceState({}, "", `/product/${active.Url}`);
    }

    document.querySelector(".product-close")
        ?.addEventListener("click", closeProductModal);
}

export function openProductModal(product, prev = "/") {
    previousPage = prev;

    const modal = document.getElementById("product-modal");
    const photos = product.Photo || [];

    currentPhotoIndex = 0;
    isAnimating = false;

    // Сбрасываем возможные CSS-классы анимации на картинке
    const mainImage = document.getElementById("main-image");
    if (mainImage) mainImage.className = "modal-main-image";

    currentPhoto = photos[0] || "";
    setMainImage(currentPhoto);

    document.getElementById("modal-title").textContent = product.Name || "";
    document.getElementById("description").textContent = product.Description || "";

    renderThumbnails(photos);
    renderCharacteristics(product.Characteristic);
    renderVariants(product.Variants);
    setupPhotoNavigation(photos);

    modal.style.display = "block";
}

function closeProductModal() {
    document.getElementById("product-modal").style.display = "none";
    
    // Удаляем слушатель клавиатуры, чтобы он не копился при повторных открытиях
    if (globalKeyHandler) {
        document.removeEventListener("keydown", globalKeyHandler);
        globalKeyHandler = null;
    }

    const lastUrl = getLastListUrl();

    if (lastUrl === "/") {
        window.location.href = "/";
        return;
    }

    history.pushState({}, "", lastUrl);

    const path = lastUrl.split("/");
    const type = path[1] || "parfume";
    const sub = path[2] || "";

    resetCatalog();
    loadProducts(type, sub, false);
}

function setMainImage(src) {
    if (!src) return;
    if (!src.startsWith("/")) src = "/" + src;
    document.getElementById("main-image").src = src;
    currentPhoto = src;
}

function renderThumbnails(photos) {
    const container = document.querySelector(".modal-thumbnails");

    container.innerHTML = photos.map((p, i) => {
        const src = p.startsWith("/") ? p : "/" + p;
        return `<img class="modal-thumbnail ${i === 0 ? "active" : ""}" 
                     src="${src}" data-src="${src}">`;
    }).join("");

    container.onclick = (e) => {
        if (isAnimating) return;
        const img = e.target.closest(".modal-thumbnail");
        if (!img) return;

        const imgs = Array.from(container.querySelectorAll(".modal-thumbnail"));
        const targetIndex = imgs.indexOf(img);
        
        if (targetIndex === currentPhotoIndex) return;

        // Определяем направление анимации по индексу кликнутой миниатюры
        const direction = targetIndex > currentPhotoIndex ? "next" : "prev";
        currentPhotoIndex = targetIndex;

        // Вызываем глобальную функцию смены с анимацией (она объявлена внутри setupPhotoNavigation)
        if (window.__triggerPhotoChange) {
            window.__triggerPhotoChange(direction);
        }
    };
}

function syncThumbnailActive() {
    const thumbs = document.querySelectorAll(".modal-thumbnail");
    thumbs.forEach((t, i) => {
        t.classList.toggle("active", i === currentPhotoIndex);
    });
}

function renderCharacteristics(data) {
    const el = document.getElementById("characteristics");

    if (!data?.length) {
        el.innerHTML = "";
        return;
    }

    el.innerHTML = `
        <table class="custom-table">
            <tbody>
                ${data.map(i => `
                    <tr>
                        <td>${i.key}</td>
                        <td>${i.value}</td>
                    </tr>
                `).join("")}
            </tbody>
        </table>
    `;
}

function renderVariants(v) {
    const select = document.getElementById("volume-select");
    if (!select || !v) return;

    const values = v.Value || [];
    const unitText = v.Unit && v.Unit.trim() !== "" ? ` ${v.Unit}` : "";

    select.innerHTML = values.map((val, i) => `
        <option value="${val}" data-price="${v.Price[i]}" data-id="${v.Id[i]}">
            ${val}${unitText} — ${v.Price[i]} ₽
        </option>
    `).join("");

    if (values.length === 1) {
        select.style.display = "none";
        updateBuyButton(v.Price[0]);
    } else {
        select.style.display = "block";

        const setPrice = () => {
            const selected = select.options[select.selectedIndex];
            updateBuyButton(selected.dataset.price);
        };

        setPrice();
        select.onchange = setPrice;
    }
}

function updateBuyButton(price) {
    const btn = document.getElementById("submit");
    if (!btn) return;

    btn.textContent = `Добавить в корзину — ${price} ₽`;
}

function setupPhotoNavigation(photos) {
    const modal = document.querySelector(".modal-main-container");
    const mainImage = document.getElementById("main-image");

    if (!modal || !mainImage || !photos?.length) return;

    // Функция, выполняющая красивый сдвиг и подмену src
    const changePhotoWithAnimation = (direction) => {
        if (isAnimating) return;
        isAnimating = true;

        // 1. Уводим старое фото в сторону
        mainImage.classList.add(direction === "next" ? "slide-left-out" : "slide-right-out");

        // 2. Ждем завершения первой фазы анимации (200ms)
        setTimeout(() => {
            let src = photos[currentPhotoIndex];
            if (!src.startsWith("/")) src = "/" + src;
            
            mainImage.src = src;
            currentPhoto = src;

            // Сбрасываем старый класс и готовим к влёту с нужной стороны (без анимации `transition: none`)
            mainImage.className = "modal-main-image";
            mainImage.classList.add(direction === "next" ? "slide-left-in-prepare" : "slide-right-in-prepare");

            // Форсируем перерисовку DOM браузером (триггер)
            mainImage.getBoundingClientRect();

            // Включаем синхронизацию миниатюр
            syncThumbnailActive();

            // Убираем подготовительный класс, запуская плавное возвращение в центр
            mainImage.classList.remove("slide-left-in-prepare", "slide-right-in-prepare");

            // Разрешаем следующий клик после полного завершения анимации влёта
            setTimeout(() => {
                isAnimating = false;
            }, 200);

        }, 200);
    };

    // Экспортируем функцию в window, чтобы renderThumbnails мог её вызывать
    window.__triggerPhotoChange = changePhotoWithAnimation;

    const nextPhoto = () => {
        if (currentPhotoIndex < photos.length - 1) {
            currentPhotoIndex++;
            changePhotoWithAnimation("next");
        }
    };

    const prevPhoto = () => {
        if (currentPhotoIndex > 0) {
            currentPhotoIndex--;
            changePhotoWithAnimation("prev");
        }
    };

    // ===== КЛИК ПО ЛЕВОЙ / ПРАВОЙ ЧАСТИ =====
    modal.onclick = (e) => {
        if (isAnimating) return;
        const rect = modal.getBoundingClientRect();
        const x = e.clientX - rect.left;

        if (x < rect.width / 2) {
            prevPhoto();
        } else {
            nextPhoto();
        }
    };

    // ===== КЛАВИАТУРА =====
    if (globalKeyHandler) {
        document.removeEventListener("keydown", globalKeyHandler);
    }
    
    globalKeyHandler = (e) => {
        if (isAnimating) return;
        if (e.key === "ArrowLeft") prevPhoto();
        if (e.key === "ArrowRight") nextPhoto();
    };

    document.addEventListener("keydown", globalKeyHandler);

    // ===== СВАЙП (МОБИЛКА) =====
    let startX = 0;

    modal.ontouchstart = (e) => {
        startX = e.touches[0].clientX;
    };

    modal.ontouchend = (e) => {
        if (isAnimating) return;
        const endX = e.changedTouches[0].clientX;
        const diff = startX - endX;

        if (Math.abs(diff) < 40) return;

        if (diff > 0) {
            nextPhoto();
        } else {
            prevPhoto();
        }
    };
}
