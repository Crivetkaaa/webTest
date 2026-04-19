let savedMainPhoto = null;
let currentPhoto = null;
let previousPage = "/";

// -------------------------
// INIT PRODUCT (SSR)
// -------------------------
export function initProduct(isSSR = false) {
    const active = window.__ACTIVE_PRODUCT__;

    // 🔥 SSR открытие
    if (isSSR && active) {
        openProductModal(active);
        history.replaceState({}, "", `/product/${active.Url}`);
    }

    // кнопка закрытия
    const btn = document.querySelector(".product-close");
    if (btn) {
        btn.addEventListener("click", closeProductModal);
    }
}

// -------------------------
// OPEN MODAL
// -------------------------
export function openProductModal(product, prev = "/") {
    previousPage = prev;

    const modal = document.getElementById("product-modal");

    const photos = product.Photo || [];

    savedMainPhoto = photos[0] || "";
    currentPhoto = savedMainPhoto;

    setMainImage(currentPhoto);

    document.getElementById("modal-title").textContent = product.Name || "";
    document.getElementById("description").textContent = product.Description || "";

    renderThumbnails(photos);
    renderCharacteristics(product.Characteristic);
    renderVariants(product.Variants);

    modal.style.display = "block";
}

// -------------------------
// CLOSE MODAL
// -------------------------
function closeProductModal() {
    const modal = document.getElementById("product-modal");
    modal.style.display = "none";

    const fromOwnSite =
        document.referrer && document.referrer.includes(location.host);

    const isReload =
        performance.getEntriesByType("navigation")[0]?.type === "reload";

    if (fromOwnSite && !isReload) {
        history.pushState({}, "", previousPage || "/");
        return;
    }

    window.location.href = "/";
}

// -------------------------
// MAIN IMAGE
// -------------------------
function setMainImage(src) {
    if (!src.startsWith("/")) src = "/" + src;

    currentPhoto = src;
    document.getElementById("main-image").src = src;
}

// -------------------------
// THUMBNAILS
// -------------------------
function renderThumbnails(photos) {
    const container = document.querySelector(".modal-thumbnails");

    container.innerHTML = photos.map((p, i) => {
        let src = p;
        if (!src.startsWith("/")) src = "/" + src;

        return `
            <img 
                class="modal-thumbnail ${i === 0 ? "active" : ""}"
                src="${src}"
                data-src="${src}"
            >
        `;
    }).join("");

    container.onclick = (e) => {
        const img = e.target.closest(".modal-thumbnail");
        if (!img) return;

        container.querySelectorAll(".modal-thumbnail")
            .forEach(el => el.classList.remove("active"));

        img.classList.add("active");

        setMainImage(img.dataset.src);
    };
}

// -------------------------
// CHARACTERISTICS
// -------------------------
function renderCharacteristics(data) {
    const container = document.getElementById("characteristics");

    if (!data || data.length === 0) {
        container.innerHTML = "";
        return;
    }

    container.innerHTML = `
        <table class="custom-table">
            <tbody>
                ${data.map(item => `
                    <tr>
                        <td>${item.key}</td>
                        <td>${item.value}</td>
                    </tr>
                `).join("")}
            </tbody>
        </table>
    `;
}

// -------------------------
// VARIANTS (SELECT)
// -------------------------
function renderVariants(variants) {
    const select = document.getElementById("volume-select");
    if (!select || !variants) return;

    const values = variants.Value || [];
    const prices = variants.Price || [];
    const unit = variants.Unit || "";

    if (values.length === 0) {
        select.innerHTML = "";
        return;
    }

    select.innerHTML = values.map((v, i) => {
        const price = prices[i] ?? 0;

        return `
            <option value="${v}">
                ${v} ${unit} — ${price} ₽
            </option>
        `;
    }).join("");
}