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
        // Сохраняем чистую цену в data-price для JS
        return `
            <option value="${v}" data-price="${price}">
                ${v} ${unit} — ${price} ₽
            </option>
        `;
    }).join("");
}

window.updateCartUI = function() {
    const cart = JSON.parse(localStorage.getItem('cart')) || [];
    const cartIcon = document.getElementById('cart-icon');
    const cartCountElement = document.getElementById('cart-count');
    const cartItemsContainer = document.getElementById('cart-items');
    const cartTotalElement = document.getElementById('cart-total');

    if (cartIcon) cartIcon.style.display = cart.length > 0 ? 'flex' : 'none';
    
    const totalQty = cart.reduce((sum, item) => sum + item.quantity, 0);
    if (cartCountElement) cartCountElement.innerText = totalQty;

    if (cartItemsContainer && cartTotalElement) {
        if (cart.length === 0) {
            cartItemsContainer.innerHTML = 'Корзина пуста';
            cartTotalElement.innerText = '0';
        } else {
            let totalSum = 0;
            cartItemsContainer.innerHTML = cart.map((item, index) => {
                totalSum += item.price * item.quantity;
                return `
            <div class="cart-item" style="display: flex; align-items: center; gap: 10px; border-bottom: 1px solid #eee; padding: 10px 0;">
                <div style="flex: 1;">
                    <!-- Ссылка на товар -->
                    <a href="${item.url}" style="text-decoration: none; color: #333; display: block;">
                        <div style="font-weight: bold; font-size: 14px; margin-bottom: 2px;">${item.title}</div>
                    </a>
                    <div style="font-size: 12px; color: #666;">${item.volume} — ${item.price} ₽</div>
                </div>
                
                <div style="display: flex; align-items: center; gap: 8px;">
                    <button onclick="changeQuantity(${index}, -1)" style="width: 25px; height: 25px; cursor: pointer; border: 1px solid #ddd; background: #f9f9f9; border-radius: 4px;">-</button>
                    <span style="min-width: 20px; text-align: center; font-weight: bold;">${item.quantity}</span>
                    <button onclick="changeQuantity(${index}, 1)" style="width: 25px; height: 25px; cursor: pointer; border: 1px solid #ddd; background: #f9f9f9; border-radius: 4px;">+</button>
                </div>

                <button onclick="removeFromCart(${index})" style="color:red; cursor:pointer; border:none; background:none; font-size: 18px; margin-left: 10px;">&times;</button>
            </div>
        `;
    }).join('');
            cartTotalElement.innerText = totalSum;
        }
    }
};


// Функция удаления
window.removeFromCart = function(index) {
    let cart = JSON.parse(localStorage.getItem('cart')) || [];
    cart.splice(index, 1);
    localStorage.setItem('cart', JSON.stringify(cart));
    window.updateCartUI();
};
window.changeQuantity = function(index, delta) {
    let cart = JSON.parse(localStorage.getItem('cart')) || [];
    
    if (cart[index]) {
        cart[index].quantity += delta;

        // Если количество стало меньше 1 — удаляем товар
        if (cart[index].quantity <= 0) {
            cart.splice(index, 1);
        }

        localStorage.setItem('cart', JSON.stringify(cart));
        window.updateCartUI();
    }
};


// Логика открытия/закрытия модалки корзины
// --- ДОБАВЛЕНИЕ В КОРЗИНУ ---
document.addEventListener('click', (e) => {
    // Слушаем клик по кнопке "Добавить", так как она может пересоздаваться
    if (e.target && e.target.id === 'submit') {
        const select = document.getElementById("volume-select");
        
        if (!select || select.options.length === 0) {
            alert("Пожалуйста, выберите вариант товара");
            return;
        }

        const selectedOption = select.options[select.selectedIndex];
        
        // Собираем данные
        const productData = {
            // Берем название и картинку прямо из модалки
            title: document.getElementById("modal-title").textContent,
            image: currentPhoto,
            url: window.location.pathname,
            // Данные из селекта (теперь берутся корректно)
            volume: selectedOption.value, 
            price: parseInt(selectedOption.dataset.price) || 0,
            quantity: 1
        };

        // Сохранение
        let cart = JSON.parse(localStorage.getItem("cart")) || [];
        
        // Проверяем, есть ли такой же товар с ТАКИМ ЖЕ объемом
        const existingItem = cart.find(item => 
            item.title === productData.title && item.volume === productData.volume
        );

        if (existingItem) {
            existingItem.quantity += 1;
        } else {
            cart.push(productData);
        }

        localStorage.setItem("cart", JSON.stringify(cart));

        // Обратная связь на кнопке
        const originalText = e.target.innerText;
        e.target.innerText = "Добавлено!";
        e.target.style.backgroundColor = "#4CAF50"; // Зеленый цвет для наглядности
        
        setTimeout(() => {
            e.target.innerText = originalText;
            e.target.style.backgroundColor = "";
        }, 1500);

        // Обновляем иконку и содержимое корзины
        window.updateCartUI();
    }
});

// --- ЛОГИКА МОДАЛЬНОГО ОКНА КОРЗИНЫ ---
document.addEventListener('DOMContentLoaded', () => {
    
    const cartIcon = document.getElementById('cart-icon');
    const cartModal = document.getElementById('cart-modal');
    const cartClose = document.querySelector('.cart-close');

    // Открытие корзины
    if (cartIcon) {
        cartIcon.addEventListener('click', () => {
            cartModal.style.display = 'flex'; // Используем flex для центрирования
            window.updateCartUI(); // Обновляем данные при открытии
        });
    }

    // Закрытие корзины на крестик
    if (cartClose) {
        cartClose.addEventListener('click', () => {
            cartModal.style.display = 'none';
        });
    }

    // Закрытие корзины при клике ВНЕ контента (на темный фон)
    window.addEventListener('click', (e) => {
        if (e.target === cartModal) {
            cartModal.style.display = 'none';
        }
    });

    // Инициализация (чтобы иконка появилась сразу, если в корзине уже что-то есть)
    window.updateCartUI();
});

// --- ОТПРАВКА ЗАКАЗА ---
const orderForm = document.getElementById('order-form');

if (orderForm) {
    orderForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const cart = JSON.parse(localStorage.getItem('cart')) || [];
        
        if (cart.length === 0) {
            alert("Ваша корзина пуста");
            return;
        }

        // Кнопка отправки (чтобы заблокировать от повторных кликов)
        const submitBtn = orderForm.querySelector('.checkout-button');
        submitBtn.disabled = true;
        submitBtn.innerText = "Отправка...";

        const formData = new FormData(orderForm);
        
        // Формируем объект для отправки
        const orderData = {
            customer: Object.fromEntries(formData.entries()), // Имя, телефон, комментарий
            items: cart,                                     // Список товаров
            total: document.getElementById('cart-total').innerText, // Итоговая сумма
            createdAt: new Date().toISOString()
        };

        try {
            // ОТПРАВКА НА СЕРВЕР
            const response = await fetch('/api/orders', { // Укажите здесь ваш URL
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(orderData)
            });

            if (response.ok) {
                const result = await response.json();
                alert("Спасибо! Ваш заказ #" + (result.id || "") + " успешно оформлен.");
                
                // Очистка после успеха
                localStorage.removeItem('cart');
                window.updateCartUI();
                document.getElementById('cart-modal').style.display = 'none';
                orderForm.reset();
            } else {
                throw new Error('Ошибка сервера');
            }
        } catch (error) {
            console.error("Ошибка при отправке заказа:", error);
            alert("Произошла ошибка при отправке. Пожалуйста, попробуйте позже.");
        } finally {
            // Возвращаем кнопку в исходное состояние
            submitBtn.disabled = false;
            submitBtn.innerText = "Оформить заказ";
        }
    });
}

