import state from "./state.js";
import { openProductModal } from "./product.js";

document.addEventListener("click", async (e) => {
    const card = e.target.closest(".product_card");
    if (!card) return;

    e.preventDefault();

    const slug = card.dataset.slug;
    if (!slug) return;

    state.previousPage = window.location.pathname;
    state.productEntryMode = "spa";

    history.pushState({}, "", `/product/${slug}`);

    const res = await fetch(`/api/product/${slug}`);
    const product = await res.json();

    openProductModal(product);
});

export function closeProductModal() {
    const modal = document.getElementById("product-modal");
    modal.style.display = "none";

    if (state.productEntryMode === "spa") {
        history.pushState({}, "", state.previousPage || "/");
        return;
    }

    window.location.href = "/";
}