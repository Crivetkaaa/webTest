
import { getLastListUrl } from "./router.js";
import { restoreRoute } from "./router.js";
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

    currentPhoto = photos[0] || "";
    setMainImage(currentPhoto);

    document.getElementById("modal-title").textContent = product.Name || "";
    document.getElementById("description").textContent = product.Description || "";

    renderThumbnails(photos);
    renderCharacteristics(product.Characteristic);
    renderVariants(product.Variants);

    modal.style.display = "block";
}


function closeProductModal() {
    document.getElementById("product-modal").style.display = "none";

    const lastUrl = getLastListUrl();

    // прямой заход на товар
    if (lastUrl === "/") {
        window.location.href = "/";
        return;
    }

    // обычный SPA возврат
    history.pushState({}, "", lastUrl);

    const path = lastUrl.split("/");
    const type = path[1] || "parfume";
    const sub = path[2] || "";

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

        container.querySelectorAll("img").forEach(i => i.classList.remove("active"));
        img.classList.add("active");

        setMainImage(img.dataset.src);
    };
}

function renderCharacteristics(data) {
    const el = document.getElementById("characteristics")

    if (!data?.length) {
        el.innerHTML = ""
        return
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
    `
}

function renderVariants(v) {
    const select = document.getElementById("volume-select");
    if (!select || !v) return;

    const values = v.Value || [];
    // Проверяем, есть ли юнит, если нет — строка пустая
    const unitText = v.Unit && v.Unit.trim() !== "" ? ` ${v.Unit}` : "";

    select.innerHTML = values.map((val, i) => `
        <option value="${val}" data-price="${v.Price[i]}" data-id="${v.Id[i]}">
            ${val}${unitText} — ${v.Price[i]} ₽
        </option>
    `).join("");

    const oldInfo = document.getElementById("single-variant-info");
    if (oldInfo) oldInfo.remove();

    if (values.length === 1) {
        select.style.display = "none";

        const infoBlock = document.createElement("div");
        infoBlock.id = "single-variant-info";
        infoBlock.className = "text-base font-semibold text-gray-800 mb-3 bg-gray-50 p-2 rounded-lg border border-gray-100 inline-block";
        
        // Корректный вывод для единственного варианта без юнита
        infoBlock.textContent = `${v.Price} ₽`;

        select.parentNode.insertBefore(infoBlock, select.nextSibling);
    } else {
        select.style.display = "block";
    }
}

