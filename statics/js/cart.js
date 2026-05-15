function getCart() {
    return JSON.parse(localStorage.getItem("cart")) || [];
}

function saveCart(cart) {
    localStorage.setItem("cart", JSON.stringify(cart));
}

export function updateCartUI() {
    const cart = getCart();

    const icon = document.getElementById("cart-icon");
    const count = document.getElementById("cart-count");

    if (!icon || !count) return;

    icon.style.display = cart.length ? "flex" : "none";

    count.innerText = cart.reduce((s, i) => s + i.quantity, 0);
}

export function initCart() {

    const icon = document.getElementById("cart-icon");
    const modal = document.getElementById("cart-modal");
    const close = document.querySelector(".cart-close");

    updateCartUI();

    icon?.addEventListener("click", () => {
        modal.style.display = "flex";
        renderCart();
    });

    close?.addEventListener("click", () => {
        modal.style.display = "none";
    });

    window.addEventListener("click", (e) => {
        if (e.target === modal) modal.style.display = "none";
    });

       // ADD TO CART (Обновленный и безопасный)
    document.addEventListener("click", (e) => {
        if (e.target.id !== "submit") return;

        const select = document.getElementById("volume-select");
        if (!select) return;

        // Проверяем, есть ли вообще доступные варианты для выбора
        if (select.options.length === 0 || select.selectedIndex === -1) {
            alert("Товар временно недоступен для заказа");
            return;
        }

        const opt = select.options[select.selectedIndex];

        const item = {
            variant_id: Number(opt.dataset.id),
            title: document.getElementById("modal-title").textContent,
            image: document.getElementById("main-image").src,
            url: window.location.pathname,
            volume: opt.value,
            price: Number(opt.dataset.price),
            quantity: 1
        };

        const cart = getCart();

        const exist = cart.find(i =>
            i.variant_id === item.variant_id &&
            i.volume === item.volume
        );

        exist ? exist.quantity++ : cart.push(item);

        saveCart(cart);
        updateCartUI();
    });


    // ORDER
    document.addEventListener("submit", async (e) => {

        if (e.target.id !== "order-form") return;

        e.preventDefault();

        const cart = getCart();
        if (!cart.length) return alert("Корзина пуста");

const order = {
    customer: Object.fromEntries(new FormData(e.target)),
    items: cart.map(i => ({
        id: 0,
        variant_id: Number(i.variant_id) || 0,
        name: i.title || "",
        price: Number(i.price) || 0,
        quantity: Number(i.quantity) || 1
    })),
    total: String(cart.reduce((s, i) => s + i.price * i.quantity, 0)),
    createdAt: new Date().toISOString()
};

        const res = await fetch("/api/orders", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(order)
        });

        if (res.ok) {
            localStorage.removeItem("cart");
            updateCartUI();
            modal.style.display = "none";
            alert("Заказ оформлен");
        }
    });
}
function renderCart() {
    const cart = getCart();
    const container = document.getElementById("cart-items");
    const total = document.getElementById("cart-total");

    if (!container) return;

    if (!cart.length) {
        container.innerHTML = "Корзина пуста";
        total.innerText = "0";
        return;
    }

    let sum = 0;

    container.innerHTML = cart.map((i, idx) => {

        sum += i.price * i.quantity;
// Внутри метода map() функции renderCart:
    const unitText = i.unit && i.unit.trim() !== "" ? ` ${i.unit}` : ""; // Добавьте i.unit в ваш объект сохранения, если нужно, или проверяйте i.volume
        return `
            <div class="cart-item">

                <div class="cart-item-info">
                    <a href="${i.url}" class="cart-item-title">
                        ${i.title}
                    </a>

                    <div class="cart-item-sub">
                        ${i.volume}${unitText} — ${i.price} ₽
                    </div>
                </div>

                <div class="cart-controls">
                    <button onclick="changeQty(${idx}, -1)">-</button>
                    <span>${i.quantity}</span>
                    <button onclick="changeQty(${idx}, 1)">+</button>
                </div>

                <button class="cart-remove" onclick="removeItem(${idx})">
                    ×
                </button>

            </div>
        `;
    }).join("");

    total.innerText = sum;
}

window.removeItem = function (i) {
    const cart = getCart();
    cart.splice(i, 1);
    saveCart(cart);
    updateCartUI();
    renderCart();
};

window.changeQty = function (index, delta) {
    const cart = getCart();

    if (!cart[index]) return;

    cart[index].quantity += delta;

    if (cart[index].quantity <= 0) {
        cart.splice(index, 1);
    }

    saveCart(cart);
    updateCartUI();
    renderCart(); // важно обновить отображение
};