import {resetCatalog } from "./catalog.js";
let currentPhotoIndex = 0;

import { getLastListUrl } from "./router.js";
import { loadProducts } from "./catalog.js";

let currentPhoto = null;
let previousPage = "/";

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
        const img = e.target.closest(".modal-thumbnail");
        if (!img) return;

        const imgs = Array.from(container.querySelectorAll(".modal-thumbnail"));

        currentPhotoIndex = imgs.indexOf(img);

        syncThumbnailActive();
        setMainImage(img.dataset.src);
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

    if (!modal || !photos?.length) return;

    const update = () => {
        const src = photos[currentPhotoIndex];
        setMainImage(src);
        syncThumbnailActive();
    };

    // ===== КЛИК ПО ЛЕВОЙ / ПРАВОЙ ЧАСТИ =====
    modal.onclick = (e) => {
        const rect = modal.getBoundingClientRect();
        const x = e.clientX - rect.left;

        if (x < rect.width / 2) {
            currentPhotoIndex = Math.max(0, currentPhotoIndex - 1);
        } else {
            currentPhotoIndex = Math.min(photos.length - 1, currentPhotoIndex + 1);
        }

        update();
    };

    // ===== КЛАВИАТУРА =====
    const keyHandler = (e) => {
        if (e.key === "ArrowLeft") {
            currentPhotoIndex = Math.max(0, currentPhotoIndex - 1);
            update();
        }

        if (e.key === "ArrowRight") {
            currentPhotoIndex = Math.min(photos.length - 1, currentPhotoIndex + 1);
            update();
        }
    };

    document.addEventListener("keydown", keyHandler);

    // ===== СВАЙП (МОБИЛКА) =====
    let startX = 0;

    modal.addEventListener("touchstart", (e) => {
        startX = e.touches[0].clientX;
    });

    modal.addEventListener("touchend", (e) => {
        const endX = e.changedTouches[0].clientX;
        const diff = startX - endX;

        if (Math.abs(diff) < 40) return;

        if (diff > 0) {
            currentPhotoIndex = Math.min(photos.length - 1, currentPhotoIndex + 1);
        } else {
            currentPhotoIndex = Math.max(0, currentPhotoIndex - 1);
        }

        update();
    });
}